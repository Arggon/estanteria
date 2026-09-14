// estanteria is a personal reading tracker: a JSON ledger + a minimal CLI.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// partitionArgs parses flags anywhere in args (the stdlib flag package stops
// at the first positional, so `add "Título" --autor X` would fail) and leaves
// the positionals in fs.Args().
func partitionArgs(fs *flag.FlagSet, args []string) error {
	var flags, positional []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if len(a) > 1 && a[0] == '-' {
			flags = append(flags, a)
			if !strings.Contains(a, "=") && i+1 < len(args) {
				next := args[i+1]
				if !strings.HasPrefix(next, "-") || isIntString(next) {
					flags = append(flags, next)
					i++
				}
			}
			continue
		}
		positional = append(positional, a)
	}
	if err := fs.Parse(flags); err != nil {
		return err
	}
	// Segunda pasada sin flags: deja los posicionales en fs.Args()/NArg().
	return fs.Parse(positional)
}

func isIntString(s string) bool {
	neg := len(s) > 1 && s[0] == '-'
	digits := strings.TrimPrefix(s, "-")
	if neg && len(digits) == 0 {
		return false
	}
	for _, r := range digits {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(digits) > 0
}

const usage = `estanteria — tracker de lectura personal

Uso:
  estanteria add <título> [--autor NOMBRE] [--paginas N]
  estanteria list [--status quiero-leer|leyendo|leido]
  estanteria status <libro> [--set quiero-leer|leyendo|leido] [--rating 1-5]

El libro se busca por id exacto o por prefijo de título (sin distinguir
mayúsculas ni acentos). El rating solo aplica con --set leido.

Variables: ESTANTERIA_FILE (default: ~/.estanteria.json)`

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Println(usage)
		return 0
	}
	cmd, rest := args[0], args[1:]
	var err error
	switch cmd {
	case "add":
		err = cmdAdd(rest)
	case "list":
		err = cmdList(rest)
	case "status":
		err = cmdStatus(rest)
	default:
		err = fmt.Errorf("comando desconocido %q", cmd)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "estanteria: %v\n", err)
		return 1
	}
	return 0
}

func cmdAdd(args []string) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	autor := fs.String("autor", "", "autor/a del libro")
	paginas := fs.Int("paginas", 0, "cantidad de páginas")
	if partitionArgs(fs, args) != nil {
		return fmt.Errorf("flags inválidas (mirá estanteria --help)")
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("uso: estanteria add <título> [--autor NOMBRE] [--paginas N]")
	}
	path, err := LedgerPath()
	if err != nil {
		return err
	}
	shelf, err := LoadShelf(path)
	if err != nil {
		return err
	}
	book, err := shelf.Add(fs.Arg(0), *autor, *paginas)
	if err != nil {
		return err
	}
	if err := SaveShelf(path, shelf); err != nil {
		return err
	}
	fmt.Printf("agregado: %s (id: %s, status: %s)\n", book.Title, book.ID, book.Status.Display())
	return nil
}

func cmdList(args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	statusRaw := fs.String("status", "", "filtrar por status")
	if partitionArgs(fs, args) != nil {
		return fmt.Errorf("flags inválidas (mirá estanteria --help)")
	}
	var status Status
	if *statusRaw != "" {
		var err error
		if status, err = ParseStatus(*statusRaw); err != nil {
			return err
		}
	}
	path, err := LedgerPath()
	if err != nil {
		return err
	}
	shelf, err := LoadShelf(path)
	if err != nil {
		return err
	}
	books := shelf.Filter(status)
	if len(books) == 0 {
		fmt.Println("(estantería vacía para este filtro)")
		return nil
	}
	for _, b := range books {
		line := fmt.Sprintf("%-14s %s %s %s",
			b.ID, padRight(trunc(b.Title, 40), 40), padRight(trunc(b.Author, 20), 20), b.Status.Display())
		if b.Rating != 0 {
			line += " " + strings.Repeat("★", b.Rating)
		}
		fmt.Println(line)
	}
	fmt.Printf("total: %d\n", len(books))
	return nil
}

// padRight pads by rune count: %-Ns pads by bytes and misaligns columns
// for accented titles.
func padRight(s string, width int) string {
	r := []rune(s)
	if len(r) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(r))
}

func cmdStatus(args []string) error {
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	setRaw := fs.String("set", "", "cambiar el status del libro")
	ratingRaw := fs.Int("rating", 0, "rating 1-5 (solo con --set leido)")
	if partitionArgs(fs, args) != nil {
		return fmt.Errorf("flags inválidas (mirá estanteria --help)")
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("uso: estanteria status <libro> [--set STATUS] [--rating 1-5]")
	}
	var newStatus Status
	if *setRaw != "" {
		var err error
		if newStatus, err = ParseStatus(*setRaw); err != nil {
			return err
		}
	}
	path, err := LedgerPath()
	if err != nil {
		return err
	}
	shelf, err := LoadShelf(path)
	if err != nil {
		return err
	}
	book, err := shelf.Find(fs.Arg(0))
	if err != nil {
		return err
	}
	if newStatus != "" {
		if err := book.SetStatus(newStatus, *ratingRaw); err != nil {
			return err
		}
		if err := SaveShelf(path, shelf); err != nil {
			return err
		}
	}
	fmt.Printf("%s — %s", book.Title, book.Status.Display())
	if book.Rating != 0 {
		fmt.Printf(" (%d/5)", book.Rating)
	}
	if book.Finished != nil {
		fmt.Printf(" — terminado: %s", book.Finished.Format("2006-01-02"))
	}
	fmt.Println()
	return nil
}

func trunc(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "…"
}
