package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func capture(t *testing.T, fn func()) (stdout, stderr string) {
	t.Helper()
	outR, outW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	errR, errW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	oldOut, oldErr := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = outW, errW
	defer func() { os.Stdout, os.Stderr = oldOut, oldErr }()
	fn()
	outW.Close()
	errW.Close()
	var outBuf, errBuf bytes.Buffer
	_, _ = io.Copy(&outBuf, outR)
	_, _ = io.Copy(&errBuf, errR)
	return outBuf.String(), errBuf.String()
}

func TestCLIEndToEnd(t *testing.T) {
	tests := []struct {
		name  string
		steps []struct {
			args     []string
			wantCode int
			wantOut  string // substring del stdout
			wantErr  string // substring del stderr
		}
	}{
		{
			name: "flujo completo add-list-status-rating",
			steps: []struct {
				args     []string
				wantCode int
				wantOut  string
				wantErr  string
			}{
				{[]string{"add", "Cazadores de microbios", "--autor", "Paul de Kruif", "--paginas", "352"}, 0, "agregado", ""},
				{[]string{"list"}, 0, "Cazadores de microbios", ""},
				{[]string{"status", "caza"}, 0, "quiero-leer", ""},
				{[]string{"status", "caza", "--set", "leyendo"}, 0, "leyendo", ""},
				{[]string{"status", "caza", "--set", "leido", "--rating", "4"}, 0, "leído (4/5)", ""},
				{[]string{"list", "--status", "leido"}, 0, "★★★★", ""},
			},
		},
		{
			name: "duplicado falla sin corromper el ledger",
			steps: []struct {
				args     []string
				wantCode int
				wantOut  string
				wantErr  string
			}{
				{[]string{"add", "Rayuela"}, 0, "agregado", ""},
				{[]string{"add", "rayuela"}, 1, "", "ya existe"},
			},
		},
		{
			name: "list con status invalido falla",
			steps: []struct {
				args     []string
				wantCode int
				wantOut  string
				wantErr  string
			}{
				{[]string{"add", "Ficciones"}, 0, "agregado", ""},
				{[]string{"list", "--status", "abandonado"}, 1, "", "status inválido"},
			},
		},
		{
			name: "rating sin terminar es rechazado",
			steps: []struct {
				args     []string
				wantCode int
				wantOut  string
				wantErr  string
			}{
				{[]string{"add", "Moby Dick"}, 0, "agregado", ""},
				{[]string{"status", "moby", "--set", "leyendo", "--rating", "5"}, 1, "", "rating solo aplica"},
			},
		},
		{
			name: "libro inexistente",
			steps: []struct {
				args     []string
				wantCode int
				wantOut  string
				wantErr  string
			}{
				{[]string{"status", "no-existe-nada"}, 1, "", "no se encontró"},
			},
		},
		{
			name: "comando desconocido",
			steps: []struct {
				args     []string
				wantCode int
				wantOut  string
				wantErr  string
			}{
				{[]string{"prestar"}, 1, "", "comando desconocido"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ESTANTERIA_FILE", filepath.Join(t.TempDir(), "ledger.json"))
			for i, step := range tt.steps {
				var out, errOut string
				var got int
				out, errOut = capture(t, func() { got = run(step.args) })
				if got != step.wantCode {
					t.Fatalf("paso %d %v: código = %d, want %d", i, step.args, got, step.wantCode)
				}
				if step.wantOut != "" && !strings.Contains(out, step.wantOut) {
					t.Errorf("paso %d %v: stdout %q no contiene %q", i, step.args, out, step.wantOut)
				}
				if step.wantErr != "" && !strings.Contains(errOut, step.wantErr) {
					t.Errorf("paso %d %v: stderr %q no contiene %q", i, step.args, errOut, step.wantErr)
				}
			}
		})
	}
}
