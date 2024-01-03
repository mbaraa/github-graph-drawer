package apis

import (
	"github-graph-drawer/log"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"time"
)

type pagesApi struct {
	templates *template.Template
	resDir    fs.FS
}

func NewPagesApi() IHandler {
	p := &pagesApi{
		templates: template.Must(template.ParseGlob("./templates/html/*")),
		resDir:    os.DirFS("./resources"),
	}
	endpoints := Endpoints{
		"GET ":           p.handleHomePage,
		"GET resources":  p.handleResources,
		"GET robots.txt": p.handleRobotsFile,
	}
	return NewHandler(endpoints, "/")
}

func (p *pagesApi) handleRobotsFile(w http.ResponseWriter, r *http.Request) {
	robotsFile, _ := os.ReadFile("./resources/robots.txt")
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write(robotsFile)
}

func (p *pagesApi) handleResources(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[1:]
	file, err := os.ReadFile(filePath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", http.DetectContentType(file))
	_, _ = w.Write(file)
}

func (p *pagesApi) handleHomePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := p.templates.ExecuteTemplate(w, "index", map[string]any{
		"CurrentYear":         time.Now().Year(),
		"DefaultCommitsCount": 80,
	}); err != nil {
		log.Errorln(err.Error())
		return
	}
}
