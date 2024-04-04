package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int    `json:"created_at"`
}

type UserAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	CreatedAt   int    `json:"created_at"`
}

type CallBackResponse struct {
	Uuid   string `json:"uuid"`
	Status string `json:"status"`
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

	r.GET("/test", FractalURL)
	r.GET("/oauth/callback", CallBack)
	r.GET("/users", VerificationsByUserIds)
	r.Run(":8080")
}

func FractalURL(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, "https://app.next.fractal.id/authorize?client_id=ne6k3g1ZTyvpJwZfxTwRu0b9jEGfc4K4AIfrjFUary0&redirect_uri=https%3A%2F%2Fapi2.bethelnet.io%2Foauth%2Fcallback&response_type=code&scope=contact%3Aread%20verification.basic%3Aread%20verification.basic.details%3Aread%20verification.liveness%3Aread%20verification.liveness.details%3Aread&state=123")
}

func CallBack(ctx *gin.Context) {

	state := ctx.Query("state")
	code := ctx.Query("code")

	fmt.Println("state: ", state)
	fmt.Println("code: ", code)

	token := ExchangeCodeToAccessToken(code)
	uuid, status := GetUserDetails(token)

	fmt.Println("Uuid: ", uuid)
	fmt.Println("Status: ", status)

	// Create CID instances using urlGC
	kycResponse := CallBackResponse{
		Uuid:   uuid,
		Status: status,
	}

	ctx.JSON(200, kycResponse)

	ctx.Redirect(http.StatusFound, "https://testnet.bethelnet.io")

}

func GetUserDetails(token string) (string, string) {

	requestURL := "https://resource.next.fractal.id/v2/users/me"

	jsonBody := []byte(`{}`)

	bodyReader := bytes.NewReader(jsonBody)

	// Create new http "POST" request
	req, err := http.NewRequest(http.MethodGet, requestURL, bodyReader)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
	}

	fmt.Println("Token from get user: ", token)

	// Set Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer "+token))

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
	fmt.Printf("Response body: %s\n", responseBody)

	resp := GetBody(responseBody)

	// Extracting UID
	uid := resp["uid"].(string)
	fmt.Println("UID:", uid)

	var status string

	// Extracting verification cases and their statuses
	verificationCases := resp["verification_cases"].([]interface{})
	for _, v := range verificationCases {
		verificationCase := v.(map[string]interface{})
		status = verificationCase["status"].(string)
		fmt.Println("Status:", status)
	}

	return uid, status
}

func ExchangeCodeToAccessToken(code string) string {

	jsonBody := []byte(`{}`)

	bodyReader := bytes.NewReader(jsonBody)

	requestURL := fmt.Sprintf("https://auth.next.fractal.id/oauth/token?client_id=ne6k3g1ZTyvpJwZfxTwRu0b9jEGfc4K4AIfrjFUary0&client_secret=rMkPTgNPJh1VeEmNzjZBqE4_VrnIk2KLjWJNy2wGJeM&code=" + code + "&grant_type=authorization_code&redirect_uri=https://api2.bethelnet.io/oauth/callback")

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

func VerificationsByUserIds(ctx *gin.Context) {

	jsonBody := []byte(`{}`)

	bodyReader := bytes.NewReader(jsonBody)

	requestURL := "https://auth.next.fractal.id/oauth/token?client_id=ne6k3g1ZTyvpJwZfxTwRu0b9jEGfc4K4AIfrjFUary0&client_secret=rMkPTgNPJh1VeEmNzjZBqE4_VrnIk2KLjWJNy2wGJeM&scope=client.stats:read&grant_type=client_credentials"

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

	var response UserAccessTokenResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		fmt.Println("Error:", err)
	}

	resp := GetAllUsers(response.AccessToken)

	ctx.JSON(200, resp)

}

func GetAllUsers(token string) string {

	requestURL := "https://resource.next.fractal.id/v2/stats/user-verifications"

	jsonBody := []byte(`{}`)

	bodyReader := bytes.NewReader(jsonBody)

	// Create new http "POST" request
	req, err := http.NewRequest(http.MethodGet, requestURL, bodyReader)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
	}

	fmt.Println("Token from get user: ", token)

	// Set Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer "+token))

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
	fmt.Printf("Response body: %s\n", responseBody)

	return string(responseBody)
}

func GetBody(res []byte) map[string]interface{} {

	//pass to string body
	sRes := string(res)

	// Replace backslashes
	sRes = strings.Replace(sRes, "\\", "", -1)

	// Decode JSON
	var data map[string]interface{}
	err := json.Unmarshal([]byte(sRes), &data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}

	return data

}
