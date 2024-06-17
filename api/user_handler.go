package api

import (
	"context"
	"errors"
	"github.com/gauss2302/hotel_management_system/internal/database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore database.UseStore
}

func NewUserHandler(userStore database.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id  = c.Params("id")
		ctx = context.Background()
	)
	user, err := h.userStore.GetUserByID(ctx, id)
	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "Not Found"})
		}

		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return ErrNotResourceNotFound("user")
	}
	return c.JSON(users)
}
