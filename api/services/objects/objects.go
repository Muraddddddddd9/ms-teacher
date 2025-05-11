package objects

import (
	"context"
	"fmt"
	"ms-teacher/api/constants"
	"strings"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetObjects(c *fiber.Ctx, db *mongo.Database, rdb *redis.Client) error {
	session := c.Cookies(constants.SessionName)

	if session == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message":  constants.ErrEntrySystem,
			"redirect": constants.RedirectPathLogin,
		})
	}

	sessionKey := fmt.Sprintf(constants.SessionKeyStart, session)
	userKey, err := rdb.Get(context.Background(), sessionKey).Result()
	if err != nil || userKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":  constants.ErrServerError,
			"redirect": constants.RedirectPathLogin,
		})
	}

	userID, err := primitive.ObjectIDFromHex(strings.Split(userKey, ":")[1])
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"messsage": constants.ErrServerError,
			"redirec":  constants.RedirectPathLogin,
		})
	}

	objectGroupRepo := mongodb.NewRepository[struct{}, models.ObjectsGroupsWithGroupAndTeacherModel](db.Collection(constants.ObjectGroupCollection))
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"teacher": userID,
			},
		},
		{
			"$lookup": bson.M{
				"from":         constants.GroupCollection,
				"localField":   "group",
				"foreignField": "_id",
				"as":           "groupData",
			},
		},
		{
			"$lookup": bson.M{
				"from":         constants.ObjectCollection,
				"localField":   "object",
				"foreignField": "_id",
				"as":           "objectData",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$groupData",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$objectData",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$project": bson.M{
				"_id":     1,
				"group":   "$groupData.group",
				"object":  "$objectData.object",
				"teacher": 1,
			},
		},
	}

	objectGroupFind, err := objectGroupRepo.AggregateAll(context.Background(), pipeline)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": constants.ErrGetData,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"objectGroups": objectGroupFind,
	})
}
