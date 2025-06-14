package evaluations

import (
	"context"
	"fmt"
	"ms-teacher/api/constants"
	"ms-teacher/api/services"
	"strconv"
	"strings"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetMyClassroomGroup(c *fiber.Ctx, db *mongo.Database, rdb *redis.Client) error {
	session := c.Cookies(constants.SessionName)

	userID, err := services.GetUserID(rdb, session)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"messsage": constants.ErrServerError,
			"redirec":  constants.RedirectPathLogin,
		})
	}

	groupRepo := mongodb.NewRepository[models.GroupsModel, struct{}](db.Collection(constants.GroupCollection))
	groupFindAll, err := groupRepo.FindAll(context.Background(), bson.M{"teacher": userID})
	if err != nil {
		services.Logging(db, "/api/common/get_my_classroom_group", c.Method(), strconv.Itoa(fiber.StatusNotFound), nil, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": constants.ErrGroupNotFound,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"groups": groupFindAll,
	})
}

func GetMyClassroomObject(c *fiber.Ctx, db *mongo.Database, rdb *redis.Client) error {
	session := c.Cookies(constants.SessionName)
	group := strings.TrimSpace(c.Params("group"))

	if session == "" || group == "" {
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

	type GroupData struct {
		Group string `bson:"group"`
	}
	groupData := GroupData{Group: group}

	_, err = primitive.ObjectIDFromHex(strings.Split(userKey, ":")[1])
	if err != nil {
		services.Logging(db, "/api/common/get_my_classroom_object", c.Method(), strconv.Itoa(fiber.StatusConflict), groupData, err)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"messsage": constants.ErrServerError,
			"redirec":  constants.RedirectPathLogin,
		})
	}

	groupID, err := primitive.ObjectIDFromHex(group)
	if err != nil {
		services.Logging(db, "/api/common/get_my_classroom_object", c.Method(), strconv.Itoa(fiber.StatusConflict), groupData, err)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"messsage": constants.ErrServerError,
		})
	}

	objectsGroupRepo := mongodb.NewRepository[struct{}, models.ObjectsGroupsWithGroupAndTeacherModel](db.Collection(constants.ObjectGroupCollection))
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"group": groupID,
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
			"$unwind": "$objectData",
		},
		{
			"$project": bson.M{
				"_id":     1,
				"object":  "$objectData.object",
				"group":   1,
				"teacher": 1,
			},
		},
	}

	objectsAggregateAll, err := objectsGroupRepo.AggregateAll(context.Background(), pipeline)
	if err != nil {
		services.Logging(db, "/api/common/get_my_classroom_object", c.Method(), strconv.Itoa(fiber.StatusNotFound), groupData, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": constants.ErrObjectNotFound,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"objects": objectsAggregateAll,
	})
}
