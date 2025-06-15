package evaluations

import (
	"context"
	"fmt"
	"log"
	"ms-teacher/api/constants"
	"ms-teacher/api/services"
	"strconv"
	"strings"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteEvaluation(c *fiber.Ctx, db *mongo.Database) error {
	session := c.Cookies(constants.SessionName)
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidData,
		})
	}

	type IdData struct {
		ID string `bson:"_id"`
	}
	idData := IdData{ID: id}

	evaluationID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		services.Logging(db, "/api/common/delete_evaluation", c.Method(), strconv.Itoa(fiber.StatusBadRequest), idData, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidData,
		})
	}

	filter := bson.M{
		"_id": evaluationID,
	}
	evaluationRepo := mongodb.NewRepository[models.EvaluationModel, struct{}](db.Collection(constants.EvaluationCollection))
	evaluationFindOne, err := evaluationRepo.FindOne(context.Background(), filter)
	if err != nil {
		services.Logging(db, "/api/common/delete_evaluation", c.Method(), strconv.Itoa(fiber.StatusNotFound), idData, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": constants.ErrEvaluationNotFound,
		})
	}

	err = evaluationRepo.DeleteOne(context.Background(), filter)
	if err != nil {
		services.Logging(db, "/api/common/delete_evaluation", c.Method(), strconv.Itoa(fiber.StatusBadRequest), idData, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrDeleteEvaluation,
		})
	}

	objectGroupRepo := mongodb.NewRepository[models.ObjectsGroupsModel, struct{}](db.Collection(constants.ObjectGroupCollection))
	objectGroupFindOne, err := objectGroupRepo.FindOne(context.Background(), bson.M{"_id": evaluationFindOne.Object})
	if err != nil {
		services.Logging(db, "/api/common/delete_evaluation", c.Method(), strconv.Itoa(fiber.StatusNotFound), idData, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": constants.ErrObjectNotFound,
		})
	}

	objectRepo := mongodb.NewRepository[models.ObjectsModel, struct{}](db.Collection(constants.ObjectCollection))
	objectFindOne, err := objectRepo.FindOne(context.Background(), bson.M{"_id": objectGroupFindOne.Object})
	if err != nil {
		services.Logging(db, "/api/common/delete_evaluation", c.Method(), strconv.Itoa(fiber.StatusNotFound), idData, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": constants.ErrObjectNotFound,
		})
	}

	studentRepo := mongodb.NewRepository[models.StudentsModel, struct{}](db.Collection(constants.StudentCollection))
	studentFindOne, err := studentRepo.FindOne(context.Background(), bson.M{"_id": evaluationFindOne.Student})
	if err != nil {
		services.Logging(db, "/api/common/delete_evaluation", c.Method(), strconv.Itoa(fiber.StatusNotFound), idData, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": constants.ErrStudentNotFound,
		})
	}

	if studentFindOne.Telegram != 0 {
		str := fmt.Sprintf("Вам удалили '%v' по предмету '%v' за %v", evaluationFindOne.Value, strings.ToUpper(objectFindOne.Object), evaluationFindOne.Date)
		err = services.NotificationSend(evaluationFindOne.Student, str, session)
		if err != nil {
			log.Println(err)
		}
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": constants.SuccDeleteEvaluation,
	})
}
