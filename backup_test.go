package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBackupLedger(t *testing.T) {
	now := time.Date(2026, 9, 14, 13, 45, 30, 0, time.UTC)

	tests := []struct {
		name        string
		ledger      string // contenido del ledger; "" = no crear el archivo
		wantErr     string
		wantCreated string // nombre de archivo esperado dentro de dir
	}{
		{
			name: "backup_correcto",
			ledger: `{"books":[{"id":"el-aleph","title":"El Aleph","status":"leido","rating":5,
				"added":"2026-09-01T10:00:00Z"}]}`,
			wantCreated: "ledger-20260914-134530.json",
		},
		{
			name:    "ledger_inexistente_es_error",
			ledger:  "",
			wantErr: "no such file or directory",
		},
		{
			name:    "ledger_corrupto_es_error_y_no_crea_backup",
			ledger:  `{"books": [`,
			wantErr: "parseando",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			srcPath := filepath.Join(dir, "ledger.json")
			if tt.ledger != "" {
				if err := os.WriteFile(srcPath, []byte(tt.ledger), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			outDir := filepath.Join(dir, "backups")

			got, err := BackupLedger(srcPath, outDir, now)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error = %v, queríamos substring %q", err, tt.wantErr)
				}
				entries, _ := os.ReadDir(outDir)
				if len(entries) != 0 {
					t.Errorf("no debía crearse nada en %s, hay %d entradas", outDir, len(entries))
				}
				return
			}
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if got != filepath.Join(outDir, tt.wantCreated) {
				t.Errorf("ruta devuelta = %q, queríamos %q", got, filepath.Join(outDir, tt.wantCreated))
			}
		})
	}
}

func TestBackupLedgerContentIsParseableAndIdentical(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "ledger.json")
	ledger := `{"books":[{"id":"el-aleph","title":"El Aleph","author":"Jorge Luis Borges",
		"pages":210,"status":"leido","rating":5,"added":"2026-09-01T10:00:00Z"}]}`
	if err := os.WriteFile(srcPath, []byte(ledger), 0o600); err != nil {
		t.Fatal(err)
	}
	outDir := filepath.Join(dir, "backups")
	now := time.Date(2026, 9, 14, 8, 0, 0, 0, time.FixedZone("X", 2*3600)) // no-UTC: el nombre va en UTC

	got, err := BackupLedger(srcPath, outDir, now)
	if err != nil {
		t.Fatal(err)
	}
	if base := filepath.Base(got); base != "ledger-20260914-060000.json" {
		t.Errorf("timestamp no UTC: %q", base)
	}

	data, err := os.ReadFile(got)
	if err != nil {
		t.Fatal(err)
	}
	// El backup parsea con la misma gramática del ledger.
	var shelf Shelf
	if err := json.Unmarshal(data, &shelf); err != nil {
		t.Fatalf("backup no parseable: %v", err)
	}
	if len(shelf.Books) != 1 || shelf.Books[0].ID != "el-aleph" || shelf.Books[0].Rating != 5 {
		t.Errorf("contenido inesperado: %s", data)
	}
	// Mismo contenido lógico que la fuente parseada.
	srcShelf, err := LoadShelf(srcPath)
	if err != nil {
		t.Fatal(err)
	}
	want, _ := srcShelf.Marshal()
	if string(data) != string(want) {
		t.Errorf("backup difiere del ledger:\n got %s\nwant %s", data, want)
	}
	// Sin residuos .tmp.
	entries, _ := os.ReadDir(outDir)
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".tmp") {
			t.Errorf("residuo temp: %s", e.Name())
		}
	}
	// Permisos 0600.
	fi, err := os.Stat(got)
	if err != nil {
		t.Fatal(err)
	}
	if perm := fi.Mode().Perm(); perm != 0o600 {
		t.Errorf("permisos = %o, queríamos 600", perm)
	}
}

func TestBackupLedgerTimestampCollisionGetsSuffix(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "ledger.json")
	if err := os.WriteFile(srcPath, []byte(`{"books":[]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	outDir := filepath.Join(dir, "backups")
	now := time.Date(2026, 9, 14, 13, 45, 30, 0, time.UTC)

	first, err := BackupLedger(srcPath, outDir, now)
	if err != nil {
		t.Fatal(err)
	}
	second, err := BackupLedger(srcPath, outDir, now)
	if err != nil {
		t.Fatal(err)
	}
	third, err := BackupLedger(srcPath, outDir, now)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		filepath.Join(outDir, "ledger-20260914-134530.json"),
		filepath.Join(outDir, "ledger-20260914-134530-1.json"),
		filepath.Join(outDir, "ledger-20260914-134530-2.json"),
	}
	got := []string{first, second, third}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("backup %d = %q, queríamos %q", i, got[i], want[i])
		}
	}
	entries, _ := os.ReadDir(outDir)
	if len(entries) != 3 {
		t.Errorf("esperábamos 3 backups, hay %d", len(entries))
	}
}
