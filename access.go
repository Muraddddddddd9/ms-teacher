package main

import (
	"context"
	"encoding/json"
	"fmt"
	"ms-teacher/api/constants"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func TeacherOnly(rdb *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()

		session := c.Cookies(constants.SessionName)
		if session == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message":  constants.ErrEntrySystem,
				"redirect": constants.RedirectPathProfile,
			})
		}

		sessionKey := fmt.Sprintf(constants.SessionKeyStart, session)

		userKey, err := rdb.Get(context.Background(), sessionKey).Result()
		if err != redis.Nil || userKey == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message":  constants.ErrEntrySystem,
				"redirect": constants.RedirectPathProfile,
			})
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message":  constants.ErrServerError,
				"redirect": constants.RedirectPathProfile,
			})
		}

		userData, err := rdb.Get(context.Background(), userKey).Bytes()
		if err == redis.Nil || userData == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message":  constants.ErrUserNotFound,
				"redirect": constants.RedirectPathProfile,
			})
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message":  constants.ErrServerError,
				"redirect": constants.RedirectPathProfile,
			})
		}

		var user struct {
			Status string `json:"status"`
		}
		err = json.Unmarshal(userData, &user)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"message":  constants.ErrGetData,
				"redirect": constants.RedirectPathProfile,
			})
		}

		if user.Status != "учитель" {
			return c.Status(301).JSON(fiber.Map{
				"redirect": constants.RedirectPathProfile,
			})
		}

		return c.Next()
	}
}
