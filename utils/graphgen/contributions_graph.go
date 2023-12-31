package graphgen

import (
	"slices"
	"time"
	"unicode"
)

type (
	// GitDate represents a git date used in `git commit --date`
	GitDate string
	// CellType represents a cell's filler.
	CellType int
	// Glyph is how a character is shown on the contributions graph.
	Glyph [][]CellType
	// GlyphSentence well, it's the string of glyphs.
	GlyphSentence []Glyph
	// Font the byte to glyph mapper.
	Font map[byte]Glyph
)

const (
	// NilCell represents a cell that's outside of the selected year.
	NilCell CellType = iota
	// EmptyCell represents a cell that's inside of the selected year,
	// but has no contributions in it.
	EmptyCell
	// OccupiedCell represents a cells that's inside of the selected year,
	// and has some juicy contributions in it.
	OccupiedCell
)

// Cell is the single element in the contributions graph's grid,
// where it contains the cell's filler and the associated git date.
type Cell struct {
	Type CellType
	Date GitDate
}

// Point is a 2D point.
type Point struct {
	X int
	Y int
}

// ContributionsGraph is a representational state of the contributions graph.
type ContributionsGraph struct {
	cells [7][53]Cell
}

// Cells returns the contributions graph's cells.
func (c *ContributionsGraph) Cells() [7][53]Cell {
	return c.cells
}

// Init initializes the contributions graph's cells for the given year,
// with empty cells, and nil cells for cells outside the given year,
// and sets each cell's corresponding git date.
func (c ContributionsGraph) Init(year int) *ContributionsGraph {
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
		yearTime = yearTime.Add(time.Hour*24 + time.Second)
		if (i+1)%7 == 0 {
			x++
		}
		y = (y + 1) % 7
	}

	return &c
}

// DrawGlyph draws the given glyps onto the contributions graph.
// and returns an occurring error, on of:
// - ErrContributionGraphGlyphOverflow
// - ErrPointOutsideContributionGraph
func (c *ContributionsGraph) DrawGlyph(g Glyph, start *Point) error {
	firstWeekAvailable, lastWeekAvailable := true, true
	for i := 0; i < len(c.cells); i++ {
		if c.cells[i][0].Type == NilCell {
			firstWeekAvailable = false
		}
		if c.cells[i][len(c.cells[0])-1].Type == NilCell {
			lastWeekAvailable = false
		}
	}

	if start.X > len(c.cells[0]) || start.X < 0 ||
		start.Y > len(c.cells) || start.Y < 0 {
		return ErrPointOutsideContributionGraph
	}
	if (!lastWeekAvailable && start.X+len(g[0]) >= len(c.cells[0])-1) ||
		start.X+len(g[0]) > len(c.cells[0]) ||
		start.Y+len(g) > len(c.cells) {
		return ErrContributionGraphGlyphOverflow
	}
	if !firstWeekAvailable && start.X <= 0 {
		start.X++
	}

	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[0]); x++ {
			c.cells[start.Y+y][start.X+x].Type = g[y][x]
		}
	}
	return nil
}

// DrawSentence calls DrawGlyph to draw a complete sentence on the contributions graph,
// and if the font is small enough it jumps to a new line drawing rest of the sentence.
func (c *ContributionsGraph) DrawSentence(gs GlyphSentence, start Point) error {
	for i, glyph := range gs {
		err := c.DrawGlyph(glyph, &start)
		if err != nil {
			return err
		}
		// HACK:
		// handle whitespace to not take more than what it needs.
		if len(glyph[0]) >= 3 {
			start.X += len(glyph[0]) + 1
		} else {
			start.X += len(glyph[0])
		}

		// move writing cursor a line down.
		if i <= len(gs)-2 &&
			start.X+len(gs[i+1][0]) >= len(c.cells[0])-1 {
			start.Y += len(glyph) + 1
			start.X = 0
		}
	}
	return nil
}

// Reset sets each occupied cell to empty cell,
// allowing the contributions graph to be used more than one time.
func (c *ContributionsGraph) Reset() {
	for y := 0; y < len(c.cells); y++ {
		for x := 0; x < len(c.cells[0]); x++ {
			if c.cells[y][x].Type == OccupiedCell {
				c.cells[y][x].Type = EmptyCell
			}
		}
	}
}

// TextToGlyphs converts the given string to a drawable glyph sentence.
func (font Font) TextToGlyphs(s string) GlyphSentence {
	glyphs := make(GlyphSentence, 0)
	for i := 0; i < len(s); i++ {
		glyph, exists := font[byte(unicode.ToUpper(rune(s[i])))]
		if !exists {
			continue
		}
		glyphs = append(glyphs, glyph)
	}
	return glyphs
}
