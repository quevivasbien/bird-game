package api

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/quevivasbien/bird-backend/db"
)

var JWT_SECRET = []byte(os.Getenv("BIRD_JWT_SECRET"))

const JWT_EXPIRE_HOURS = 12
const JWT_COOKIE_NAME = "access_token"

type JWTPayload struct {
	Name       string
	Admin      bool
	ExpireTime int64
}

func getToken(user db.User) (string, time.Time, error) {
	expireTime := time.Now().Add(time.Hour * JWT_EXPIRE_HOURS)
	claims := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":   user.Name,
			"admin": user.Admin,
			"exp":   expireTime.Unix(),
		},
	)
	token, err := claims.SignedString(JWT_SECRET)
	return token, expireTime, err
}

func SetTokenCookie(c *fiber.Ctx, user db.User) error {
	jwt, expireTime, err := getToken(user)
	if err != nil {
		return err
	}
	c.Cookie(&fiber.Cookie{
		Name:     JWT_COOKIE_NAME,
		Value:    jwt,
		Expires:  expireTime,
		HTTPOnly: true,
		Secure:   true,
	})
	return nil
}

func ClearTokenCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     JWT_COOKIE_NAME,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HTTPOnly: true,
		Secure:   true,
	})
}

func UnloadTokenCookie(c *fiber.Ctx) (JWTPayload, error) {
	cookie := c.Cookies(JWT_COOKIE_NAME)
	token, err := jwt.Parse(cookie, func(t *jwt.Token) (interface{}, error) {
		return JWT_SECRET, nil
	})
	if err != nil {
		return JWTPayload{}, fmt.Errorf("Error parsing jwt from request cookie")
	}
	payload := token.Claims.(jwt.MapClaims)
	return JWTPayload{
		Name:       payload["sub"].(string),
		Admin:      payload["admin"].(bool),
		ExpireTime: payload["exp"].(int64),
	}, nil
}
