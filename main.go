package main

import (
	"context"
	"fmt"
	"log"
	"ms-teacher/api/constants"
	"ms-teacher/api/services/evaluations"
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
		AllowMethods:     "GET, POST, DELETE",
		AllowCredentials: true,
	}))

	app.Get("/api/teacher/get_evaluation/:group/:object", Access(rdb, []string{constants.AdminStatus, constants.RestrictedAdminStatus, constants.TeacherStatus}), func(c *fiber.Ctx) error {
		return evaluations.GetEvaluation(c, db)
	})

	app.Post("/api/teacher/send_evaluation", Access(rdb, []string{constants.AdminStatus, constants.RestrictedAdminStatus, constants.TeacherStatus}), func(c *fiber.Ctx) error {
		return evaluations.SendEvaluation(c, db)
	})

	app.Delete("/api/teacher/delete_evaluation/:id", Access(rdb, []string{constants.AdminStatus, constants.RestrictedAdminStatus, constants.TeacherStatus}), func(c *fiber.Ctx) error {
		return evaluations.DeleteEvaluation(c, db)
	})

	app.Get("/api/teacher/get_my_classroom_group", Access(rdb, []string{constants.AdminStatus, constants.RestrictedAdminStatus, constants.TeacherStatus}), func(c *fiber.Ctx) error {
		return evaluations.GetMyClassroomGroup(c, db, rdb)
	})

	app.Get("/api/teacher/get_my_classroom_object/:group", Access(rdb, []string{constants.AdminStatus, constants.RestrictedAdminStatus, constants.TeacherStatus}), func(c *fiber.Ctx) error {
		return evaluations.GetMyClassroomObject(c, db, rdb)
	})

	app.Listen(fmt.Sprintf("%v", cfg.PROJECT_PORT))
}
