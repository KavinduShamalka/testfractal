package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var conf *oauth2.Config

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://*", "https://*", "*", "https://testnet.bethelnet.io"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.GET("/test", Test)
	r.GET("/oauth/callback", CallBack)
	r.Run(":8080")
}

func Test(c *gin.Context) {

	conf = &oauth2.Config{
		ClientID: "ne6k3g1ZTyvpJwZfxTwRu0b9jEGfc4K4AIfrjFUary0",
	}

	redirectURI := "http://localhost:8080/callback"
	encodedRedirectURI := url.QueryEscape(redirectURI)

	val := url.Values{}
	val.Set("client_id", conf.ClientID) // Add client_id to query params
	val.Set("redirect_uri", encodedRedirectURI)
	val.Set("response_type", "code")
	val.Set("state", "xyzABC123")

	authorizationURL := fmt.Sprintf("https://app.next.fractal.id/authorize?%s", val.Encode())

	fmt.Println("Url: ", authorizationURL)

	c.Redirect(http.StatusFound, "https://app.next.fractal.id/authorize?client_id=ne6k3g1ZTyvpJwZfxTwRu0b9jEGfc4K4AIfrjFUary0&redirect_uri=http%3A%2F%2Flocalhost%3A8090%2Foauth%2Fcallback&response_type=code&scope=contact%3Aread%20verification.basic%3Aread%20verification.basic.details%3Aread%20verification.liveness%3Aread%20verification.liveness.details%3Aread")
}

func CallBack(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	fmt.Println("state: ", state)
	fmt.Println("code: ", code)

	c.JSON(http.StatusOK, gin.H{
		"state": state,
		"code":  code,
	})
}
