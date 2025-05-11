package evaluations

import (
	"context"
	"ms-teacher/api/constants"
	"strings"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StudentMinimal struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	Name       string             `bson:"name" json:"name"`
	Surname    string             `bson:"surname" json:"surname"`
	Patronymic string             `bson:"patronymic" json:"patronymic"`
}

type ObjectGroupMinimal struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	Object string             `bson:"object" json:"object"`
}

func GetEvaluation(c *fiber.Ctx, db *mongo.Database) error {
	group := c.Params("group")
	object := c.Params("object")

	if strings.TrimSpace(group) == "" || strings.TrimSpace(object) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrGetData,
		})
	}

	objectId, err := primitive.ObjectIDFromHex(object)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid object ID format",
		})
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": objectId,
			},
		},
		{
			"$lookup": bson.M{
				"from":         constants.ObjectCollection,
				"localField":   "object",
				"foreignField": "_id",
				"as":           "object_info",
			},
		},
		{
			"$unwind": "$object_info",
		},
		{
			"$project": bson.M{
				"_id":        0,
				"objectName": "$object_info.object",
			},
		},
	}

	type Result struct {
		ObjectName string `bson:"objectName"`
	}

	objectGroupRepo := mongodb.NewRepository[struct{}, Result](db.Collection(constants.ObjectGroupCollection))
	objectName, err := objectGroupRepo.AggregateAll(context.Background(), pipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error",
		})
	}

	groupRepo := mongodb.NewRepository[models.GroupsModel, struct{}](db.Collection(constants.GroupCollection))
	groupData, err := groupRepo.FindOne(context.Background(), bson.M{"group": group})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrStatusNotFound,
		})
	}

	studentRepo := mongodb.NewRepository[StudentMinimal, struct{}](db.Collection(constants.StudentCollection))
	students, err := studentRepo.FindAll(context.Background(), bson.M{"group": groupData.ID})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrStudentNotFound,
		})
	}

	evaluationRepo := mongodb.NewRepository[models.EvaluationModel, models.EvaluationModelWithStudent](db.Collection(constants.EvaluationCollection))
	evaluation, err := evaluationRepo.FindAll(context.Background(), bson.M{"object": objectId})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constants.ErrServerError,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"object":     objectName[0].ObjectName,
		"students":   students,
		"evaluation": evaluation,
	})
}
