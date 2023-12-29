package graphgen

import "errors"

var (
	ErrPointOutsideContributionGraph  = errors.New("point is outside the contribution graph")
	ErrContributionGraphGlyphOverflow = errors.New("glyph overflows the contribution graph at this point")
)
