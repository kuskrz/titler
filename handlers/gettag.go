package handlers

import (
	"encoding/json"
	"kus/krzysztof/titler/logging"
	"kus/krzysztof/titler/requests"
	"net/http"
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

	res, err := http.Get(input.Url)
	if err == nil {
		res.Body.Close()
	} else {
		logging.Log(logging.ERROR, err.Error())
	}
}
