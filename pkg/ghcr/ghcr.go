package ghcr

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

func GetGhcrToken(githubAppID, githubInstallationID, privateKeyPath string) string {
	// Load the private key from a PEM file
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatalf("Failed to read private key file: %v", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Create a new JWT token for the GitHub App
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Issuer:    githubAppID, // Replace with your GitHub App's App ID
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)), // Token valid for 10 minutes
	})

	// Sign the token with the private key
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
	}

	// Setup HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.github.com/app/installations/"+githubInstallationID+"/access_tokens", nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+signedToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read and print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var response TokenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Failed to unmarshal response: %v", err)
	}

	return response.Token
}

func CreateDockerConfigJson(usename, token string) string {
	auth := usename + ":" + token
	authB64 := base64.StdEncoding.EncodeToString([]byte(auth))
	return `{"auths": {"ghcr.io": {"auth": "` + authB64 + `"}}}`
}
