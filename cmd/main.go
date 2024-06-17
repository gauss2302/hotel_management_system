package main

import (
	"context"
	"flag"
	"github.com/gauss2302/hotel_management_system/api/userHandler"
	"github.com/gauss2302/hotel_management_system/internal/database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const dburi = "mongodb://localhost:27017"

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})

	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(),
		options.Client().ApplyURI(dburi))

	if err != nil {
		log.Fatal(err)
	}

	userHandler := userHandler.NewUserHandler(database.NewMongoUserStore(client))

	app := fiber.New()
	apiV1 := app.Group("/api/v1")

	apiV1.Post("/user", userHandler.HandlePostUSer)
	apiV1.Get("/user", userHandler.HandleGetusers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)

	app.Listen(*listenAddr)
}
