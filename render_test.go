package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestCmdRender(t *testing.T) {
	dir := t.TempDir()
	ledger := filepath.Join(dir, "ledger.json")
	out := filepath.Join(dir, "out", "index.html")
	// 0644 en el ledger de entrada: render nunca debe mutarlo.
	content := []byte(`{"books":[{"id":"el-aleph","title":"El Aleph","author":"Jorge Luis Borges","pages":210,"status":"leido","rating":5,"added":"2025-11-03T10:15:00Z","finished":"2025-11-20T22:30:00Z"},{"id":"rayuela","title":"Rayuela","author":"Julio Cortázar","pages":640,"status":"leyendo","added":"2025-12-01T09:00:00Z"}]}`)
	if err := os.WriteFile(ledger, content, 0o644); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(ledger)
	if err != nil {
		t.Fatal(err)
	}

	// --out requerido.
	if err := cmdRender(nil); err == nil {
		t.Fatal("render sin --out debería fallar")
	}

	if err := cmdRender([]string{"--out", out, "--file", ledger}); err != nil {
		t.Fatalf("cmdRender: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("leyendo el render: %v", err)
	}
	html := string(data)
	for _, want := range []string{
		"<!DOCTYPE html>",
		"</html>",
		"El Aleph",
		"Jorge Luis Borges",
		"Rayuela",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("el HTML no contiene %q", want)
		}
	}
	if fi, err := os.Stat(out); err != nil {
		t.Fatal(err)
	} else if fi.Mode().Perm() != 0o644 {
		t.Errorf("permisos del render = %v, want -rw-r--r--", fi.Mode().Perm())
	}

	// Escritura atómica: sin .tmp residuales en el directorio de salida.
	entries, err := os.ReadDir(filepath.Dir(out))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".tmp") {
			t.Errorf("temp residual: %s", e.Name())
		}
	}

	// El ledger de origen queda intacto byte a byte.
	after, err := os.ReadFile(ledger)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Error("render mutó el ledger de origen")
	}
}

func TestCmdRenderOverwrite(t *testing.T) {
	dir := t.TempDir()
	ledger := filepath.Join(dir, "ledger.json")
	if err := os.WriteFile(ledger, []byte(`{"books":[{"id":"x","title":"Ficciones","status":"quiero-leer","added":"2026-01-01T00:00:00Z"}]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "index.html")
	if err := cmdRender([]string{"--out", out, "--file", ledger}); err != nil {
		t.Fatal(err)
	}
	// Segunda pasada sobre el mismo --out: rename lo pisa sin error.
	if err := cmdRender([]string{"--out", out, "--file", ledger}); err != nil {
		t.Fatalf("segundo render: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile("(?s)<h1>estantería</h1>").Match(data) {
		t.Error("el HTML regenerado no es la vista esperada")
	}
}

func TestCmdRenderMissingFile(t *testing.T) {
	dir := t.TempDir()
	// Ledger inexistente == estantería vacía (no error), igual que el resto.
	out := filepath.Join(dir, "index.html")
	if err := cmdRender([]string{"--out", out, "--file", filepath.Join(dir, "noexiste.json")}); err != nil {
		t.Fatalf("render con ledger inexistente: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "(vacía)") {
		t.Error("ledger vacío debería renderizar (vacía)")
	}
}
