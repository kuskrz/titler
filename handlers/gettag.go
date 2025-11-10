package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"kus/krzysztof/titler/environment"
	"kus/krzysztof/titler/logging"
	"kus/krzysztof/titler/requests"
	"kus/krzysztof/titler/responses"

	"golang.org/x/net/html"
)

func GetTag(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logging.Log(logging.DEBUG, "GetTag handler")

	var input requests.GetTag
	var response responses.ResponseData
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logging.Log(logging.ERROR, "GetTag cannot parse input")
		response = responses.ResponseData{Status: "failure", Tag: environment.EnvVars["TAG"], TagValue: "GetTag cannot parse input"}
		response.SendResponse(w)
		return
	}
	logging.Logf(logging.DEBUG, "Request from: %s %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)
	logging.Log(logging.DEBUG, input.Url)

	resp, err := http.Get(input.Url)
	if err != nil {
		logging.Log(logging.ERROR, err.Error())
		response = responses.ResponseData{Status: "failure", Tag: environment.EnvVars["TAG"], TagValue: err.Error()}
		response.SendResponse(w)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Log(logging.ERROR, err.Error())
		response = responses.ResponseData{Status: "failure", Tag: environment.EnvVars["TAG"], TagValue: err.Error()}
		response.SendResponse(w)
		return
	}

	var isRightToken bool
	tokenizer := html.NewTokenizer(strings.NewReader(string(body)))
MAINLOOP:
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			logging.Log(logging.ERROR, "cannot find HTML text token")
			response = responses.ResponseData{Status: "failure", Tag: environment.EnvVars["TAG"], TagValue: "cannot find HTML text token"}
			break MAINLOOP
		case html.StartTagToken:
			token := tokenizer.Token()
			isRightToken = token.Data == environment.EnvVars["TAG"]
		case html.TextToken:
			if isRightToken {
				token := tokenizer.Token()
				response = responses.ResponseData{Status: "success", Tag: environment.EnvVars["TAG"], TagValue: token.String()}
				break MAINLOOP
			}
			isRightToken = false
		default:
			isRightToken = false
		}
	}
	response.SendResponse(w)
}
