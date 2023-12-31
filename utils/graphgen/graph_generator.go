package graphgen

import (
	"bytes"
	"html/template"
	"io"
	"slices"
	"strings"
	ttemplate "text/template"
)

// ContributionsGraphGenerator is an interface that represents a general contributions graph generator.
type ContributionsGraphGenerator interface {
	// SetFont sets the used font in the generation process,
	// available fonts:
	// - Font3x3; holds 3x3 glyphs.
	// - Font3x5; holds 3x5 glyphs.
	SetFont(gm Font)
	// GetFinalForm receives a string, and returns an io.Reader, with the resulting graph,
	// and an occurring errpr,
	// the reason behind using io.Reader, is that the output is 99.9% of time is a file,
	// for now the current implementations are files, so returning an io.Reader, makes more sense
	// than a string, which will be converted to an io.Reader, or os.File, so this saves conversion time,
	// and it looks a lot more neater.
	GetFinalForm(text string) (io.Reader, error)
}

// GeneratorType is a selector between the current implementations of ContributionsGraphGenerator,
// to be used with the NewContributionsGraphGenerator factory function.
type GeneratorType int

const (
	// HtmlGeneratorType makes NewContributionsGraphGenerator return an HtmlContributionsGraphGenerator
	HtmlGeneratorType GeneratorType = iota
	// CheatScriptGeneratorType makes NewContributionsGraphGenerator return an CheatScriptContributionsGraphGenerator
	CheatScriptGeneratorType
)

// NewContributionsGraphGenerator is a factory method that returns a new instance implemnting ContributionsGraphGenerator,
// based on the given GeneratorType.
func NewContributionsGraphGenerator(t GeneratorType, cg *ContributionsGraph) ContributionsGraphGenerator {
	switch t {
	case HtmlGeneratorType:
		return &htmlContributionsGraphGenerator{
			cg:   cg,
			font: Font3x5,
			cellFiller: map[CellType]string{
				NilCell:      "nilCell",
				EmptyCell:    "emptyCell",
				OccupiedCell: "occupiedCell",
			},
		}
	case CheatScriptGeneratorType:
		return &cheatScriptContributionsGraphGenerator{
			cg:   cg,
			font: Font3x5,
		}
	default:
		return nil
	}
}

type htmlContributionsGraphGenerator struct {
	cg         *ContributionsGraph
	font       Font
	cellFiller map[CellType]string
}

func (h *htmlContributionsGraphGenerator) SetFont(font Font) {
	h.font = font
}

func (h *htmlContributionsGraphGenerator) GetFinalForm(text string) (io.Reader, error) {
	// get a usable sentence
	sentence := h.font.TextToGlyphs(text)
	// TODO:
	// handle different fonts' sizes.
	err := h.cg.DrawSentence(sentence, Point{0, 0})
	if err != nil {
		return nil, err
	}

	// map css classes
	classes := make([][]string, len(h.cg.Cells()))
	for i, row := range h.cg.Cells() {
		classes[i] = make([]string, len(row))
		for j, cell := range row {
			classes[i][j] = h.cellFiller[cell.Type]
		}
	}

	// FIX:
	// fix this illegal floating template instance.
	tmpl := template.Must(template.ParseGlob("./templates/html/*"))
	buf := bytes.NewBuffer([]byte{})
	err = tmpl.ExecuteTemplate(buf, "center_piece", map[string]any{
		"Cells": classes,
		"Msg":   text,
	})
	if err != nil {
		return nil, err
	}

	return buf, nil
}

type cheatScriptContributionsGraphGenerator struct {
	cg   *ContributionsGraph
	font Font
}

func (c *cheatScriptContributionsGraphGenerator) SetFont(font Font) {
	c.font = font
}

func (c *cheatScriptContributionsGraphGenerator) GetFinalForm(text string) (io.Reader, error) {
	// get a usable sentence
	sentence := c.font.TextToGlyphs(text)
	// TODO:
	// handle different fonts' sizes.
	err := c.cg.DrawSentence(sentence, Point{0, 0})
	if err != nil {
		return nil, err
	}

	gitDates := make([]string, 0)
	for _, day := range c.cg.Cells() {
		for _, weekDay := range day {
			if weekDay.Type == NilCell || weekDay.Type == EmptyCell {
				continue
			}
			gitDates = append(gitDates, string(weekDay.Date))
		}
	}
	slices.Sort(gitDates)

	// FIX:
	// fix this illegal floating template instance.
	tmpl := ttemplate.Must(ttemplate.ParseFiles("./templates/text/generate-commits.sh"))
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer([]byte{})
	err = tmpl.ExecuteTemplate(buf, "generate_commits_script", map[string]string{
		"Dates": strings.Join(gitDates, " "),
	})
	if err != nil {
		return nil, err
	}

	return buf, nil
}
