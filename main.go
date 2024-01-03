package main

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github-graph-drawer/apis"
	"github-graph-drawer/config"
	"github-graph-drawer/db"
	"github-graph-drawer/log"
	"github-graph-drawer/utils/emailsched"
	"github-graph-drawer/utils/graphgen"

	"github.com/robfig/cron"
)

var (
	//go:embed resources/*
	res embed.FS
)

func main() {
	// scheduleEmails()
	// startServer()
	for _, api := range []apis.IHandler{
		apis.NewPagesApi(),
		apis.NewGraphGeneratorApi(),
	} {
		log.Infof("mounting the api %s\n", api.Prefix())
		http.Handle(api.Prefix(), api)
	}

	http.ListenAndServe(":8080", nil)
}

func scheduleEmails() {
	cronie := cron.New()
	err := cronie.AddFunc("0 0 * * * *", func() {
		emailSchedules, err := db.GetEmailSchedules(time.Now())
		if err != nil {
			log.Errorln(err.Error())
		}
		for _, es := range emailSchedules {
			err = emailsched.SendDailyCommitsEmail(es)
			if err != nil {
				log.Errorln(err.Error())
			}
			err = db.DeleteEmailSchedule(es.Id)
			if err != nil {
				log.Errorln(err.Error())
			}
		}
	})
	if err != nil {
		log.Errorln(err.Error())
	}
	cronie.Start()
}

func generateToken() string {
	sha256 := sha256.New()
	sha256.Write([]byte(time.Now().String()))
	return hex.EncodeToString(sha256.Sum(nil))
}

func startServer() {
	http.HandleFunc("/schedule-emails", func(w http.ResponseWriter, r *http.Request) {
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
		err = emailsched.SendScheduleConfirmationEmail(db.EmailRequest{
			Dates:        dates,
			Email:        email[0],
			Token:        generateToken(),
			Operation:    db.StartSchedule,
			CreatedAt:    time.Now().UnixMilli(),
			Message:      msg[0],
			CommitsCount: commitsCount,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorln(err.Error())
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<span style=\"color: white;\">Done, check your email<span>"))
	})

	http.HandleFunc("/confirm-email", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if len(token) == 0 {
			log.Warningln("someone tried to confirm their email with an empty token")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		reqs, err := db.GetEmailRequests(token)
		if err != nil {
			log.Errorln(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = emailsched.ConfirmDailySchedule(reqs[0])
		if err != nil {
			log.Errorln(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/unsubscribe-email", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if len(token) == 0 {
			log.Warningln("someone tried to unsubscribe their email with an empty token")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		reqs, err := db.GetEmailRequests(token)
		if err != nil {
			log.Errorln(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = db.DeleteEmailRequests(reqs[0].Email)
		if err != nil {
			log.Errorln(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	log.Infoln("server started at http://localhost:" + config.Config().Port)
	log.Fatalln(string(log.ErrorLevel), http.ListenAndServe(":"+config.Config().Port, nil))
}
