package auth

import (
	"net/http"
	"strings"

	authorizer "github.com/authorizerdev/authorizer-go"
)

func authenticate(request *http.Request) bool {
	defaultHeaders := map[string]string{}
	client, err := authorizer.NewAuthorizerClient("c87ad9f9-e076-429f-b175-777e73570a9b", "http://localhost:8082/", "", defaultHeaders)

	if err != nil {
		panic(err)
	}

	header := request.Header.Get("Authorization")
	token := strings.Split(header, " ")

	if len(token) < 2 || token[1] == "" {
		return false
	}

	response, err := client.ValidateJWTToken(&authorizer.ValidateJWTTokenInput{
		TokenType: authorizer.TokenTypeIDToken,
		Token:     token[1],
	})

	if err != nil {
		panic(err)
	}

	return response.IsValid
}
