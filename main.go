package main

import (
	"net/http"
	"time"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if users[username] != password {
		return echo.ErrUnauthorized
	}

	claims := &jwtCustomClaims{
		username,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)

	return c.String(http.StatusOK, "Welcome "+claims.Name+"!")
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/api/login", login)

	e.GET("/api", accessible)

	r := e.Group("/api/restricted")

	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte("secret"),
	}

	r.Use(middleware.JWTWithConfig(config))
	r.GET("", restricted)

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.HasPrefix(c.Request().URL.Path, "/api") {
				return next(c)
			}

			file := "lbjwt/dist" + c.Request().URL.Path
			if _, err := os.Stat(file); os.IsNotExist(err) {
				file = "lbjwt/dist/index.html"
			}

			return c.File(file)
		}
	})

	e.Start(":1323")
}