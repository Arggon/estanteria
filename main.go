// estanteria is a personal reading tracker: a JSON ledger + a minimal CLI.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
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
  estanteria search <consulta>
  estanteria status <libro> [--set quiero-leer|leyendo|leido] [--rating 1-5]
  estanteria search <consulta>
  estanteria stats
  estanteria serve [--addr 127.0.0.1:8080]
  estanteria backup [--dir RUTA]

El libro se busca por id exacto o por prefijo de título (sin distinguir
mayúsculas ni acentos). El rating solo aplica con --set leido. El backup
crea una copia fechada y atómica del ledger (nunca lo modifica).

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
	case "search":
		err = cmdSearch(rest)
	case "stats":
		err = cmdStats(rest)
	case "serve":
		err = cmdServe(rest)
	case "backup":
		err = cmdBackup(rest)
	default:
		err = fmt.Errorf("comando desconocido %q", cmd)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "estanteria: %v\n", err)
		return 1
	}
	return 0
}

// cmdServe arranca la vista web de SOLO LECTURA (ADR 0003): bindea a
// 127.0.0.1 por defecto y recarga el ledger en cada request vía ServeHandler.
func cmdServe(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	addr := fs.String("addr", "127.0.0.1:8080", "dirección de escucha")
	if partitionArgs(fs, args) != nil {
		return fmt.Errorf("flags inválidas (mirá estanteria --help)")
	}
	path, err := LedgerPath()
	if err != nil {
		return err
	}
	fmt.Printf("estantería en http://%s (solo lectura, ledger: %s)\n", *addr, path)
	return http.ListenAndServe(*addr, ServeHandler(path))
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

// cmdSearch lists every book whose title or author contains the query
// (case/accent-insensitive). It is a pure read: the ledger is never saved.
func cmdSearch(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("uso: estanteria search <consulta>")
	}
	path, err := LedgerPath()
	if err != nil {
		return err
	}
	shelf, err := LoadShelf(path)
	if err != nil {
		return err
	}
	books := shelf.Search(args[0])
	if len(books) == 0 {
		fmt.Printf("(sin resultados para %q)\n", args[0])
	} else {
		for _, b := range books {
			line := fmt.Sprintf("%-14s %s %s %s",
				b.ID, padRight(trunc(b.Title, 40), 40), padRight(trunc(b.Author, 20), 20), b.Status.Display())
			if b.Rating != 0 {
				line += " " + strings.Repeat("★", b.Rating)
			}
			fmt.Println(line)
		}
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

// cmdBackup makes a dated, atomic copy of the ledger. It only ever reads the
// main ledger; the copy goes under dir (default ~/.estanteria-backups).
func cmdBackup(args []string) error {
	fs := flag.NewFlagSet("backup", flag.ContinueOnError)
	dirFlag := fs.String("dir", "", "directorio destino de los backups")
	if partitionArgs(fs, args) != nil {
		return fmt.Errorf("flags inválidas (mirá estanteria --help)")
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("uso: estanteria backup [--dir RUTA]")
	}
	dir := *dirFlag
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("no se pudo resolver el home: %w", err)
		}
		dir = filepath.Join(home, ".estanteria-backups")
	}
	path, err := LedgerPath()
	if err != nil {
		return err
	}
	backupPath, err := BackupLedger(path, dir, time.Now())
	if err != nil {
		return err
	}
	fmt.Println(backupPath)
	return nil
}

// cmdStats prints a read-only aggregate of the shelf. It only loads the
// ledger; it never saves (pure read — stats nunca muta el ledger).
// Months without pages are omitted from the output.
func cmdStats(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("uso: estanteria stats")
	}
	path, err := LedgerPath()
	if err != nil {
		return err
	}
	shelf, err := LoadShelf(path)
	if err != nil {
		return err
	}
	st := ComputeStats(shelf.Books, time.Now())
	fmt.Printf("total: %d\n", st.Total)
	fmt.Printf("quiero-leer: %d\n", st.ByStatus[StatusQuieroLeer])
	fmt.Printf("leyendo: %d\n", st.ByStatus[StatusLeyendo])
	fmt.Printf("leído: %d\n", st.ByStatus[StatusLeido])
	fmt.Printf("rating promedio: %.2f/5\n", st.AvgRating)
	if len(st.PagesPerMonth) == 0 {
		fmt.Println("páginas por mes: (sin datos)")
		return nil
	}
	fmt.Println("páginas por mes:")
	months := make([]string, 0, len(st.PagesPerMonth))
	for m := range st.PagesPerMonth {
		months = append(months, m)
	}
	sort.Strings(months)
	for _, m := range months {
		fmt.Printf("  %s: %d\n", m, st.PagesPerMonth[m])
	}
	return nil
}
