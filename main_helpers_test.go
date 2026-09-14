package main

import (
	"strings"
	"testing"
)

func TestPadRight(t *testing.T) {
	tests := []struct {
		name  string
		input string
		width int
		want  int // rune length of the output
	}{
		{"ascii corto", "abc", 10, 10},
		{"acentos paddean por rune", "Crónica", 10, 10},
		{"ya excede el ancho", "Demasiado largo", 5, 15},
		{"exacto", "justo1234", 9, 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := padRight(tt.input, tt.width)
			if n := len([]rune(got)); n != tt.want {
				t.Errorf("padRight(%q, %d) tiene %d runes, want %d", tt.input, tt.width, n, tt.want)
			}
			if !strings.HasPrefix(got, tt.input) {
				t.Errorf("padRight(%q, %d) = %q, no conserva el original", tt.input, tt.width, got)
			}
		})
	}
}
