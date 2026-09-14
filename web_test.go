package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func mustWriteLedger(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "ledger.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestServeHandlerGet(t *testing.T) {
	path := mustWriteLedger(t, `{"books":[
		{"id":"el-aleph","title":"El Aleph","author":"Jorge Luis Borges","pages":210,"status":"leido","rating":5}
	]}`)
	srv := httptest.NewServer(ServeHandler(path))
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Fatalf("Content-Type = %q, want text/html; charset=utf-8", ct)
	}
	buf := make([]byte, 4096)
	n, _ := resp.Body.Read(buf)
	body := string(buf[:n])
	for _, want := range []string{"El Aleph", "Jorge Luis Borges", "leído", "★★★★★"} {
		if !strings.Contains(body, want) {
			t.Errorf("body no contiene %q", want)
		}
	}
}

func TestServeHandlerMethodNotAllowed(t *testing.T) {
	path := mustWriteLedger(t, `{"books":[]}`)
	srv := httptest.NewServer(ServeHandler(path))
	defer srv.Close()

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		req, err := http.NewRequest(method, srv.URL, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("%s: status = %d, want 405", method, resp.StatusCode)
		}
	}
}

func TestServeHandlerEmptyLedger(t *testing.T) {
	path := mustWriteLedger(t, `{"books":[]}`)
	srv := httptest.NewServer(ServeHandler(path))
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	buf := make([]byte, 4096)
	n, _ := resp.Body.Read(buf)
	if !strings.Contains(string(buf[:n]), "(vacía)") {
		t.Errorf("ledger vacío: body no contiene \"(vacía)\"")
	}
}

func TestRenderHTMLEscapesMaliciousTitle(t *testing.T) {
	books := []Book{{
		ID:     "malo",
		Title:  `<script>alert("x")</script> & <b>negrita</b> café`,
		Status: StatusLeyendo,
	}}
	body := string(RenderHTML(books, time.Now()))
	if strings.Contains(body, "<script>alert") {
		t.Errorf("el título malicioso aparece sin escapar")
	}
	if !strings.Contains(body, "&lt;script&gt;") {
		t.Errorf("body no contiene el título escapado &lt;script&gt;")
	}
	if !strings.Contains(body, "café") {
		t.Errorf("los acentos del título se perdieron")
	}
}

func TestRenderHTMLSections(t *testing.T) {
	fin := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	books := []Book{
		{ID: "a", Title: "Quiero", Status: StatusQuieroLeer},
		{ID: "b", Title: "Leyendo ahora", Status: StatusLeyendo},
		{ID: "c", Title: "Terminado", Status: StatusLeido, Rating: 4, Finished: &fin, Pages: 300},
	}
	body := string(RenderHTML(books, time.Now()))
	for _, want := range []string{"quiero-leer", "leyendo", "leído", "★★★★", "300 págs.", "terminado: 2026-09-01", "total: 3", "promedio de ratings: 4.0"} {
		if !strings.Contains(body, want) {
			t.Errorf("body no contiene %q", want)
		}
	}
}
