package main

import (
	"strings"
	"testing"
)

func TestSlug(t *testing.T) {
	tests := []struct {
		title string
		want  string
	}{
		{"El Aleph", "el-aleph"},
		{"CAZADORES DE MICROBIOS", "cazadores-de-microbios"},
		{"Crónica de una muerte", "cronica-de-una-muerte"},
		{"  Espacios   múltiples  ", "espacios-multiples"},
		{"Signos!!! ¿y símbolos?", "signos-y-simbolos"},
		{"100 años de soledad", "100-anos-de-soledad"},
	}
	for _, tt := range tests {
		if got := slug(tt.title); got != tt.want {
			t.Errorf("slug(%q) = %q, want %q", tt.title, got, tt.want)
		}
	}
}

func newTestShelf() *Shelf {
	s := NewShelf()
	s.Books = append(s.Books,
		Book{ID: "el-aleph", Title: "El Aleph", Author: "Jorge Luis Borges", Status: StatusLeido, Rating: 5},
		Book{ID: "la-osiade", Title: "La Osiade", Author: "María Dueñas", Status: StatusLeyendo},
		Book{ID: "ficciones", Title: "Ficciones", Status: StatusQuieroLeer},
	)
	return s
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		author  string
		pages   int
		wantErr string
	}{
		{"ok", "Rayuela", "Julio Cortázar", 600, ""},
		{"titulo vacio", "   ", "", 0, "título"},
		{"paginas negativas", "Algo", "", -5, "negativas"},
		{"duplicado", "el aleph", "", 0, "ya existe"},
		{"solo signos", "!!!", "", 0, "id válido"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestShelf()
			before := len(s.Books)
			book, err := s.Add(tt.title, tt.author, tt.pages)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("Add(%q) error = %v, want contiene %q", tt.title, err, tt.wantErr)
				}
				if len(s.Books) != before {
					t.Errorf("el shelf cambió a pesar del error: %d -> %d", before, len(s.Books))
				}
				return
			}
			if err != nil {
				t.Fatalf("Add(%q) error inesperado: %v", tt.title, err)
			}
			if book.Status != StatusQuieroLeer {
				t.Errorf("status inicial = %q, want quiero-leer", book.Status)
			}
			if len(s.Books) != before+1 {
				t.Errorf("el libro no quedó agregado")
			}
		})
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantID  string
		wantErr bool
	}{
		{"por id exacto", "ficciones", "ficciones", false},
		{"por prefijo case-insensitive", "el al", "el-aleph", false},
		{"por prefijo con acentos en la query", "la osi", "la-osiade", false},
		{"titulo con acento en el shelf matchea por folding", "cronica", "cronica-del-pajaro-que-da-cuerda-al-mundo", false},
		{"sin match", "quijote", "", true},
		{"query vacia", "   ", "", true},
	}
	s := newTestShelf()
	// título con acento para probar el folding del lado del shelf
	if _, err := s.Add("Crónica del pájaro que da cuerda al mundo", "", 0); err != nil {
		t.Fatalf("setup: %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := s.Find(tt.query)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Find(%q) = %v, want error", tt.query, b)
				}
				return
			}
			if err != nil {
				t.Fatalf("Find(%q) error inesperado: %v", tt.query, err)
			}
			if b.ID != tt.wantID {
				t.Errorf("Find(%q) = %q, want %q", tt.query, b.ID, tt.wantID)
			}
		})
	}
}

func TestFindAmbiguousPrefixIsDeterministic(t *testing.T) {
	// Dos libros comparten el prefijo "la": el resultado no debe depender
	// del orden del ledger, gana el (título, id) lexicográficamente menor.
	s := NewShelf()
	if _, err := s.Add("Zorba el griego", "", 0); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Add("La sombra del viento", "", 0); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Add("La casa de los espíritus", "", 0); err != nil {
		t.Fatal(err)
	}
	b, err := s.Find("la")
	if err != nil {
		t.Fatalf("Find(la): %v", err)
	}
	if b.ID != "la-casa-de-los-espiritus" {
		t.Errorf("Find(la) = %q, want la-casa-de-los-espiritus (el menor lexicográfico)", b.ID)
	}
}

func TestFilter(t *testing.T) {
	s := newTestShelf()
	if got := len(s.Filter("")); got != 3 {
		t.Errorf("Filter(all) = %d libros, want 3", got)
	}
	if got := len(s.Filter(StatusLeido)); got != 1 {
		t.Errorf("Filter(leido) = %d libros, want 1", got)
	}
	leidos := s.Filter(StatusLeido)
	if leidos[0].ID != "el-aleph" {
		t.Errorf("Filter(leido) trajo %q", leidos[0].ID)
	}
	if got := len(s.Filter(Status("no-existe"))); got != 0 {
		t.Errorf("Filter(status inventado) = %d libros, want 0", got)
	}
}
