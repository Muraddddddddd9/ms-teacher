package services

import (
	"context"
	"fmt"
	"ms-teacher/api/constants"
	"strings"
	"time"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserID(rdb *redis.Client, session string) (primitive.ObjectID, error) {
	if session == "" {
		return primitive.NilObjectID, fmt.Errorf(constants.ErrEntrySystem)
	}

	sessionKey := fmt.Sprintf(constants.SessionKeyStart, session)
	userKey, err := rdb.Get(context.Background(), sessionKey).Result()
	if err != nil || userKey == "" {
		return primitive.NilObjectID, fmt.Errorf(constants.ErrServerError)
	}

	userID, err := primitive.ObjectIDFromHex(strings.Split(userKey, ":")[1])
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf(constants.ErrServerError)
	}

	return userID, err
}

func Logging(db *mongo.Database, api, method, status string, data any, errData any) {
	document := models.Log{
		API:    api,
		Method: method,
		Status: status,
		Data:   data,
		Date:   time.Now().Local().Format("2006-01-02 15:04:05 MST"),
		Error:  errData,
	}
	logRepo := mongodb.NewRepository[models.Log, struct{}](db.Collection(constants.LogsCollection))
	_, err := logRepo.InsertOne(context.Background(), &document)
	if err != nil {
		log.Errorf(constants.ErrDataLogging)
	}
}
