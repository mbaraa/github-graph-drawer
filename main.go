package main

import (
	"embed"
	"io"
	"net/http"
	"text/template"
	"time"

	"github-graph-drawer/log"
	"github-graph-drawer/utils/graphgen"
)

var (
	//go:embed resources/*
	res embed.FS
)

func main() {
	templates := template.Must(template.ParseGlob("./templates/html/*"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if err := templates.ExecuteTemplate(w, "index", nil); err != nil {
			log.Errorln(err.Error())
			return
		}
	})

	http.HandleFunc("/contribution-graph", func(w http.ResponseWriter, r *http.Request) {
		msg, exists := r.URL.Query()["msg"]
		if !exists {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		gg := graphgen.NewContributionsGraphGenerator(
			graphgen.HtmlGeneratorType,
			graphgen.ContributionsGraph{}.Init(time.Now().Year()),
		)

		buf, err := gg.GetFinalForm(msg[0])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorln(err.Error())
			return
		}

		w.Header().Set("Content-Type", "text/html")
		_, _ = io.Copy(w, buf)
	})

	http.Handle("/resources/", http.FileServer(http.FS(res)))

	log.Infoln("server started...")
	log.Fatalln(string(log.ErrorLevel), http.ListenAndServe(":8080", nil))
}
