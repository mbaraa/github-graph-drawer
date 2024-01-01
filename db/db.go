package db

import (
	"context"
	"errors"
	"github-graph-drawer/config"
	"github-graph-drawer/log"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx               context.Context
	client            *mongo.Client
	emailRequestsColl *mongo.Collection
	emailScheduleColl *mongo.Collection

	ErrInvalidCreatedAtDate = errors.New("invalid created at date")
)

func init() {
	ctx = context.Background()
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(config.Config().DbUri))
	if err != nil {
		log.Errorln(err.Error())
	}
	emailRequestsColl = client.Database("github-graph-drawer-db").Collection("emailRequests")
	emailScheduleColl = client.Database("github-graph-drawer-db").Collection("emailSchedule")
}

type EmailRequestOperation string

const (
	StartSchedule EmailRequestOperation = "start"
	StopSchedule  EmailRequestOperation = "stop"
)

type EmailRequest struct {
	Dates        []string              `bson:"dates,omitempty"`
	Email        string                `bson:"email,omitempty"`
	Token        string                `bson:"token,omitempty"`
	Operation    EmailRequestOperation `bson:"operation,omitempty"`
	CreatedAt    int64                 `bson:"createdAt"`
	Message      string                `bson:"message,omitempty"`
	CommitsCount int                   `bson:"commitsCount,omitempty"`
}

type EmailScedule struct {
	Email        string `bson:"email,omitempty"`
	Token        string `bson:"token,omitempty"`
	Date         string `bson:"date,omitempty"`
	Message      string `bson:"message,omitempty"`
	CommitsCount int    `bson:"commitsCount,omitempty"`
	CreatedAt    int64  `bson:"createdAt"`
}

func InsertEmailRequest(er EmailRequest) error {
	if er.CreatedAt == 0 || er.CreatedAt < time.Now().Unix() {
		return ErrInvalidCreatedAtDate
	}
	_, err := emailRequestsColl.InsertOne(ctx, er)
	if err != nil {
		return err
	}
	return nil
}

func GetEmailRequests(email string) (ers []EmailRequest, err error) {
	cursor, err := emailRequestsColl.Find(ctx, bson.D{{"email", bson.D{{"$eq", email}}}})
	if err != nil {
		return
	}

	for cursor.Next(ctx) {
		var result EmailRequest
		if err := cursor.Decode(&result); err != nil {
			log.Errorln(err)
			continue
		}
		ers = append(ers, result)
	}
	if err := cursor.Err(); err != nil {
		log.Errorln(err)
		return nil, err
	}

	if ers != nil {
		slices.SortFunc(ers, func(a, b EmailRequest) int {
			return int(a.CreatedAt) - int(b.CreatedAt)
		})
	}

	return
}

func DeleteEmailRequests(email string) error {
	_, err := emailRequestsColl.DeleteMany(ctx, bson.D{{"email", bson.D{{"$eq", email}}}})
	if err != nil {
		return err
	}
	return nil
}

func InsertEmailSchedule(es EmailScedule) error {
	if es.CreatedAt == 0 || es.CreatedAt < time.Now().Unix() {
		return ErrInvalidCreatedAtDate
	}
	_, err := emailScheduleColl.InsertOne(ctx, es)
	if err != nil {
		return err
	}
	return nil
}

func GetEmailSchedules(date time.Time) (ess []EmailScedule, err error) {
	cursor, err := emailScheduleColl.Find(ctx, bson.D{{"createdAt", bson.D{{"$lte", date.Unix()}}}})
	if err != nil {
		return
	}

	for cursor.Next(ctx) {
		var result EmailScedule
		if err := cursor.Decode(&result); err != nil {
			log.Errorln(err)
			continue
		}
		ess = append(ess, result)
	}
	if err := cursor.Err(); err != nil {
		log.Errorln(err)
		return nil, err
	}

	return
}
