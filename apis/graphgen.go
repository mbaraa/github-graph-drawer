package apis

import (
	"github-graph-drawer/log"
	"github-graph-drawer/utils/graphgen"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"
)

type graphGeneratorApi struct {
	templates *template.Template
}

func NewGraphGeneratorApi() IHandler {
	g := &graphGeneratorApi{}
	endpoints := Endpoints{
		"GET preview": g.handlePreviewContributionsGraph,
		"GET script":  g.handleGenerateCheatScript,
	}
	return NewHandler(endpoints, "/graphgen/")
}

func (g *graphGeneratorApi) handlePreviewContributionsGraph(w http.ResponseWriter, r *http.Request) {
	msg, exists := r.URL.Query()["msg"]
	if !exists {
		log.Warningln("someone sent a bad request...")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	font := r.URL.Query().Get("font")
	year := time.Now().Year()
	if parsedYear, err := strconv.Atoi(r.URL.Query().Get("year")); err == nil {
		year = parsedYear
	}
	commitsCount := 80
	if parsedCommitsCount, err := strconv.Atoi(r.URL.Query().Get("commits-count")); err == nil {
		commitsCount = parsedCommitsCount
	}

	gg := graphgen.NewContributionsGraphGenerator(
		graphgen.HtmlGeneratorType,
		graphgen.ContributionsGraph{}.Init(year),
	)

	switch font {
	case "3x3":
		gg.SetFont(graphgen.Font3x3)
	case "3x5":
		gg.SetFont(graphgen.Font3x5)
	}

	buf, err := gg.GetFinalForm(msg[0], commitsCount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_, _ = io.Copy(w, buf)
}

func (g *graphGeneratorApi) handleGenerateCheatScript(w http.ResponseWriter, r *http.Request) {
	msg, exists := r.URL.Query()["msg"]
	if !exists {
		log.Warningln("someone sent a bad request...")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	font := r.URL.Query().Get("font")
	year := time.Now().Year()
	if parsedYear, err := strconv.Atoi(r.URL.Query().Get("year")); err == nil {
		year = parsedYear
	}
	commitsCount := 80
	if parsedCommitsCount, err := strconv.Atoi(r.URL.Query().Get("commits-count")); err == nil {
		commitsCount = parsedCommitsCount
	}

	gg := graphgen.NewContributionsGraphGenerator(
		graphgen.CheatScriptGeneratorType,
		graphgen.ContributionsGraph{}.Init(year),
	)

	switch font {
	case "3x3":
		gg.SetFont(graphgen.Font3x3)
	case "3x5":
		gg.SetFont(graphgen.Font3x5)
	}

	buf, err := gg.GetFinalForm(msg[0], commitsCount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = io.Copy(w, buf)
}
