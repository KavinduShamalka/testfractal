package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

/*
	{
	    "access_token": "RccWJCQxFmJ0Zg3O0-0C5wPuT2uV39k9BaWTUvz4nzI",
	    "token_type": "Bearer",
	    "expires_in": 7200,
	    "refresh_token": "PPync9yNKbdV41jnYjzsAmVfiVf2EVtpuXK6dPipL-Q",
	    "scope": "contact:read verification.basic:read verification.basic.details:read verification.liveness:read verification.liveness.details:read",
	    "created_at": 1712127620
	}
*/
type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int    `json:"created_at"`
}

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
	userDetails := GetUserDetails(token)

	ctx.JSON(http.StatusOK, gin.H{
		"state":       state,
		"code":        code,
		"token":       token,
		"userDetails": userDetails,
	})

}

func GetUserDetails(token string) []byte {

	requestURL := "https://resource.next.fractal.id/v2/users/me"

	jsonBody := []byte(`{}`)

	bodyReader := bytes.NewReader(jsonBody)

	// Create new http "POST" request
	req, err := http.NewRequest(http.MethodGet, requestURL, bodyReader)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
	}

	// Set Headers
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", fmt.Sprintf("Basic "+token))

	// Create client
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	// Make Http request and get response
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		fmt.Printf("Internal Server Error: %v\n", http.StatusInternalServerError)
	}

	defer res.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		fmt.Printf("Internal Server Error: %v\n", http.StatusInternalServerError)
	}

	// Print the response body
	fmt.Printf("%s\n", responseBody)

	return responseBody
}

func ExchangeCodeToAccessToken(code string) string {

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

	var response AccessTokenResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return response.AccessToken
}
