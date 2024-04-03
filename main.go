package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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

func Test(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, "https://app.next.fractal.id/authorize?client_id=ne6k3g1ZTyvpJwZfxTwRu0b9jEGfc4K4AIfrjFUary0&redirect_uri=https%3A%2F%2Fapi2.bethelnet.io%2Foauth%2Fcallback&response_type=code&scope=contact%3Aread%20verification.basic%3Aread%20verification.basic.details%3Aread%20verification.liveness%3Aread%20verification.liveness.details%3Aread&state=123")
}

func CallBack(ctx *gin.Context) {
	state := ctx.Query("state")
	code := ctx.Query("code")

	fmt.Println("state: ", state)
	fmt.Println("code: ", code)

	token := ExchangeCodeToAccessToken(code)

	ctx.JSON(http.StatusOK, gin.H{
		"state": state,
		"code":  code,
		"token": token,
	})

}

func ExchangeCodeToAccessToken(code string) []byte {

	client_id := "ne6k3g1ZTyvpJwZfxTwRu0b9jEGfc4K4AIfrjFUary0"
	client_secret := "rMkPTgNPJh1VeEmNzjZBqE4_VrnIk2KLjWJNy2wGJeM"
	grant_type := "authorization_code"
	redirect_uri := "https://testnet.bethelnet.io/"

	jsonBody := []byte(`{}`)

	bodyReader := bytes.NewReader(jsonBody)

	requestURL := "https://auth.next.fractal.id/oauth/token?" + client_id + "&" + client_secret + "&" + code + "&" + grant_type + "&" + redirect_uri

	fmt.Println("Code: ", code)
	fmt.Println("Request URL: ", requestURL)

	// Create new http "POST" request
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)

	}

	defer res.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
	}

	return responseBody
}
