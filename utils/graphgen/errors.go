package graphgen

import "errors"

var (
	// ErrPointOutsideContributionGraph is an occurring error when a wanted point is outside the contributions graph.
	ErrPointOutsideContributionGraph = errors.New("point is outside the contribution graph")
	// ErrContributionGraphGlyphOverflow is an occurring when the given glyph is drawn outside the contribution graph,
	// ie starts inside of it and ends outside it.
	// or that the glyph draws on a nil cell.
	ErrContributionGraphGlyphOverflow = errors.New("glyph overflows the contribution graph at this point")
)
