package api

import (
	"errors"
	"fmt"
	"github.com/gauss2302/hotel_management_system/internal/database"
	"github.com/gauss2302/hotel_management_system/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
	"time"
)

type AuthHandler struct {
	userStore database.UserStore
}

func NewAuthHandler(userStore database.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

type genericResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func CreateTokenFromUser(user *models.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 12).Unix()

	claims := jwt.MapClaims{
		"userID":  user.ID,
		"email":   user.Email,
		"expires": expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		fmt.Println("Failed To Sign Token With Secret", err)
	}

	return tokenString
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return InvalidCredentials(c)
		}
		return err
	}

	if !models.IsValidPassword(user.EncryptedPassword, params.Password) {
		return InvalidCredentials(c)
	}

	resp := AuthResponse{
		User:  user,
		Token: CreateTokenFromUser(user),
	}
	return c.JSON(resp)
}

func InvalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResp{
		Type: "error",
		Msg:  "Invalid Credentials",
	})
}
