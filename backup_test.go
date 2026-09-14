package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
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

// writeBackup crea en dir un backup con nombre ledger-<ts>.json.
func writeBackup(t *testing.T, dir, ts string) string {
	t.Helper()
	path := filepath.Join(dir, "ledger-"+ts+".json")
	if err := os.WriteFile(path, []byte(`{"books":[]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestPruneBackups(t *testing.T) {
	tests := []struct {
		name       string
		timestamps []string // backups pre-creados, en orden de creación
		keep       int
		mkdir      bool     // false = dir no existe
		wantGone   int      // cantidad de eliminados
		wantLeft   []string // nombres que deben sobrevivir
	}{
		{
			name:       "cinco_backups_keep_dos_borra_los_tres_mas_viejos",
			timestamps: []string{"20260910-000000", "20260911-000000", "20260912-000000", "20260913-000000", "20260914-000000"},
			keep:       2,
			mkdir:      true,
			wantGone:   3,
			wantLeft:   []string{"ledger-20260913-000000.json", "ledger-20260914-000000.json"},
		},
		{
			name:       "keep_cero_no_borra_nada",
			timestamps: []string{"20260910-000000", "20260911-000000"},
			keep:       0,
			mkdir:      true,
			wantGone:   0,
			wantLeft:   []string{"ledger-20260910-000000.json", "ledger-20260911-000000.json"},
		},
		{
			name:       "keep_negativo_no_borra_nada",
			timestamps: []string{"20260910-000000"},
			keep:       -3,
			mkdir:      true,
			wantGone:   0,
			wantLeft:   []string{"ledger-20260910-000000.json"},
		},
		{
			name:       "keep_mayor_que_cantidad_no_borra_nada",
			timestamps: []string{"20260910-000000", "20260911-000000"},
			keep:       10,
			mkdir:      true,
			wantGone:   0,
			wantLeft:   []string{"ledger-20260910-000000.json", "ledger-20260911-000000.json"},
		},
		{
			name:     "dir_inexistente_es_noop_sin_error",
			keep:     2,
			mkdir:    false,
			wantGone: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			dir := filepath.Join(root, "backups")
			if tt.mkdir {
				if err := os.MkdirAll(dir, 0o700); err != nil {
					t.Fatal(err)
				}
			}
			created := map[string]bool{}
			for _, ts := range tt.timestamps {
				created[filepath.Base(writeBackup(t, dir, ts))] = true
			}
			// Otros archivos que NUNCA deben ser tocados.
			otros := []string{"otro.txt", "ledger.json", "ledger-nota.md", "ledger-.json"}
			if tt.mkdir {
				for _, n := range otros {
					if err := os.WriteFile(filepath.Join(dir, n), []byte("x"), 0o600); err != nil {
						t.Fatal(err)
					}
				}
			}

			removed, err := PruneBackups(dir, tt.keep)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if len(removed) != tt.wantGone {
				t.Errorf("eliminados = %d (%v), queríamos %d", len(removed), removed, tt.wantGone)
			}
			if !tt.mkdir {
				return // nada más que verificar en el no-op
			}
			// Verificación sobre el filesystem.
			entries, err := os.ReadDir(dir)
			if err != nil {
				t.Fatal(err)
			}
			got := map[string]bool{}
			for _, e := range entries {
				got[e.Name()] = true
			}
			for _, n := range tt.wantLeft {
				if !got[n] {
					t.Errorf("falta el backup que debía sobrevivir: %s", n)
				}
			}
			for _, n := range otros {
				if !got[n] {
					t.Errorf("PruneBackups borró un archivo que no matchea el patrón: %s", n)
				}
			}
			// Cantidad exacta de ledger-*.json restantes.
			remaining := 0
			for n := range got {
				if isLedgerBackupName(n) {
					remaining++
				}
			}
			if want := len(tt.timestamps) - tt.wantGone; remaining != want {
				t.Errorf("quedan %d ledger-*.json, queríamos %d", remaining, want)
			}
		})
	}
}

func TestPruneBackupsUnorderedNames(t *testing.T) {
	dir := t.TempDir()
	// Nombres insertados en desorden; el orden por nombre es el cronológico.
	created := map[string]bool{}
	for _, ts := range []string{"20260913-000000", "20260901-000000", "20260920-000000", "20260907-000000"} {
		created[filepath.Base(writeBackup(t, dir, ts))] = true
	}
	removed, err := PruneBackups(dir, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(removed) != 2 {
		t.Fatalf("eliminados = %v, queríamos 2", removed)
	}
	for _, r := range removed {
		base := filepath.Base(r)
		if base != "ledger-20260901-000000.json" && base != "ledger-20260907-000000.json" {
			t.Errorf("borró un backup que debía sobrevivir: %s", base)
		}
	}
	entries, _ := os.ReadDir(dir)
	var got []string
	for _, e := range entries {
		got = append(got, e.Name())
	}
	sort.Strings(got)
	want := []string{"ledger-20260913-000000.json", "ledger-20260920-000000.json"}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("sobreviviente %d = %s, queríamos %s", i, got[i], want[i])
		}
	}
}

func TestPruneBackupsSuffixCollisionOrdering(t *testing.T) {
	dir := t.TempDir()
	// Misma colisión de segundo que produce BackupLedger: -1 y -2 ordenan
	// después del nombre base.
	for _, n := range []string{"ledger-20260914-134530.json", "ledger-20260914-134530-1.json", "ledger-20260914-134530-2.json"} {
		if err := os.WriteFile(filepath.Join(dir, n), []byte(`{}`), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	removed, err := PruneBackups(dir, 1)
	if err != nil {
		t.Fatal(err)
	}
	// Lexicográficamente '-' < '.', así que los sufijos ordenan ANTES que
	// el nombre base: sobrevive el base, se eliminan -1 y -2.
	if len(removed) != 2 || filepath.Base(removed[0]) != "ledger-20260914-134530-1.json" ||
		filepath.Base(removed[1]) != "ledger-20260914-134530-2.json" {
		t.Errorf("eliminados = %v, queríamos -1 y -2 eliminados, base sobrevive", removed)
	}
}
