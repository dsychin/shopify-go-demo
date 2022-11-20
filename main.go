package main

import (
	"fmt"
	"net/http"
	"net/url"
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
	e.Static("/", "static")
	e.GET("/auth/", auth)
	e.GET("/callback/", authCallback)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func auth(c echo.Context) error {
	shop := c.QueryParam("shop")

	err := verifyHmac(c.QueryParam("hmac"))
	if err != nil {
		return err
	}

	permissionUrl, err := url.Parse(fmt.Sprintf("https://%s/admin/oauth/authorize", shop))
	if err != nil {
		return err
	}
	callbackUrl := fmt.Sprintf("https://%s/callback/", c.Request().Host)
	nonce := "mynonce" // TODO: this should be generated and checked against in the callback
	q := permissionUrl.Query()
	q.Set("scope", "read_products")
	q.Set("client_id", os.Getenv("SHOPIFY_CLIENT_ID"))
	q.Set("redirect_uri", callbackUrl)
	q.Set("nonce", nonce)
	permissionUrl.RawQuery = q.Encode()

	return c.Redirect(http.StatusTemporaryRedirect, permissionUrl.String())
}

func authCallback(c echo.Context) error {
	err := verifyHmac(c.QueryParam("hmac"))
	if err != nil {
		return err
	}

	// TODO: set up store's account on your database

	u, err := url.Parse("/")
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("host", c.QueryParam("host"))
	u.RawQuery = q.Encode()

	return c.Redirect(http.StatusTemporaryRedirect, u.String())
}

func verifyHmac(hmac string) error {
	// TODO: You will need to implement the HMAC verification logic here
	// Refer to the documentation below
	// https://shopify.dev/apps/auth/oauth/getting-started#step-2-verify-the-installation-request
	return nil
}
