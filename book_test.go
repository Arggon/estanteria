package main

import (
	"testing"
	"time"
)

func TestParseStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Status
		wantErr bool
	}{
		{"canonico quiero-leer", "quiero-leer", StatusQuieroLeer, false},
		{"leyendo", "leyendo", StatusLeyendo, false},
		{"canonico leido", "leido", StatusLeido, false},
		{"con acento", "leído", StatusLeido, false},
		{"mayusculas y espacios", "  LEYENDO ", StatusLeyendo, false},
		{"quiero leer sin guion", "quiero leer", StatusQuieroLeer, false},
		{"invalido", "abandonado", "", true},
		{"vacio", "   ", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStatus(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseStatus(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseStatus(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestStatusDisplay(t *testing.T) {
	tests := []struct {
		status Status
		want   string
	}{
		{StatusQuieroLeer, "quiero-leer"},
		{StatusLeyendo, "leyendo"},
		{StatusLeido, "leído"},
	}
	for _, tt := range tests {
		if got := tt.status.Display(); got != tt.want {
			t.Errorf("Display(%q) = %q, want %q", tt.status, got, tt.want)
		}
	}
}

func TestSetStatus(t *testing.T) {
	finished := func() *time.Time {
		f := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
		return &f
	}
	tests := []struct {
		name        string
		setup       func(*Book)
		newStatus   Status
		rating      int
		wantErr     bool
		wantStatus  Status
		wantRating  int
		wantFinshed bool
	}{
		{
			name:        "quiero-leer a leyendo",
			newStatus:   StatusLeyendo,
			wantStatus:  StatusLeyendo,
			wantFinshed: false,
		},
		{
			name:        "leyendo a leido con rating",
			newStatus:   StatusLeido,
			rating:      4,
			wantStatus:  StatusLeido,
			wantRating:  4,
			wantFinshed: true,
		},
		{
			name:        "leido sin rating fija finished",
			newStatus:   StatusLeido,
			wantStatus:  StatusLeido,
			wantRating:  0,
			wantFinshed: true,
		},
		{
			name:        "rating con status leyendo es rechazado",
			newStatus:   StatusLeyendo,
			rating:      5,
			wantErr:     true,
			wantStatus:  StatusLeyendo,
			wantRating:  0,
			wantFinshed: false,
		},
		{
			name:        "rating 6 fuera de rango",
			newStatus:   StatusLeido,
			rating:      6,
			wantErr:     true,
			wantStatus:  StatusLeyendo,
			wantRating:  0,
			wantFinshed: false,
		},
		{
			name: "leido a leyendo limpia rating y finished",
			setup: func(b *Book) {
				b.Status = StatusLeido
				b.Rating = 5
				b.Finished = finished()
			},
			newStatus:   StatusLeyendo,
			wantStatus:  StatusLeyendo,
			wantRating:  0,
			wantFinshed: false,
		},
		{
			name: "leido a leido actualiza rating",
			setup: func(b *Book) {
				b.Status = StatusLeido
				b.Rating = 2
			},
			newStatus:   StatusLeido,
			rating:      5,
			wantStatus:  StatusLeido,
			wantRating:  5,
			wantFinshed: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Book{ID: "x", Title: "X", Status: StatusLeyendo}
			if tt.setup != nil {
				tt.setup(b)
			}
			err := b.SetStatus(tt.newStatus, tt.rating)
			if (err != nil) != tt.wantErr {
				t.Fatalf("SetStatus(%q, %d) error = %v, wantErr %v", tt.newStatus, tt.rating, err, tt.wantErr)
			}
			if b.Status != tt.wantStatus {
				t.Errorf("status = %q, want %q", b.Status, tt.wantStatus)
			}
			if b.Rating != tt.wantRating {
				t.Errorf("rating = %d, want %d", b.Rating, tt.wantRating)
			}
			if (b.Finished != nil) != tt.wantFinshed {
				t.Errorf("finished = %v, wantFinished %v", b.Finished, tt.wantFinshed)
			}
		})
	}
}
