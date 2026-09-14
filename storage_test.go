package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeShelf(t *testing.T, path string, shelf *Shelf) {
	t.Helper()
	if err := SaveShelf(path, shelf); err != nil {
		t.Fatalf("SaveShelf: %v", err)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ledger.json")
	in := NewShelf()
	fin := time.Date(2026, 9, 10, 12, 0, 0, 0, time.UTC)
	in.Books = append(in.Books, Book{
		ID: "padre-rico-padre-pobre", Title: "Padre rico, padre pobre",
		Author: "Robert Kiyosaki", Pages: 336, Status: StatusLeido, Rating: 3,
		Added: fin, Finished: &fin,
	})
	writeShelf(t, path, in)

	out, err := LoadShelf(path)
	if err != nil {
		t.Fatalf("LoadShelf: %v", err)
	}
	if len(out.Books) != 1 {
		t.Fatalf("round trip devolvió %d libros, want 1", len(out.Books))
	}
	got := out.Books[0]
	if got.ID != in.Books[0].ID || got.Title != in.Books[0].Title ||
		got.Status != StatusLeido || got.Rating != 3 || got.Pages != 336 {
		t.Errorf("round trip distinto: %+v", got)
	}
	if got.Finished == nil || !got.Finished.Equal(fin) {
		t.Errorf("finished = %v, want %v", got.Finished, fin)
	}
}

func TestLoadShelfMissingFileIsEmptyShelf(t *testing.T) {
	path := filepath.Join(t.TempDir(), "no-existe.json")
	shelf, err := LoadShelf(path)
	if err != nil {
		t.Fatalf("LoadShelf(missing) error = %v, want nil", err)
	}
	if shelf == nil || len(shelf.Books) != 0 {
		t.Errorf("LoadShelf(missing) = %+v, want shelf vacío", shelf)
	}
}

func TestLoadShelfCorruptJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "roto.json")
	if err := os.WriteFile(path, []byte("{esto no es json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadShelf(path); err == nil {
		t.Fatal("LoadShelf(corrupto) error = nil, want error")
	}
}

func TestSaveShelfIsAtomicNoTempLeftovers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ledger.json")

	writeShelf(t, path, NewShelf())
	writeShelf(t, path, func() *Shelf {
		s := NewShelf()
		s.Add("Titulo de prueba", "", 10)
		return s
	}())

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var leftovers []string
	for _, e := range entries {
		if strings.Contains(e.Name(), ".tmp") {
			leftovers = append(leftovers, e.Name())
		}
	}
	if len(leftovers) != 0 {
		t.Errorf("quedaron archivos temp tras el save: %v", leftovers)
	}
	shelf, err := LoadShelf(path)
	if err != nil || len(shelf.Books) != 1 {
		t.Errorf("segundo save no reemplazó el contenido: %+v err=%v", shelf, err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("permisos del ledger = %v, want -rw-------", info.Mode().Perm())
	}
}
