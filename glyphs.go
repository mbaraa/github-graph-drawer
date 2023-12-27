package main

var (
	font3x3Glyphs = GlyphMapper{
		'A': {
			{EmptyCell, Occupied4, EmptyCell},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
		},
		'B': {
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, Occupied4, Occupied4},
		},
		'C': {
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, EmptyCell, EmptyCell},
			{Occupied4, Occupied4, Occupied4},
		},
		'D': {
			{Occupied4, Occupied4, EmptyCell},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, EmptyCell},
		},
		'E': {
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, Occupied4, EmptyCell},
			{Occupied4, Occupied4, Occupied4},
		},
		'F': {
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, Occupied4, EmptyCell},
			{Occupied4, EmptyCell, EmptyCell},
		},
		'G': {
			{Occupied4, Occupied4, EmptyCell},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, Occupied4, Occupied4},
		},
		'H': {
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
		},
		'I': {
			{EmptyCell, Occupied4, EmptyCell},
			{EmptyCell, Occupied4, EmptyCell},
			{EmptyCell, Occupied4, EmptyCell},
		},
		'J': {
			{EmptyCell, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, Occupied4},
		},
		'K': {
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, EmptyCell},
			{Occupied4, EmptyCell, Occupied4},
		},
		'L': {
			{Occupied4, EmptyCell, EmptyCell},
			{Occupied4, EmptyCell, EmptyCell},
			{Occupied4, Occupied4, Occupied4},
		},
		'M': {
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
		},
		'N': {
			{EmptyCell, EmptyCell, Occupied4},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, EmptyCell, EmptyCell},
		},
		'O': {
			{EmptyCell, Occupied4, EmptyCell},
			{Occupied4, EmptyCell, Occupied4},
			{EmptyCell, Occupied4, EmptyCell},
		},
		'P': {
			{Occupied4, Occupied4, EmptyCell},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, EmptyCell, EmptyCell},
		},
		'Q': {
			{EmptyCell, Occupied4, EmptyCell},
			{Occupied4, EmptyCell, Occupied4},
			{EmptyCell, Occupied4, Occupied4},
		},
		'R': {
			{Occupied4, Occupied4, EmptyCell},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
		},
		'S': {
			{EmptyCell, Occupied4, Occupied4},
			{EmptyCell, Occupied4, EmptyCell},
			{Occupied4, Occupied4, EmptyCell},
		},
		'T': {
			{Occupied4, Occupied4, Occupied4},
			{EmptyCell, Occupied4, EmptyCell},
			{EmptyCell, Occupied4, EmptyCell},
		},
		'U': {
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, Occupied4},
		},
		'V': {
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{EmptyCell, Occupied4, EmptyCell},
		},
		'W': {
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, Occupied4, Occupied4},
		},
		'X': {
			{Occupied4, EmptyCell, Occupied4},
			{EmptyCell, Occupied4, EmptyCell},
			{Occupied4, EmptyCell, Occupied4},
		},
		'Y': {
			{Occupied4, EmptyCell, Occupied4},
			{EmptyCell, Occupied4, EmptyCell},
			{EmptyCell, Occupied4, EmptyCell},
		},
		'Z': {
			{Occupied4, Occupied4, EmptyCell},
			{EmptyCell, Occupied4, EmptyCell},
			{EmptyCell, Occupied4, Occupied4},
		},
		' ': {
			{EmptyCell},
			{EmptyCell},
			{EmptyCell},
		},
	}

	font3x5Glyphs = GlyphMapper{
		'A': {
			{EmptyCell, Occupied4, EmptyCell},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
		},
		'B': {
			{Occupied4, Occupied4, EmptyCell},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, EmptyCell},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, EmptyCell},
		},
		'C': {
			{EmptyCell, Occupied4, Occupied4},
			{Occupied4, EmptyCell, EmptyCell},
			{Occupied4, EmptyCell, EmptyCell},
			{Occupied4, EmptyCell, EmptyCell},
			{EmptyCell, Occupied4, Occupied4},
		},
		'D': {
			{Occupied4, Occupied4, EmptyCell},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, EmptyCell},
		},
		'E': {
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, EmptyCell, EmptyCell},
			{Occupied4, Occupied4, EmptyCell},
			{Occupied4, EmptyCell, EmptyCell},
			{Occupied4, Occupied4, Occupied4},
		},
		'F': {
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, EmptyCell, EmptyCell},
			{Occupied4, Occupied4, EmptyCell},
			{Occupied4, EmptyCell, EmptyCell},
			{Occupied4, EmptyCell, EmptyCell},
		},
		'G': {
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, EmptyCell, EmptyCell},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, Occupied4},
		},
		'H': {
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
		},
		'I': {
			{Occupied4, Occupied4, Occupied4},
			{EmptyCell, Occupied4, EmptyCell},
			{EmptyCell, Occupied4, EmptyCell},
			{EmptyCell, Occupied4, EmptyCell},
			{Occupied4, Occupied4, Occupied4},
		},
		'J': {
			{EmptyCell, Occupied4, Occupied4},
			{EmptyCell, EmptyCell, Occupied4},
			{EmptyCell, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, Occupied4},
		},
		'K': {
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, EmptyCell},
			{Occupied4, EmptyCell, EmptyCell},
			{Occupied4, Occupied4, EmptyCell},
			{Occupied4, EmptyCell, Occupied4},
		},
		'L': {
			{Occupied4, EmptyCell, EmptyCell},
			{Occupied4, EmptyCell, EmptyCell},
			{Occupied4, EmptyCell, EmptyCell},
			{Occupied4, EmptyCell, EmptyCell},
			{Occupied4, Occupied4, Occupied4},
		},
		'M': {
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
		},
		'N': {
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
		},
		'O': {
			{EmptyCell, Occupied4, EmptyCell},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{EmptyCell, Occupied4, EmptyCell},
		},
		'P': {
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, EmptyCell, EmptyCell},
			{Occupied4, EmptyCell, EmptyCell},
		},
		'Q': {
			{EmptyCell, EmptyCell, EmptyCell},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, Occupied4},
			{EmptyCell, EmptyCell, Occupied4},
		},
		'R': {
			{Occupied4, Occupied4, EmptyCell},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, EmptyCell},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
		},
		'S': {
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, EmptyCell, EmptyCell},
			{EmptyCell, Occupied4, EmptyCell},
			{EmptyCell, EmptyCell, Occupied4},
			{Occupied4, Occupied4, Occupied4},
		},
		'T': {
			{Occupied4, Occupied4, Occupied4},
			{EmptyCell, Occupied4, EmptyCell},
			{EmptyCell, Occupied4, EmptyCell},
			{EmptyCell, Occupied4, EmptyCell},
			{EmptyCell, Occupied4, EmptyCell},
		},
		'U': {
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, Occupied4},
		},
		'V': {
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{EmptyCell, Occupied4, EmptyCell},
			{EmptyCell, Occupied4, EmptyCell},
		},
		'W': {
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, Occupied4, Occupied4},
			{Occupied4, Occupied4, Occupied4},
		},
		'X': {
			{Occupied4, EmptyCell, Occupied4},
			{EmptyCell, Occupied4, EmptyCell},
			{EmptyCell, Occupied4, EmptyCell},
			{EmptyCell, Occupied4, EmptyCell},
			{Occupied4, EmptyCell, Occupied4},
		},
		'Y': {
			{Occupied4, EmptyCell, Occupied4},
			{Occupied4, EmptyCell, Occupied4},
			{EmptyCell, Occupied4, EmptyCell},
			{EmptyCell, Occupied4, EmptyCell},
			{EmptyCell, Occupied4, EmptyCell},
		},
		'Z': {
			{Occupied4, Occupied4, Occupied4},
			{EmptyCell, EmptyCell, Occupied4},
			{EmptyCell, Occupied4, EmptyCell},
			{Occupied4, EmptyCell, EmptyCell},
			{Occupied4, Occupied4, Occupied4},
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
