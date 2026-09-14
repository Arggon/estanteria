package main

import (
	"fmt"
	"strings"
	"time"
	"unicode"
)

// Status is the reading state of a book. Canonical values are ASCII so they
// work as CLI flags; display uses the accented Spanish word.
type Status string

const (
	StatusQuieroLeer Status = "quiero-leer"
	StatusLeyendo    Status = "leyendo"
	StatusLeido      Status = "leido"
)

// ParseStatus accepts the canonical ASCII form and the accented display form.
func ParseStatus(s string) (Status, error) {
	n := foldAccents(strings.ToLower(strings.TrimSpace(s)))
	switch n {
	case "quiero-leer", "quiero leer":
		return StatusQuieroLeer, nil
	case "leyendo":
		return StatusLeyendo, nil
	case "leido", "leído":
		return StatusLeido, nil
	default:
		return "", fmt.Errorf("status inválido %q (quiero-leer | leyendo | leido)", s)
	}
}

func (s Status) Display() string {
	if s == StatusLeido {
		return "leído"
	}
	return string(s)
}

func foldAccents(s string) string {
	out := make([]rune, 0, len(s))
	for _, c := range s {
		r := unicode.ToLower(c)
		switch r {
		case 'á', 'à', 'ä':
			r = 'a'
		case 'é', 'è', 'ë':
			r = 'e'
		case 'í', 'ì', 'ï':
			r = 'i'
		case 'ó', 'ò', 'ö':
			r = 'o'
		case 'ú', 'ù', 'ü':
			r = 'u'
		case 'ñ':
			r = 'n'
		}
		out = append(out, r)
	}
	return string(out)
}

// Book is one entry of the reading tracker.
type Book struct {
	ID       string     `json:"id"`
	Title    string     `json:"title"`
	Author   string     `json:"author,omitempty"`
	Pages    int        `json:"pages,omitempty"`
	Status   Status     `json:"status"`
	Rating   int        `json:"rating,omitempty"`
	Added    time.Time  `json:"added"`
	Finished *time.Time `json:"finished,omitempty"`
}

// SetStatus moves the book to the new status with rating semantics:
// a rating is only meaningful when finishing (leido), and leaving leido
// clears a stale rating.
func (b *Book) SetStatus(s Status, rating int) error {
	if rating != 0 && s != StatusLeido {
		return fmt.Errorf("el rating solo aplica al terminar (status %q)", s.Display())
	}
	if rating != 0 && (rating < 1 || rating > 5) {
		return fmt.Errorf("rating %d fuera de rango (1-5)", rating)
	}
	if b.Status == StatusLeido && s != StatusLeido {
		b.Rating = 0
		b.Finished = nil
	}
	b.Status = s
	if s == StatusLeido {
		if rating != 0 {
			b.Rating = rating
		}
		now := time.Now().UTC().Truncate(time.Second)
		b.Finished = &now
	} else {
		b.Finished = nil
	}
	return nil
}
