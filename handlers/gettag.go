package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"kus/krzysztof/titler/environment"
	"kus/krzysztof/titler/logging"
	"kus/krzysztof/titler/requests"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func GetTag(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logging.Log(logging.DEBUG, "GetTag handler")

	var input requests.GetTag
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logging.Log(logging.ERROR, "GetTag cannot parse input")
		return
	}
	logging.Logf(logging.DEBUG, "Request from: %s %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)
	logging.Log(logging.DEBUG, input.Url)

	resp, err := http.Get(input.Url)
	if err != nil {
		logging.Log(logging.ERROR, err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Log(logging.ERROR, err.Error())
		return
	}

	tokenizer := html.NewTokenizer(strings.NewReader(string(body)))
	var isRightToken bool
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			logging.Log(logging.ERROR, "cannot parse HTML token")
			return
		case html.StartTagToken:
			token := tokenizer.Token()
			isRightToken = token.Data == environment.EnvVars["TAG"]
		case html.TextToken:
			if isRightToken {
				token := tokenizer.Token()
				fmt.Println(token)
				return
			}
			isRightToken = false
		default:
			isRightToken = false
		}
	}
}
