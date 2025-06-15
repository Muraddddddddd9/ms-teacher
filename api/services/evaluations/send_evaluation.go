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
	"go.mongodb.org/mongo-driver/mongo"
)

func SendEvaluation(c *fiber.Ctx, db *mongo.Database) error {
	session := c.Cookies(constants.SessionName)
	var evaluationData models.EvaluationModel

	if err := c.BodyParser(&evaluationData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidData,
		})
	}

	fileds := map[string]string{
		"value": evaluationData.Value,
		"date":  evaluationData.Date,
	}

	for _, v := range fileds {
		if v == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": constants.ErrInvalidData,
			})
		}
	}

	studentRepo := mongodb.NewRepository[models.StudentsModel, struct{}](db.Collection(constants.StudentCollection))
	studentFindOne, err := studentRepo.FindOne(context.Background(), bson.M{"_id": evaluationData.Student})
	if err != nil {
		services.Logging(db, "/api/common/send_evaluation", c.Method(), strconv.Itoa(fiber.StatusBadRequest), evaluationData, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": constants.ErrStudentNotFound,
		})
	}

	objectGroupRepo := mongodb.NewRepository[models.ObjectsGroupsModel, struct{}](db.Collection(constants.ObjectGroupCollection))
	objectGroupFindOne, err := objectGroupRepo.FindOne(context.Background(), bson.M{"_id": evaluationData.Object})
	if err != nil {
		services.Logging(db, "/api/common/send_evaluation", c.Method(), strconv.Itoa(fiber.StatusBadRequest), evaluationData, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": constants.ErrObjectNotFound,
		})
	}

	evaluationRepo := mongodb.NewRepository[models.EvaluationModel, struct{}](db.Collection(constants.EvaluationCollection))
	_, err = evaluationRepo.InsertOne(context.Background(), &evaluationData)
	if err != nil {
		services.Logging(db, "/api/common/send_evaluation", c.Method(), strconv.Itoa(fiber.StatusBadRequest), evaluationData, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrSendEvaluation,
		})
	}

	objectRepo := mongodb.NewRepository[models.ObjectsModel, struct{}](db.Collection(constants.ObjectCollection))
	objectFindOne, err := objectRepo.FindOne(context.Background(), bson.M{"_id": objectGroupFindOne.Object})
	if err != nil {
		services.Logging(db, "/api/common/send_evaluation", c.Method(), strconv.Itoa(fiber.StatusBadRequest), evaluationData, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": constants.ErrObjectNotFound,
		})
	}

	if studentFindOne.Telegram != 0 {
		str := fmt.Sprintf("Вам поставили '%v' по предмету '%v' за %v", evaluationData.Value, strings.ToUpper(objectFindOne.Object), evaluationData.Date)
		err = services.NotificationSend(evaluationData.Student, str, session)
		if err != nil {
			log.Println(err)
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": constants.SuccSendEvaluation,
	})
}
