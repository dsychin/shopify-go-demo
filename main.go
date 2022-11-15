package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.AddTrailingSlash())

	// Routes
	e.GET("/auth/", auth)
	e.GET("/callback/", authCallback)
	e.GET("/dashboard/", dashboard)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func auth(c echo.Context) error {
	shop := c.QueryParam("shop")

	permissionUrl := fmt.Sprintf("https://%s/admin/oauth/authorize?client_id=%s&scope=read_products&redirect_uri=%s&state=mynonce", shop, os.Getenv("SHOPIFY_CLIENT_ID"), "http://localhost:8080/callback/")

	return c.Redirect(http.StatusTemporaryRedirect, permissionUrl)
}

func authCallback(c echo.Context) error {
	return c.Redirect(http.StatusTemporaryRedirect, "/dashboard/")
}

func dashboard(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the dashboard!")
}
