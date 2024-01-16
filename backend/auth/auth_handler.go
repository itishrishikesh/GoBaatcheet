package auth

import (
	"log"
	"net/http"

	authorizer "github.com/authorizerdev/authorizer-go"
)

func Authenticate(request *http.Request) bool {
	defaultHeaders := map[string]string{}
	client, err := authorizer.NewAuthorizerClient("c87ad9f9-e076-429f-b175-777e73570a9b", "http://localhost:8082/", "", defaultHeaders)

	if err != nil {
		panic(err)
	}

	token := request.URL.Query().Get("access_token")

	log.Println("Access token", token)

	if token == "" {
		return false
	}

	response, err := client.ValidateJWTToken(&authorizer.ValidateJWTTokenInput{
		TokenType: authorizer.TokenTypeIDToken,
		Token:     token,
	})

	if err != nil {
		panic(err)
	}

	log.Println("Validated Token", response)

	return response.IsValid
}
