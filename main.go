package main

import (
	"embed"
	"io"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github-graph-drawer/db"
	"github-graph-drawer/log"
	"github-graph-drawer/utils/graphgen"
)

var (
	//go:embed resources/*
	res embed.FS
)

func main() {
	_ = db.EmailMsg{}
	templates := template.Must(template.ParseGlob("./templates/html/*"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robots.txt" {
			tmpl, _ := template.ParseFS(res, "resources/robots.txt")
			w.Header().Set("Content-Type", "text/plain")
			_ = tmpl.Execute(w, nil)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		if err := templates.ExecuteTemplate(w, "index", nil); err != nil {
			log.Errorln(err.Error())
			return
		}
	})

	http.HandleFunc("/contribution-graph", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("/generate-script", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("/schedule-emails", func(w http.ResponseWriter, r *http.Request) {
		msg, exists := r.URL.Query()["msg"]
		if !exists {
			log.Warningln("someone sent a bad request...")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = msg

		// TODO: add zis to ze generator interface
		commitCount, exists := r.URL.Query()["commit-count"]
		if !exists {
			log.Warningln("someone sent a bad request...")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = commitCount

		email, exists := r.URL.Query()["email"]
		if !exists {
			log.Warningln("someone sent a bad request...")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = email
		log.Info("email: ", email)

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<span style=\"color: white;\">NOT IMPLEMENTED :)</span>"))
	})

	http.Handle("/resources/", http.FileServer(http.FS(res)))

	log.Infoln("server started...")
	log.Fatalln(string(log.ErrorLevel), http.ListenAndServe(":8080", nil))
}
