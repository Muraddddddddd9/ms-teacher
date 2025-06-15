package services

import (
	"context"
	"ms-teacher/api/constants"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SendMessage struct {
	Email       string `bson:"email" json:"email"`
	Description string `bson:"description" json:"description"`
}

func MessageContest(c *fiber.Ctx, db *mongo.Database) error {
	var session = c.Cookies(constants.SessionName)
	var sendMessage SendMessage

	if err := c.BodyParser(&sendMessage); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidData,
		})
	}

	studntRepo := mongodb.NewRepository[models.StudentsModel, struct{}](db.Collection(constants.StudentCollection))
	studentFindOne, err := studntRepo.FindOne(context.Background(), bson.M{"email": sendMessage.Email})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": constants.ErrStudentNotFound,
		})
	}

	if studentFindOne.Telegram > 0 {
		err := NotificationSend(studentFindOne.ID, sendMessage.Description, session)
		if err != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": constants.ErrSendMessage,
			})
		}
	} else {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Студент не подключил телеграм",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Сообщение отправлено",
	})
}
