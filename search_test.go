package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func searchTestShelf() *Shelf {
	s := NewShelf()
	s.Books = append(s.Books,
		Book{ID: "el-aleph", Title: "El Aleph", Author: "Jorge Luis Borges", Status: StatusLeido, Rating: 5},
		Book{ID: "cronica-de-una-muerte", Title: "Crónica de una muerte", Author: "Gabriel García Márquez", Status: StatusLeyendo},
		Book{ID: "la-osiade", Title: "La Osiade", Author: "María Dueñas", Status: StatusLeyendo},
		Book{ID: "ficciones", Title: "Ficciones", Author: "Jorge Luis Borges", Status: StatusQuieroLeer},
	)
	return s
}

func TestSearch(t *testing.T) {
	ids := func(books []Book) []string {
		out := make([]string, 0, len(books))
		for _, b := range books {
			out = append(out, b.ID)
		}
		return out
	}
	eq := func(a, b []string) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
	tests := []struct {
		name  string
		query string
		want  []string
	}{
		{"match por titulo", "aleph", []string{"el-aleph"}},
		{"match por autor", "borges", []string{"el-aleph", "ficciones"}},
		{"query con acentos sobre datos con acentos", "Crónica", []string{"cronica-de-una-muerte"}},
		{"query sin acentos sobre datos con acentos", "cronica", []string{"cronica-de-una-muerte"}},
		{"query con acentos sobre datos sin acentos", "maría", []string{"la-osiade"}},
		{"contains no solo prefijo", "muerte", []string{"cronica-de-una-muerte"}},
		{"sin match", "zzz", []string{}},
		{"query vacia", "", []string{}},
		{"query solo whitespace", "   ", []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := searchTestShelf()
			before := len(s.Books)
			got := s.Search(tt.query)
			if !eq(ids(got), tt.want) {
				t.Errorf("Search(%q) = %v, want %v", tt.query, ids(got), tt.want)
			}
			if len(s.Books) != before {
				t.Errorf("Search(%q) mutó el shelf: %d libros, want %d", tt.query, len(s.Books), before)
			}
			// Orden por título.
			for i := 1; i < len(got); i++ {
				if got[i-1].Title > got[i].Title {
					t.Errorf("Search(%q) no está ordenado por título: %q > %q", tt.query, got[i-1].Title, got[i].Title)
				}
			}
		})
	}
}

// TestSearchCLIDoesNotMutateLedger verifies the pure-read invariant end to
// end: after `search`, the ledger file is byte-for-byte identical.
func TestSearchCLIDoesNotMutateLedger(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ledger.json")
	t.Setenv("ESTANTERIA_FILE", path)
	if code := run([]string{"add", "Crónica de una muerte", "--autor", "García Márquez"}); code != 0 {
		t.Fatalf("add exit = %d, want 0", code)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	stdout, _ := capture(t, func() {
		if code := run([]string{"search", "cronica"}); code != 0 {
			t.Errorf("search exit = %d, want 0", code)
		}
	})
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Error("search modificó el ledger")
	}
	if !strings.Contains(stdout, "Crónica de una muerte") || !strings.Contains(stdout, "total: 1") {
		t.Errorf("search stdout = %q, want fila + total: 1", stdout)
	}
}
