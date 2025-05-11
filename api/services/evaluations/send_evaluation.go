package evaluations

import (
	"context"
	"ms-teacher/api/constants"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SendEvaluation(c *fiber.Ctx, db *mongo.Database) error {
	var evaluationData models.EvaluationModel

	if err := c.BodyParser(&evaluationData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Не вернный ввод данных",
		})
	}

	studentRepo := mongodb.NewRepository[models.StudentsModel, struct{}](db.Collection(constants.StudentCollection))
	_, err := studentRepo.FindOne(context.Background(), bson.M{"_id": evaluationData.Student})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Пользователь не найден",
		})
	}

	objectGroupRepo := mongodb.NewRepository[models.ObjectsGroupsModel, struct{}](db.Collection(constants.ObjectGroupCollection))
	_, err = objectGroupRepo.FindOne(context.Background(), bson.M{"_id": evaluationData.Object})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Группа не найдена",
		})
	}

	evaluationRepo := mongodb.NewRepository[models.EvaluationModel, struct{}](db.Collection(constants.EvaluationCollection))
	_, err = evaluationRepo.InsertOne(context.Background(), &evaluationData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Ошибка в добалвении оценки",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Оценка отправлена",
	})
}
