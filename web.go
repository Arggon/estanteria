package main

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"
)

// RenderHTML produce la vista HTML autocontenida de la estantería. Es una
// función PURA sobre los libros (ADR 0004): serve la sirve por request y
// render la escribirá como snapshot. Todo lo que viene del ledger se escapa
// con html.EscapeString — los títulos pueden contener acentos, ángulos o
// markup malicioso.
func RenderHTML(books []Book, generatedAt time.Time) []byte {
	groups := []struct {
		status  Status
		display string
	}{
		{StatusQuieroLeer, "quiero-leer"},
		{StatusLeyendo, "leyendo"},
		{StatusLeido, "leído"},
	}

	var b strings.Builder
	b.WriteString(`<!DOCTYPE html>
<html lang="es">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>estantería</title>
<style>
body{font-family:system-ui,sans-serif;max-width:42rem;margin:2rem auto;padding:0 1rem;color:#222;line-height:1.5}
h1{border-bottom:2px solid #888}
h2{margin-top:2rem;font-size:1.1rem;color:#555}
ul{list-style:none;padding:0}
li{margin:0.4rem 0}
.meta{color:#666}
footer{margin-top:3rem;border-top:1px solid #ccc;padding-top:0.5rem;color:#666;font-size:0.9rem}
</style>
</head>
<body>
<h1>estantería</h1>
`)
	if len(books) == 0 {
		b.WriteString("<p>(vacía)</p>\n")
	}
	for _, g := range groups {
		shown := filterSorted(books, g.status)
		b.WriteString("<h2>" + html.EscapeString(g.display) + "</h2>\n")
		if len(shown) == 0 {
			b.WriteString("<p class=\"meta\">(vacía)</p>\n")
			continue
		}
		b.WriteString("<ul>\n")
		for _, bk := range shown {
			b.WriteString("<li>")
			b.WriteString("<strong>" + html.EscapeString(bk.Title) + "</strong>")
			if bk.Author != "" {
				b.WriteString(" — " + html.EscapeString(bk.Author))
			}
			if bk.Pages > 0 {
				b.WriteString(fmt.Sprintf(" <span class=\"meta\">(%d págs.)</span>", bk.Pages))
			}
			if bk.Rating != 0 {
				b.WriteString(" <span class=\"meta\">" + strings.Repeat("★", bk.Rating) + "</span>")
			}
			if bk.Finished != nil {
				b.WriteString(" <span class=\"meta\">terminado: " + bk.Finished.Format("2006-01-02") + "</span>")
			}
			b.WriteString("</li>\n")
		}
		b.WriteString("</ul>\n")
	}

	total := len(books)
	rated, sum := 0, 0
	for _, bk := range books {
		if bk.Rating != 0 {
			rated++
			sum += bk.Rating
		}
	}
	b.WriteString("<footer>\n")
	fmt.Fprintf(&b, "<p>total: %d libro(s)", total)
	if rated > 0 {
		fmt.Fprintf(&b, " — promedio de ratings: %.1f / 5", float64(sum)/float64(rated))
	}
	b.WriteString("</p>\n")
	fmt.Fprintf(&b, "<p class=\"meta\">generado: %s — solo lectura</p>\n", generatedAt.UTC().Format("2006-01-02 15:04 UTC"))
	b.WriteString("</footer>\n</body>\n</html>\n")
	return []byte(b.String())
}

// filterSorted devuelve los libros de un status, ordenados por título
// (mismo criterio que Shelf.Filter) sin tocar el slice entrante.
func filterSorted(books []Book, status Status) []Book {
	var out []Book
	for _, bk := range books {
		if bk.Status == status {
			out = append(out, bk)
		}
	}
	// Filter ya ordena cuando la entrada viene del shelf; mantener el orden
	// por título aquí hace RenderHTML independiente del orden de entrada.
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].Title < out[j-1].Title; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// ServeHandler envuelve RenderHTML en un handler de SOLO LECTURA (ADR 0003):
// recarga el ledger en cada request y nunca escribe.
func ServeHandler(path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "solo lectura: método no permitido", http.StatusMethodNotAllowed)
			return
		}
		shelf, err := LoadShelf(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(RenderHTML(shelf.Books, time.Now()))
	})
}
