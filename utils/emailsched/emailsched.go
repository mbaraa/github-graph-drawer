package emailsched

import (
	"github-graph-drawer/db"
	"strings"
)

func SendScheduleConfirmationEmail(er db.EmailRequest) error {
	buf := new(strings.Builder)
	err := templates.ExecuteTemplate(buf, "schedule_confirmation", map[string]any{
		"CommitsCount":     er.CommitsCount,
		"Msg":              er.Message,
		"ConfirmationLink": "https://github-graph-drawer.mbaraa.com/confirm-email?token=" + er.Token,
	})
	if err != nil {
		return err
	}

	return sendEmail("GitHub Graph Drawer Schedule Confirmation", buf.String(), er.Email)

}

func SendDailyCommitsEmail(es db.EmailScedule) error {
	buf := new(strings.Builder)
	err := templates.ExecuteTemplate(buf, "daily_commits_reminder", map[string]any{
		"CommitsCount":        es.CommitsCount,
		"Msg":                 es.Message,
		"UnsubscriptionEmail": "https://github-graph-drawer.mbaraa.com/unsubscribe-email?token=" + es.Token,
	})
	if err != nil {
		return err
	}

	return sendEmail("GitHub Graph Drawer Daily Reminder", buf.String(), es.Email)
}

func UnsubscribeFromTheThing(er db.EmailRequest) error {
	return nil
}
