package graphgen

import (
	"bytes"
	"html/template"
	"io"
	"slices"
	"strings"
	ttemplate "text/template"
)

type ContributionsGraphGenerator interface {
	SetFont(gm Font)
	GetFinalForm(text string) (io.Reader, error)
}

type GeneratorType int

const (
	HtmlGeneratorType GeneratorType = iota
	CheatScriptGeneratorType
)

func NewContributionsGraphGenerator(t GeneratorType, cg *ContributionsGraph) ContributionsGraphGenerator {
	switch t {
	case HtmlGeneratorType:
		return &HtmlContributionsGraphGenerator{
			cg:   cg,
			font: font3x5Glyphs,
			cellFiller: map[CellType]string{
				NilCell:      "nilCell",
				EmptyCell:    "emptyCell",
				OccupiedCell: "occupiedCell",
			},
		}
	case CheatScriptGeneratorType:
		return &CheatScriptContributionsGraphGenerator{
			cg:   cg,
			font: font3x5Glyphs,
		}
	default:
		return nil
	}
}

type HtmlContributionsGraphGenerator struct {
	cg         *ContributionsGraph
	font       Font
	cellFiller map[CellType]string
}

func (h *HtmlContributionsGraphGenerator) SetFont(font Font) {
	h.font = font
}

func (h *HtmlContributionsGraphGenerator) GetFinalForm(text string) (io.Reader, error) {
	// get a usable sentence
	sentence := h.font.TextToGlyphs(text)
	err := h.cg.DrawSentence(sentence, Point{0, 1})
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

type CheatScriptContributionsGraphGenerator struct {
	cg   *ContributionsGraph
	font Font
}

func (c *CheatScriptContributionsGraphGenerator) SetFont(font Font) {
	c.font = font
}

func (c *CheatScriptContributionsGraphGenerator) GetFinalForm(text string) (io.Reader, error) {
	// get a usable sentence
	sentence := c.font.TextToGlyphs(text)
	err := c.cg.DrawSentence(sentence, Point{0, 1})
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
