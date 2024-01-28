package auth

import (
	"GoBaatcheet/constants"
	"net/http"

	authorizer "github.com/authorizerdev/authorizer-go"
)

const authorizerServer = "http://localhost:8082/"       // Todo: Move this to centralized configuration server
const clientId = "c87ad9f9-e076-429f-b175-777e73570a9b" // Todo: Move this to centralized configuration server

func Authenticate(request *http.Request) bool {
	defaultHeaders := map[string]string{}
	client, err := authorizer.NewAuthorizerClient(clientId, authorizerServer, constants.EmptyStr, defaultHeaders)
	if err != nil {
		panic(err)
	}
	token := request.URL.Query().Get("access_token")
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
	return response.IsValid
}
