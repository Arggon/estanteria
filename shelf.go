package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Shelf is the whole collection: the ledger document and its operations.
type Shelf struct {
	Books []Book `json:"books"`
}

func NewShelf() *Shelf { return &Shelf{Books: []Book{}} }

func ParseShelf(data []byte) (*Shelf, error) {
	var s Shelf
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	if s.Books == nil {
		s.Books = []Book{}
	}
	return &s, nil
}

func (s *Shelf) Marshal() ([]byte, error) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

// slug turns a title into an id: lowercase, accent-folded, non-alphanumeric
// runs collapsed to a single dash.
func slug(title string) string {
	folded := foldAccents(strings.ToLower(strings.TrimSpace(title)))
	var b strings.Builder
	lastDash := false
	for _, r := range folded {
		switch {
		case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		default:
			if !lastDash && b.Len() > 0 {
				b.WriteRune('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

// Add appends a book; the id is the slug of the title and must be unique.
func (s *Shelf) Add(title, author string, pages int) (*Book, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, fmt.Errorf("el título no puede estar vacío")
	}
	if pages < 0 {
		return nil, fmt.Errorf("las páginas no pueden ser negativas")
	}
	id := slug(title)
	if id == "" {
		return nil, fmt.Errorf("el título %q no produce un id válido", title)
	}
	if _, err := s.Find(id); err == nil {
		return nil, fmt.Errorf("ya existe un libro con id %q", id)
	}
	book := Book{
		ID:     id,
		Title:  title,
		Author: strings.TrimSpace(author),
		Pages:  pages,
		Status: StatusQuieroLeer,
		Added:  time.Now().UTC().Truncate(time.Second),
	}
	s.Books = append(s.Books, book)
	return &s.Books[len(s.Books)-1], nil
}

// Find matches by exact id first, then by case/accent-insensitive title
// prefix. Among prefix matches it picks the lexicographically smallest
// (title, id) so the result does not depend on ledger order. Empty query is
// an error.
func (s *Shelf) Find(query string) (*Book, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return nil, fmt.Errorf("consulta de búsqueda vacía")
	}
	folded := foldAccents(strings.ToLower(q))
	var best *Book
	for i := range s.Books {
		b := &s.Books[i]
		if b.ID == q {
			return b, nil
		}
		title := foldAccents(strings.ToLower(b.Title))
		if strings.HasPrefix(title, folded) {
			if best == nil || b.Title < best.Title || (b.Title == best.Title && b.ID < best.ID) {
				best = b
			}
		}
	}
	if best != nil {
		return best, nil
	}
	return nil, fmt.Errorf("no se encontró ningún libro para %q", query)
}

// Filter returns the books matching a status (nil Status = all),
// sorted by title for stable output.
func (s *Shelf) Filter(status Status) []Book {
	out := make([]Book, 0, len(s.Books))
	for _, b := range s.Books {
		if status == "" || b.Status == status {
			out = append(out, b)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Title < out[j].Title })
	return out
}
