package db

import (
	"context"
	"errors"
	"github-graph-drawer/config"
	"github-graph-drawer/log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx                 context.Context
	client              *mongo.Client
	scheduleRequestColl *mongo.Collection
	dailyScheduleColl   *mongo.Collection

	ErrInvalidCreatedAtDate = errors.New("invalid created at date")
)

func init() {
	ctx = context.Background()
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(config.Config().DbUri))
	if err != nil {
		log.Errorln(err.Error())
	}
	scheduleRequestColl = client.Database("github-graph-drawer-db").Collection("ScheduleRequest")
	dailyScheduleColl = client.Database("github-graph-drawer-db").Collection("DailySchedule")
}

type EmailContent struct {
	CommitsCount int    `bson:"commitsCount,omitempty"`
	Message      string `bson:"message,omitempty"`
	Year         string `bson:"year"`
}

type ScheduleRequest struct {
	Email             string       `bson:"email,omitempty"`
	Dates             []string     `bson:"dates,omitempty"`
	ConfirmationToken string       `bson:"confirmationToken,omitempty"`
	Content           EmailContent `bson:"content"`
	CreatedAt         int64        `bson:"createdAt"`
}

type DailySchedule struct {
	Id               string       `bson:"_id,omitempty"`
	Email            string       `bson:"email,omitempty"`
	Date             string       `bson:"date,omitempty"`
	CancelationToken string       `bson:"cancelationToken,omitempty"`
	Content          EmailContent `bson:"content"`
	CreatedAt        int64        `bson:"createdAt"`
}

// ok
func InsertScheduleRequest(sr ScheduleRequest) error {
	if sr.CreatedAt == 0 {
		return ErrInvalidCreatedAtDate
	}
	_, err := scheduleRequestColl.InsertOne(ctx, sr)
	if err != nil {
		return err
	}
	return nil
}

// ok
func InsertDailySchedule(ds DailySchedule) error {
	if ds.CreatedAt == 0 {
		return ErrInvalidCreatedAtDate
	}
	_, err := dailyScheduleColl.InsertOne(ctx, ds)
	if err != nil {
		return err
	}
	return nil
}

// ok
func InsertDailySchedules(dss []DailySchedule) error {
	documents := make([]any, 0)
	for _, ds := range dss {
		documents = append(documents, ds)
	}
	_, err := dailyScheduleColl.InsertMany(ctx, documents)
	if err != nil {
		return err
	}
	return nil
}

// ok
func GetScheduleRequestByEmailAndToken(token string) (er ScheduleRequest, err error) {
	filter := bson.D{{Key: "confirmationToken", Value: bson.D{{Key: "$eq", Value: token}}}}
	result := scheduleRequestColl.FindOne(ctx, filter)
	if result.Err() != nil {
		return
	}
	err = result.Decode(&er)
	if err != nil {
		return
	}
	return
}

// ok
func GetDailySchedulesByTimestamp(t time.Time) (dss []DailySchedule, err error) {
	cursor, err := dailyScheduleColl.Find(ctx, bson.D{
		{Key: "createdAt", Value: bson.D{{Key: "$lte", Value: t.Unix()}}},
	})
	if err != nil {
		return
	}
	for cursor.Next(ctx) {
		var result DailySchedule
		if err := cursor.Decode(&result); err != nil {
			log.Errorln(err)
			continue
		}
		dss = append(dss, result)
	}
	if err := cursor.Err(); err != nil {
		log.Errorln(err)
		return nil, err
	}
	currentDayEmails := make([]DailySchedule, 0)
	for _, ds := range dss {
		iDidntHaveANameForThis := time.Unix(ds.CreatedAt, 0)
		if err != nil {
			continue
		}
		if iDidntHaveANameForThis.Day() == t.Day() &&
			iDidntHaveANameForThis.Month() == t.Month() &&
			iDidntHaveANameForThis.Year() == t.Year() {
			currentDayEmails = append(currentDayEmails, ds)
		}
	}
	return currentDayEmails, nil
}

// ok
func DeleteScheduleRequestByEmailAndToken(email, token string) error {
	filter := bson.D{{Key: "$and", Value: bson.A{
		bson.D{{Key: "confirmationToken", Value: bson.D{{Key: "$eq", Value: token}}}},
		bson.D{{Key: "email", Value: bson.D{{Key: "$eq", Value: email}}}},
	}}}
	_, err := scheduleRequestColl.DeleteOne(ctx, filter)
	return err
}

// ok
func DeleteDailyScheduleById(id string) error {
	primitiveId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$eq", Value: primitiveId}}}}
	_, err = dailyScheduleColl.DeleteOne(ctx, filter)
	return err
}

// ok
func DeleteDailySchedulesByEmailAndToken(token string) error {
	filter := bson.D{{Key: "cancelationToken", Value: bson.D{{Key: "$eq", Value: token}}}}
	_, err := dailyScheduleColl.DeleteMany(ctx, filter)
	return err
}
