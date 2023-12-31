package api

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/quevivasbien/bird-game/db"
)

var JWT_SECRET = []byte(os.Getenv("BIRD_JWT_SECRET"))

const JWT_EXPIRE_HOURS = 12
const JWT_COOKIE_NAME = "jwt_token"

type JWTPayload struct {
	Name       string `json:"name"`
	Admin      bool   `json:"admin"`
	ExpireTime int64  `json:"expireTime"`
}

type MissingToken struct{}

func (t MissingToken) Error() string {
	return "Request is missing JWT authentication cookie"
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

func SetTokenCookie(c *fiber.Ctx, user db.User) (JWTPayload, error) {
	jwt, expireTime, err := getToken(user)
	if err != nil {
		return JWTPayload{}, err
	}
	c.Cookie(&fiber.Cookie{
		Name:     JWT_COOKIE_NAME,
		Value:    jwt,
		Expires:  expireTime,
		HTTPOnly: true,
		Secure:   true,
		Path:     "/",
	})
	return JWTPayload{
		Name:       user.Name,
		Admin:      user.Admin,
		ExpireTime: expireTime.Unix(),
	}, nil
}

func ClearTokenCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:   JWT_COOKIE_NAME,
		Value:  "",
		MaxAge: -1,
	})
}

func UnloadTokenCookie(c *fiber.Ctx) (JWTPayload, error) {
	cookie := c.Cookies(JWT_COOKIE_NAME)
	if cookie == "" {
		return JWTPayload{}, MissingToken{}
	}
	token, err := jwt.Parse(cookie, func(t *jwt.Token) (interface{}, error) {
		return JWT_SECRET, nil
	})
	if err != nil {
		return JWTPayload{}, fmt.Errorf("Error parsing jwt from request cookie: %v", err)
	}
	payload := token.Claims.(jwt.MapClaims)
	name := payload["sub"].(string)
	if name == "" {
		return JWTPayload{}, fmt.Errorf("Empty name in parsed JWT payload")
	}
	return JWTPayload{
		Name:       name,
		Admin:      payload["admin"].(bool),
		ExpireTime: int64(payload["exp"].(float64)),
	}, nil
}
