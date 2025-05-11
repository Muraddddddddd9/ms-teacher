package main

import (
	"context"
	"fmt"
	"log"
	"ms-teacher/api/constants"
	"ms-teacher/api/services/evaluations"
	"ms-teacher/api/services/objects"
	loconfig "ms-teacher/config"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/data/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	db, err := mongodb.Connect()
	if err == nil {
		log.Println(constants.SuccConnectMongo)
	}
	defer db.Client().Disconnect(context.Background())

	rdb, err := redis.Connect()
	if err == nil {
		log.Println(constants.SuccConnectRedis)
	}
	defer rdb.Close()

	cfg, err := loconfig.LoadLocalConfig()
	if err != nil {
		log.Print(err)
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.ORIGIN_URL,
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	app.Get("/api/teacher/get_objects", TeacherOnly(rdb), func(c *fiber.Ctx) error {
		return objects.GetObjects(c, db, rdb)
	})

	app.Get("/api/teacher/get_evaluation/:group/:object", TeacherOnly(rdb), func(c *fiber.Ctx) error {
		return evaluations.GetEvaluation(c, db)
	})

	app.Post("/api/teacher/send_evaluation", TeacherOnly(rdb), func(c *fiber.Ctx) error {
		return evaluations.SendEvaluation(c, db)
	})

	app.Patch("/api/teacher/update_evaluation", TeacherOnly(rdb), func(c *fiber.Ctx) error {
		return nil
	})

	app.Delete("/api/teacher/delete_evaluation", TeacherOnly(rdb), func(c *fiber.Ctx) error {
		return nil
	})

	app.Post("/api/teacher/create_contests", TeacherOnly(rdb), func(c *fiber.Ctx) error {
		return nil
	})

	app.Listen(fmt.Sprintf("localhost%v", cfg.PROJECT_PORT))
}
