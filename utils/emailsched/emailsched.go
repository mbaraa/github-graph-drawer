package emailsched

import (
	"github-graph-drawer/db"
	"strings"
	"time"

	"github.com/google/uuid"
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
	err = db.InsertEmailRequest(er)
	if err != nil {
		return err
	}

	return sendEmail("GitHub Graph Drawer Schedule Confirmation", buf.String(), er.Email)
}

func ConfirmDailySchedule(er db.EmailRequest) error {
	err := db.DeleteEmailRequests(er.Token)
	if err != nil {
		return err
	}
	return StartScheduleForIdkItsLateIHaveNoIdeaWhyImStillAwakeDotDotDot(er)
}

func SendDailyCommitsEmail(es db.EmailSchedule) error {
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

func StartScheduleForIdkItsLateIHaveNoIdeaWhyImStillAwakeDotDotDot(er db.EmailRequest) error {
	for _, date := range er.Dates {
		err := db.InsertEmailSchedule(db.EmailSchedule{
			Id:           uuid.NewString(),
			Email:        er.Email,
			Token:        er.Token,
			Date:         date,
			Message:      er.Message,
			CommitsCount: er.CommitsCount,
			CreatedAt:    time.Now().UnixMilli(),
		})
		if err != nil {
			continue
		}
	}
	return nil
}

func UnsubscribeFromTheThing(email string) error {
	return db.DeleteEmailRequests(email)
}
