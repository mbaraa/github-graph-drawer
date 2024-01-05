package main

import (
	"net/http"
	"time"

	"github-graph-drawer/apis"
	"github-graph-drawer/log"
	"github-graph-drawer/utils/emailsched"

	"github.com/robfig/cron"
)

func main() {
	scheduleEmails()
	startServer()
}

func scheduleEmails() {
	cronie := cron.New()
	err := cronie.AddFunc("0 0 * * * *", func() {
		log.Infoln("Sending emails")
		sent, total, err := emailsched.SendDailySchedulesEmail(time.Now())
		if err != nil {
			log.Errorln(err)
		}
		if sent != total {
			log.Warningf("%d of the %d emails weren't sent\n", (total - sent), total)
		}
		log.Infoln("Done sending emails...")
	})
	if err != nil {
		log.Errorln(err.Error())
	}
	cronie.Start()
}

func startServer() {
	for _, api := range []apis.IHandler{
		apis.NewPagesApi(),
		apis.NewGraphGeneratorApi(),
		apis.NewEmailApi(),
	} {
		log.Infof("mounting the api %s\n", api.Prefix())
		http.Handle(api.Prefix(), api)
	}
	log.Errorln(http.ListenAndServe(":8080", nil))
}
