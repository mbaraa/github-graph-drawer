package emailsched

import (
	"crypto/sha256"
	"encoding/hex"
	"github-graph-drawer/db"
	"strings"
	"time"
)

func SendScheduleConfirmationEmail(sr db.ScheduleRequest) error {
	buf := new(strings.Builder)
	err := templates.ExecuteTemplate(buf, "schedule_confirmation", map[string]any{
		"CommitsCount":     sr.Content.CommitsCount,
		"Msg":              sr.Content.Message,
		"Year":             sr.Content.Year,
		"ConfirmationLink": "https://github-graph-drawer.mbaraa.com/email/confirm?token=" + sr.ConfirmationToken,
	})
	if err != nil {
		return err
	}
	err = db.InsertScheduleRequest(sr)
	if err != nil {
		return err
	}
	return sendEmail("GitHub Graph Drawer Schedule Confirmation", buf.String(), sr.Email)
}

func ConfirmEmail(token string) error {
	scheduleRequest, err := db.GetScheduleRequestByEmailAndToken(token)
	if err != nil {
		return err
	}
	for _, date := range scheduleRequest.Dates {
		d, _ := time.Parse("2006-01-02", date)
		err := db.InsertDailySchedule(db.DailySchedule{
			Email:            scheduleRequest.Email,
			Date:             date,
			CancelationToken: scheduleRequest.ConfirmationToken,
			Content: db.EmailContent{
				CommitsCount: scheduleRequest.Content.CommitsCount,
				Message:      scheduleRequest.Content.Message,
				Year:         scheduleRequest.Content.Year,
			},
			CreatedAt: d.Unix(),
		})
		if err != nil {
			return err
		}
	}
	return db.DeleteScheduleRequestByEmailAndToken(scheduleRequest.Email, token)
}

func Unsubscribe(token string) error {
	return db.DeleteDailySchedulesByEmailAndToken(token)
}

func SendDailySchedulesEmail(time time.Time) (sent, total int, err error) {
	dailySchedules, err := db.GetDailySchedulesByTimestamp(time)
	if err != nil {
		return
	}
	for _, dailySchedule := range dailySchedules {
		buf := new(strings.Builder)
		err = templates.ExecuteTemplate(buf, "daily_commits_reminder", map[string]any{
			"CommitsCount":        dailySchedule.Content.CommitsCount,
			"Msg":                 dailySchedule.Content.Message,
			"UnsubscriptionEmail": "https://github-graph-drawer.mbaraa.com/email/unsubscribe?token=" + dailySchedule.CancelationToken,
		})
		if err != nil {
			return
		}
		err = sendEmail("GitHub Graph Drawer Daily Reminder", buf.String(), dailySchedule.Email)
		if err != nil {
			return
		}
		err = db.DeleteDailyScheduleById(dailySchedule.Id)
		if err != nil {
			return
		}
		sent++
	}
	return sent, len(dailySchedules), err
}

func generateToken() string {
	sha256 := sha256.New()
	sha256.Write([]byte(time.Now().String()))
	return hex.EncodeToString(sha256.Sum(nil))
}
