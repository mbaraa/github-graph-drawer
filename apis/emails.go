package apis

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github-graph-drawer/db"
	"github-graph-drawer/log"
	"github-graph-drawer/utils/emailsched"
	"github-graph-drawer/utils/graphgen"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type emailsApi struct{}

func NewEmailApi() IHandler {
	e := &emailsApi{}
	endpoints := Endpoints{
		"GET schedule":    e.handleScheduleEmails,
		"GET confirm":     e.handleConfirmEmail,
		"GET unsubscribe": e.handleUnsubscribeEmails,
	}
	return NewHandler(endpoints, "/email/")
}

func (e *emailsApi) handleScheduleEmails(w http.ResponseWriter, r *http.Request) {
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

	email, exists := r.URL.Query()["email"]
	if !exists {
		log.Warningln("someone sent a bad request...")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	gg := graphgen.NewContributionsGraphGenerator(
		graphgen.EmailScheduleGeneratorType,
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
	strDates, _ := io.ReadAll(buf)
	dates := strings.Split(string(strDates), " ")
	err = emailsched.SendScheduleConfirmationEmail(db.ScheduleRequest{
		Email:             email[0],
		Dates:             dates,
		ConfirmationToken: generateToken(),
		Content: db.EmailContent{
			CommitsCount: commitsCount,
			Message:      msg[0],
			Year:         fmt.Sprint(year),
		},
		CreatedAt: time.Now().Unix(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte("<span style=\"color: white;\">Done, check your email<span>"))
}

func (e *emailsApi) handleConfirmEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if len(token) == 0 {
		log.Warningln("someone tried to confirm their email with an empty token")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := emailsched.ConfirmEmail(token)
	if err != nil {
		log.Errorln(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte("wohoo, get ready for the daily email reminding you about the commits."))
}

func (e *emailsApi) handleUnsubscribeEmails(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if len(token) == 0 {
		log.Warningln("someone tried to unsubscribe their email with an empty token")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := emailsched.Unsubscribe(token)
	if err != nil {
		log.Errorln(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte("unsubscribed!"))
}

/*
	http.HandleFunc("/schedule-emails", func(w http.ResponseWriter, r *http.Request) {
	})

	http.HandleFunc("/confirm-email", func(w http.ResponseWriter, r *http.Request) {
	})

	http.HandleFunc("/unsubscribe-email", func(w http.ResponseWriter, r *http.Request) {
	})
*/

func generateToken() string {
	sha256 := sha256.New()
	sha256.Write([]byte(time.Now().String()))
	return hex.EncodeToString(sha256.Sum(nil))
}
