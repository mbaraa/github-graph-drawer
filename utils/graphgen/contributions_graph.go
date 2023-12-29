package graphgen

import (
	"slices"
	"time"
	"unicode"
)

type (
	GitDate       string
	CellType      int
	Glyph         [][]CellType
	GlyphSentence []Glyph
	Font          map[byte]Glyph
)

const (
	NilCell CellType = iota
	EmptyCell
	OccupiedCell
)

type Cell struct {
	Type CellType
	Date GitDate
}

type Point struct {
	X int
	Y int
}

type ContributionsGraph struct {
	cells [7][53]Cell
}

func (c *ContributionsGraph) Cells() [7][53]Cell {
	return c.cells
}

func (c ContributionsGraph) Init(year int) *ContributionsGraph {
	for i := 0; i < len(c.cells); i++ {
		for j := 0; j < len(c.cells[0]); j++ {
			date := ""

			c.cells[i][j] = Cell{
				Type: NilCell,
				Date: GitDate(date),
			}
		}
	}

	yearTime := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	firstWeekdayOfYear := yearTime.Weekday()
	weekdays := []time.Weekday{time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday}
	weekdayIndex := slices.Index(weekdays, firstWeekdayOfYear)

	// initializing cells before the first day
	for i := 0; i < weekdayIndex; i++ {
		c.cells[0][i] = Cell{
			Type: NilCell,
			Date: "",
		}
	}

	// initializing cells after the first day
	for i, x, y := weekdayIndex, 0, weekdayIndex; i < len(c.cells)*len(c.cells[0]); i++ {
		cell := NilCell
		if yearTime.Year() == year {
			cell = EmptyCell
		}
		c.cells[y][x] = Cell{
			Type: cell,
			Date: GitDate(yearTime.Format("2006-01-02T15:04:05")),
		}
		if (i+1)%7 == 0 {
			x++
		}
		y = (y + 1) % 7
	}

	return &c
}

func (c *ContributionsGraph) DrawGlyph(g Glyph, start Point) error {
	if start.X > len(c.cells[0]) || start.X < 0 ||
		start.Y > len(c.cells) || start.Y < 0 {
		return ErrPointOutsideContributionGraph
	}
	if start.X+len(g[0]) > len(c.cells[0]) || start.Y+len(g) > len(c.cells) {
		return ErrContributionGraphGlyphOverflow
	}

	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[0]); x++ {
			c.cells[start.Y+y][start.X+x].Type = g[y][x]
		}
	}
	return nil
}

func (c *ContributionsGraph) DrawSentence(gs GlyphSentence, start Point) error {
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

func (c *ContributionsGraph) Reset() {
	for y := 0; y < len(c.cells); y++ {
		for x := 0; x < len(c.cells[0]); x++ {
			if c.cells[y][x].Type == OccupiedCell {
				c.cells[y][x].Type = EmptyCell
			}
		}
	}
}

func (font Font) TextToGlyphs(s string) GlyphSentence {
	glyphs := make(GlyphSentence, len(s))
	for i, chr := range s {
		glyphs[i] = font[byte(unicode.ToUpper(chr))]
	}
	return glyphs
}
