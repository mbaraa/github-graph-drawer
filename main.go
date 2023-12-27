package main

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"text/template"
	"time"
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

type ContributionGraph struct {
	Cells [7][52]Cell
}

func (c ContributionGraph) New() *ContributionGraph {
	for i := 0; i < len(c.Cells); i++ {
		for j := 0; j < len(c.Cells[0]); j++ {
			c.Cells[i][j] = EmptyCell
		}
	}
	return &c
}

type Point struct {
	X int
	Y int
}

func (c *ContributionGraph) DrawGlyph(g Glyph, start Point) error {
	if start.X > len(c.Cells[0]) || start.X < 0 ||
		start.Y > len(c.Cells) || start.Y < 0 {
		return ErrPointOutsideContributionGraph
	}
	if start.X+len(g[0]) > len(c.Cells[0]) || start.Y+len(g) > len(c.Cells) {
		return ErrContributionGraphGlyphOverflow
	}

	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[0]); x++ {
			c.Cells[start.Y+y][start.X+x] = g[y][x]
		}
	}
	return nil
}

func (c *ContributionGraph) DrawSentence(gs GlyphSentence, start Point) error {
	for _, glyph := range gs {
		err := c.DrawGlyph(glyph, start)
		if err != nil {
			return err
		}
		if len(glyph[0]) >= 3 {
			start.X += len(glyph[0]) + 1
		} else {
			start.X += len(glyph[0])
		}
	}
	return nil
}

type Glyph [][]Cell

type GlyphSentence []Glyph

type GlyphMapper map[byte]Glyph

func (font GlyphMapper) TextToGlyphs(s string) GlyphSentence {
	glyphs := make(GlyphSentence, len(s))
	for i, chr := range s {
		glyphs[i] = font[byte(unicode.ToUpper(chr))]
	}
	return glyphs
}

func main() {
	// server()
	script()
}

func script() {
	sentence := font3x5Glyphs.TextToGlyphs("Hola Soy Baraa")
	cg := ContributionGraph{}.New()
	err := cg.DrawSentence(sentence, Point{0, 0})
	if err != nil {
		panic(err)
	}

	days := getContributionsForYear(2023)
	for k, v := range days {
		if k.X == 1 {
			fmt.Println(k, v)
		}
	}
	for _, day := range cg.Cells {
		for _, weekDay := range day {
			fmt.Print(weekDay)
		}
		fmt.Println()
	}

	gitDates := make([]string, 0)
	for y, day := range cg.Cells {
		for x, weekDay := range day {
			if weekDay == NilCell || weekDay == EmptyCell {
				continue
			}
			gitDates = append(gitDates, string(days[Point{x, y}]))
		}
	}

	slices.Sort(gitDates)

	fmt.Println(textFS.ReadDir("."))
	textTemplates, err := template.ParseFS(textFS, "templates/text/generate-commits.sh")
	if err != nil {
		panic(err)
	}

	out, _ := os.Create("out.sh")
	defer out.Close()
	textTemplates.Execute(out, map[string]string{
		"Dates": strings.Join(gitDates, " "),
	})
}

var (
	//go:embed templates/html
	res embed.FS

	pages = map[string]string{
		"/": "templates/html/index.html",
	}
)

var cellFiller = map[Cell]string{
	NilCell:   "nilCell",
	EmptyCell: "emptyCell",
	Occupied1: "occupied1",
	Occupied2: "occupied2",
	Occupied3: "occupied3",
	Occupied4: "occupied4",
}

func mapCellsToCSS(cg ContributionGraph) [][]string {
	classes := make([][]string, len(cg.Cells))
	for i, row := range cg.Cells {
		classes[i] = make([]string, len(row))
		for j, cell := range row {
			classes[i][j] = cellFiller[cell]
		}
	}
	return classes
}

type GitDate string

func getContributionsForYear(year int) map[Point]GitDate {
	yearTime := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	numDays := 365
	if yearTime.Year()%4 == 0 && yearTime.Year()%100 != 0 {
		numDays++
	}

	firstWeekdayOfYear := yearTime.Weekday()
	weekdays := []time.Weekday{time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday}
	weekdayIndex := slices.Index(weekdays, firstWeekdayOfYear)

	days := make(map[Point]GitDate)
	point := Point{0, weekdayIndex}

	for i := 0; i < numDays; i++ {
		days[point] = GitDate(yearTime.Format("2006-01-02T15:04:05"))
		yearTime = yearTime.Add(time.Hour*24 + time.Second)
		if (i+1)%7 == 0 {
			point.X++
		}
		point.Y = (point.Y + 1) % 7
	}

	return days
}

func server() {
	tt := template.Must(template.ParseGlob("./templates/html/*"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		sentence := font3x5Glyphs.TextToGlyphs("Hola Soy Baraa")
		cg := ContributionGraph{}.New()
		err := cg.DrawSentence(sentence, Point{0, 0})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := tt.ExecuteTemplate(w, "index", map[string]any{
			"Cells": mapCellsToCSS(*cg),
		}); err != nil {
			fmt.Println(err)
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
