package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// LedgerPath returns the ledger file: $ESTANTERIA_FILE or ~/.estanteria.json.
func LedgerPath() (string, error) {
	if p := os.Getenv("ESTANTERIA_FILE"); p != "" {
		return p, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("no se pudo resolver el home: %w", err)
	}
	return filepath.Join(home, ".estanteria.json"), nil
}

// LoadShelf reads the ledger; a missing file is an empty shelf.
func LoadShelf(path string) (*Shelf, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return NewShelf(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("leyendo %s: %w", path, err)
	}
	shelf, err := ParseShelf(data)
	if err != nil {
		return nil, fmt.Errorf("parseando %s: %w", path, err)
	}
	return shelf, nil
}

// SaveShelf writes the ledger atomically: temp file in the same directory,
// fsync, then rename over the target so readers never see a partial file.
func SaveShelf(path string, shelf *Shelf) error {
	data, err := shelf.Marshal()
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	f, err := os.CreateTemp(dir, ".estanteria-*.tmp")
	if err != nil {
		return fmt.Errorf("creando temp: %w", err)
	}
	tmp := f.Name()
	defer os.Remove(tmp) // no-op after a successful rename
	if _, err := f.Write(data); err != nil {
		f.Close()
		return fmt.Errorf("escribiendo temp: %w", err)
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return fmt.Errorf("fsync temp: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("cerrando temp: %w", err)
	}
	if err := os.Chmod(tmp, 0o600); err != nil {
		return fmt.Errorf("chmod temp: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("renombrando sobre %s: %w", path, err)
	}
	return nil
}
