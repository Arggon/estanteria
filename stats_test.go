package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func bt(year, month, day int) *time.Time {
	t := time.Date(year, time.Month(month), day, 12, 0, 0, 0, time.UTC)
	return &t
}

func TestComputeStats(t *testing.T) {
	// now = 2026-09-15; ventana de 12 meses: 2025-10 .. 2026-09.
	now := time.Date(2026, 9, 15, 10, 0, 0, 0, time.UTC)
	tests := []struct {
		name          string
		books         []Book
		wantTotal     int
		wantByStatus  map[Status]int
		wantAvgRating float64
		wantPages     map[string]int
	}{
		{
			name:          "estantería vacía",
			books:         []Book{},
			wantTotal:     0,
			wantByStatus:  map[Status]int{},
			wantAvgRating: 0,
			wantPages:     map[string]int{},
		},
		{
			name: "sin ratings el promedio es 0",
			books: []Book{
				{ID: "a", Status: StatusLeido, Finished: bt(2026, 9, 1), Pages: 100},
				{ID: "b", Status: StatusLeido, Finished: bt(2026, 8, 1), Pages: 200},
			},
			wantTotal:     2,
			wantByStatus:  map[Status]int{StatusLeido: 2},
			wantAvgRating: 0,
			wantPages:     map[string]int{"2026-09": 100, "2026-08": 200},
		},
		{
			name: "promedio solo sobre libros con rating",
			books: []Book{
				{ID: "a", Status: StatusLeido, Rating: 5, Finished: bt(2026, 9, 1), Pages: 100},
				{ID: "b", Status: StatusLeido, Rating: 3, Finished: bt(2026, 9, 2), Pages: 50},
				{ID: "c", Status: StatusLeyendo},
			},
			wantTotal:     3,
			wantByStatus:  map[Status]int{StatusLeido: 2, StatusLeyendo: 1},
			wantAvgRating: 4.0,
			wantPages:     map[string]int{"2026-09": 150},
		},
		{
			name: "libros sin pages ni finished cuentan en totales pero no suman páginas",
			books: []Book{
				{ID: "a", Status: StatusQuieroLeer},
				{ID: "b", Status: StatusLeido, Rating: 4, Finished: bt(2026, 9, 1)}, // sin pages
				{ID: "c", Status: StatusLeyendo, Pages: 300},                        // sin finished
			},
			wantTotal:     3,
			wantByStatus:  map[Status]int{StatusQuieroLeer: 1, StatusLeyendo: 1, StatusLeido: 1},
			wantAvgRating: 4.0,
			wantPages:     map[string]int{},
		},
		{
			name: "distribución entre meses acumula",
			books: []Book{
				{ID: "a", Status: StatusLeido, Finished: bt(2026, 3, 10), Pages: 100},
				{ID: "b", Status: StatusLeido, Finished: bt(2026, 3, 25), Pages: 250},
				{ID: "c", Status: StatusLeido, Finished: bt(2026, 4, 5), Pages: 90},
			},
			wantTotal:     3,
			wantByStatus:  map[Status]int{StatusLeido: 3},
			wantAvgRating: 0,
			wantPages:     map[string]int{"2026-03": 350, "2026-04": 90},
		},
		{
			name: "borde de la ventana de 12 meses",
			books: []Book{
				{ID: "dentro-inicio", Status: StatusLeido, Finished: bt(2025, 10, 1), Pages: 10},
				{ID: "fuera-antes", Status: StatusLeido, Finished: bt(2025, 9, 30), Pages: 20},
				{ID: "dentro-fin", Status: StatusLeido, Finished: bt(2026, 9, 30), Pages: 30},
				{ID: "fuera-futuro", Status: StatusLeido, Finished: bt(2026, 10, 1), Pages: 40},
			},
			wantTotal:     4,
			wantByStatus:  map[Status]int{StatusLeido: 4},
			wantAvgRating: 0,
			wantPages:     map[string]int{"2025-10": 10, "2026-09": 30},
		},
		{
			name: "promedio con redondeo a 2 decimales",
			books: []Book{
				{ID: "a", Status: StatusLeido, Rating: 5, Finished: bt(2026, 9, 1)},
				{ID: "b", Status: StatusLeido, Rating: 4, Finished: bt(2026, 9, 1)},
				{ID: "c", Status: StatusLeido, Rating: 4, Finished: bt(2026, 9, 1)},
			},
			wantTotal:     3,
			wantByStatus:  map[Status]int{StatusLeido: 3},
			wantAvgRating: 4.33,
			wantPages:     map[string]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeStats(tt.books, now)
			if got.Total != tt.wantTotal {
				t.Errorf("Total = %d, want %d", got.Total, tt.wantTotal)
			}
			if len(got.ByStatus) != len(tt.wantByStatus) {
				t.Fatalf("ByStatus = %v, want %v", got.ByStatus, tt.wantByStatus)
			}
			for s, n := range tt.wantByStatus {
				if got.ByStatus[s] != n {
					t.Errorf("ByStatus[%s] = %d, want %d", s, got.ByStatus[s], n)
				}
			}
			if got.AvgRating != tt.wantAvgRating {
				t.Errorf("AvgRating = %v, want %v", got.AvgRating, tt.wantAvgRating)
			}
			if len(got.PagesPerMonth) != len(tt.wantPages) {
				t.Fatalf("PagesPerMonth = %v, want %v", got.PagesPerMonth, tt.wantPages)
			}
			for m, n := range tt.wantPages {
				if got.PagesPerMonth[m] != n {
					t.Errorf("PagesPerMonth[%s] = %d, want %d", m, got.PagesPerMonth[m], n)
				}
			}
		})
	}
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestCLIStats(t *testing.T) {
	t.Run("stats no muta el ledger e imprime agregados", func(t *testing.T) {
		ledger := filepath.Join(t.TempDir(), "ledger.json")
		t.Setenv("ESTANTERIA_FILE", ledger)
		if got := run([]string{"add", "El Aleph", "--autor", "Jorge Luis Borges", "--paginas", "210"}); got != 0 {
			t.Fatalf("add: código = %d", got)
		}
		if got := run([]string{"status", "el aleph", "--set", "leido", "--rating", "5"}); got != 0 {
			t.Fatalf("status: código = %d", got)
		}
		before := mustRead(t, ledger)
		var out string
		var got int
		out, _ = capture(t, func() { got = run([]string{"stats"}) })
		if got != 0 {
			t.Fatalf("stats: código = %d", got)
		}
		if after := mustRead(t, ledger); before != after {
			t.Errorf("stats mutó el ledger")
		}
		for _, want := range []string{"total: 1", "leído: 1", "rating promedio: 5.00/5"} {
			if !strings.Contains(out, want) {
				t.Errorf("stdout %q no contiene %q", out, want)
			}
		}
	})
	t.Run("stats sobre estantería vacía", func(t *testing.T) {
		t.Setenv("ESTANTERIA_FILE", filepath.Join(t.TempDir(), "ledger.json"))
		var out string
		var got int
		out, _ = capture(t, func() { got = run([]string{"stats"}) })
		if got != 0 {
			t.Fatalf("stats: código = %d", got)
		}
		if !strings.Contains(out, "total: 0") {
			t.Errorf("stdout %q no contiene %q", out, "total: 0")
		}
	})
}
