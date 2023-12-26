package main

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"unicode"
)

var (
	//go:embed templates/text
	textFS embed.FS

	//go:embed templates/html
	htmlFS embed.FS
)

var (
	ErrPointOutsideContributionGraph  = errors.New("point is outside the contribution graph")
	ErrContributionGraphGlyphOverflow = errors.New("glyph overflows the contribution graph at this point")
)

type Cell int

const (
	NilCell Cell = iota
	EmptyCell
	Occupied1
	Occupied2
	Occupied3
	Occupied4
)

type ContributionGraph [7][52]Cell

type Point struct {
	X int
	Y int
}

type Glyph [][]Cell

func (f Glyph) Draw(cg *ContributionGraph, start Point) error {
	if start.X > len(cg[0]) || start.X < 0 ||
		start.Y > len(cg) || start.Y < 0 {
		return ErrPointOutsideContributionGraph
	}
	if start.X+len(f[0]) > len(cg[0]) || start.Y+len(f) > len(cg) {
		return ErrContributionGraphGlyphOverflow
	}

	for y := 0; y < len(f); y++ {
		for x := 0; x < len(f[0]); x++ {
			cg[start.Y+y][start.X+x] = f[y][x]
		}
	}

	return nil
}

type GlyphMapper map[byte]Glyph

func (font GlyphMapper) textToGlyphs(s string) []Glyph {
	glyphs := make([]Glyph, len(s))
	for i, chr := range s {
		glyphs[i] = font[byte(unicode.ToUpper(chr))]
	}
	return glyphs
}

var (
	font3x3Glyphs = GlyphMapper{
		'A': {
			{EmptyCell, Occupied1, EmptyCell},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
		},
		'B': {
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, Occupied1, Occupied1},
		},
		'C': {
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, EmptyCell, EmptyCell},
			{Occupied1, Occupied1, Occupied1},
		},
		'D': {
			{Occupied1, Occupied1, EmptyCell},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, EmptyCell},
		},
		'E': {
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, Occupied1, EmptyCell},
			{Occupied1, Occupied1, Occupied1},
		},
		'F': {
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, Occupied1, EmptyCell},
			{Occupied1, EmptyCell, EmptyCell},
		},
		'G': {
			{Occupied1, Occupied1, EmptyCell},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, Occupied1, Occupied1},
		},
		'H': {
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
		},
		'I': {
			{EmptyCell, Occupied1, EmptyCell},
			{EmptyCell, Occupied1, EmptyCell},
			{EmptyCell, Occupied1, EmptyCell},
		},
		'J': {
			{EmptyCell, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, Occupied1},
		},
		'K': {
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, EmptyCell},
			{Occupied1, EmptyCell, Occupied1},
		},
		'L': {
			{Occupied1, EmptyCell, EmptyCell},
			{Occupied1, EmptyCell, EmptyCell},
			{Occupied1, Occupied1, Occupied1},
		},
		'M': {
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
		},
		'N': {
			{EmptyCell, EmptyCell, Occupied1},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, EmptyCell, EmptyCell},
		},
		'O': {
			{EmptyCell, Occupied1, EmptyCell},
			{Occupied1, EmptyCell, Occupied1},
			{EmptyCell, Occupied1, EmptyCell},
		},
		'P': {
			{Occupied1, Occupied1, EmptyCell},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, EmptyCell, EmptyCell},
		},
		'Q': {
			{EmptyCell, Occupied1, EmptyCell},
			{Occupied1, EmptyCell, Occupied1},
			{EmptyCell, Occupied1, Occupied1},
		},
		'R': {
			{Occupied1, Occupied1, EmptyCell},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
		},
		'S': {
			{EmptyCell, Occupied1, Occupied1},
			{EmptyCell, Occupied1, EmptyCell},
			{Occupied1, Occupied1, EmptyCell},
		},
		'T': {
			{Occupied1, Occupied1, Occupied1},
			{EmptyCell, Occupied1, EmptyCell},
			{EmptyCell, Occupied1, EmptyCell},
		},
		'U': {
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, Occupied1},
		},
		'V': {
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{EmptyCell, Occupied1, EmptyCell},
		},
		'W': {
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, Occupied1, Occupied1},
		},
		'X': {
			{Occupied1, EmptyCell, Occupied1},
			{EmptyCell, Occupied1, EmptyCell},
			{Occupied1, EmptyCell, Occupied1},
		},
		'Y': {
			{Occupied1, EmptyCell, Occupied1},
			{EmptyCell, Occupied1, EmptyCell},
			{EmptyCell, Occupied1, EmptyCell},
		},
		'Z': {
			{Occupied1, Occupied1, EmptyCell},
			{EmptyCell, Occupied1, EmptyCell},
			{EmptyCell, Occupied1, Occupied1},
		},
		' ': {
			{EmptyCell},
			{EmptyCell},
			{EmptyCell},
		},
	}

	font3x5Glyphs = GlyphMapper{
		'A': {
			{EmptyCell, Occupied1, EmptyCell},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
		},
		'B': {
			{Occupied1, Occupied1, EmptyCell},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, EmptyCell},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, EmptyCell},
		},
		'C': {
			{EmptyCell, Occupied1, Occupied1},
			{Occupied1, EmptyCell, EmptyCell},
			{Occupied1, EmptyCell, EmptyCell},
			{Occupied1, EmptyCell, EmptyCell},
			{EmptyCell, Occupied1, Occupied1},
		},
		'D': {
			{Occupied1, Occupied1, EmptyCell},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, EmptyCell},
		},
		'E': {
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, EmptyCell, EmptyCell},
			{Occupied1, Occupied1, EmptyCell},
			{Occupied1, EmptyCell, EmptyCell},
			{Occupied1, Occupied1, Occupied1},
		},
		'F': {
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, EmptyCell, EmptyCell},
			{Occupied1, Occupied1, EmptyCell},
			{Occupied1, EmptyCell, EmptyCell},
			{Occupied1, EmptyCell, EmptyCell},
		},
		'G': {
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, EmptyCell, EmptyCell},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, Occupied1},
		},
		'H': {
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
		},
		'I': {
			{Occupied1, Occupied1, Occupied1},
			{EmptyCell, Occupied1, EmptyCell},
			{EmptyCell, Occupied1, EmptyCell},
			{EmptyCell, Occupied1, EmptyCell},
			{Occupied1, Occupied1, Occupied1},
		},
		'J': {
			{EmptyCell, Occupied1, Occupied1},
			{EmptyCell, EmptyCell, Occupied1},
			{EmptyCell, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, Occupied1},
		},
		'K': {
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, EmptyCell},
			{Occupied1, EmptyCell, EmptyCell},
			{Occupied1, Occupied1, EmptyCell},
			{Occupied1, EmptyCell, Occupied1},
		},
		'L': {
			{Occupied1, EmptyCell, EmptyCell},
			{Occupied1, EmptyCell, EmptyCell},
			{Occupied1, EmptyCell, EmptyCell},
			{Occupied1, EmptyCell, EmptyCell},
			{Occupied1, Occupied1, Occupied1},
		},
		'M': {
			{EmptyCell, EmptyCell, EmptyCell},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
		},
		'N': {
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
		},
		'O': {
			{EmptyCell, Occupied1, EmptyCell},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{EmptyCell, Occupied1, EmptyCell},
		},
		'P': {
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, EmptyCell, EmptyCell},
			{Occupied1, EmptyCell, EmptyCell},
		},
		'Q': {
			{EmptyCell, EmptyCell, EmptyCell},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, Occupied1},
			{EmptyCell, EmptyCell, Occupied1},
		},
		'R': {
			{Occupied1, Occupied1, EmptyCell},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, EmptyCell},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
		},
		'S': {
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, EmptyCell, EmptyCell},
			{EmptyCell, Occupied1, EmptyCell},
			{EmptyCell, EmptyCell, Occupied1},
			{Occupied1, Occupied1, Occupied1},
		},
		'T': {
			{Occupied1, Occupied1, Occupied1},
			{EmptyCell, Occupied1, EmptyCell},
			{EmptyCell, Occupied1, EmptyCell},
			{EmptyCell, Occupied1, EmptyCell},
			{EmptyCell, Occupied1, EmptyCell},
		},
		'U': {
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, Occupied1},
		},
		'V': {
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{EmptyCell, Occupied1, EmptyCell},
			{EmptyCell, Occupied1, EmptyCell},
		},
		'W': {
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, Occupied1, Occupied1},
			{Occupied1, Occupied1, Occupied1},
		},
		'X': {
			{Occupied1, EmptyCell, Occupied1},
			{EmptyCell, Occupied1, EmptyCell},
			{EmptyCell, Occupied1, EmptyCell},
			{EmptyCell, Occupied1, EmptyCell},
			{Occupied1, EmptyCell, Occupied1},
		},
		'Y': {
			{Occupied1, EmptyCell, Occupied1},
			{Occupied1, EmptyCell, Occupied1},
			{EmptyCell, Occupied1, EmptyCell},
			{EmptyCell, Occupied1, EmptyCell},
			{EmptyCell, Occupied1, EmptyCell},
		},
		'Z': {
			{Occupied1, Occupied1, Occupied1},
			{EmptyCell, EmptyCell, Occupied1},
			{EmptyCell, Occupied1, EmptyCell},
			{Occupied1, EmptyCell, EmptyCell},
			{Occupied1, Occupied1, Occupied1},
		},
		' ': {
			{EmptyCell},
			{EmptyCell},
			{EmptyCell},
			{EmptyCell},
			{EmptyCell},
		},
	}
)

func main() {
	cg := ContributionGraph{}
	glyphs := font3x3Glyphs.textToGlyphs("Hello World")
	point := Point{1, 1}
	for _, glyph := range glyphs {
		err := glyph.Draw(&cg, point)
		if err != nil {
			fmt.Println(err)
		}
		point.X += len(glyph[0]) + 1
	}

	for _, day := range cg {
		for _, week := range day {
			if week != Occupied1 {
				fmt.Print(" ")
				continue
			}
			fmt.Print(week)
		}
		fmt.Println()
	}

	return

	fmt.Println(textFS.ReadDir("."))
	textTemplates, err := template.ParseFS(textFS, "templates/text/generate-commits.sh")
	if err != nil {
		panic(err)
	}

	out, _ := os.Create("out.sh")
	defer out.Close()
	textTemplates.Execute(out, map[string]string{
		"Dates": "1 2 2 3 4",
	})
}

var (
	//go:embed templates/html
	res embed.FS

	pages = map[string]string{
		"/": "templates/html/index.html",
	}
)

func main2() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page, ok := pages[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		tpl, err := template.ParseFS(res, page)
		if err != nil {
			log.Printf("page %s not found in pages cache...", r.RequestURI)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		data := map[string]interface{}{
			"userAgent": r.UserAgent(),
		}
		if err := tpl.Execute(w, data); err != nil {
			return
		}
	})

	http.FileServer(http.FS(res))

	log.Println("server started...")
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		panic(err)
	}
}
