package main

import (
	"kus/krzysztof/titler/environment"
	"kus/krzysztof/titler/handlers"
	"kus/krzysztof/titler/httpclient"
	"kus/krzysztof/titler/logging"
	"net/http"
)

func main() {
	logging.InitLogging()
	environment.InitEnv()
	logging.Log(logging.ERROR, "MNM: "+"Start \xF0\x9F\xAA\xBF")
	for n, v := range environment.EnvVars {
		logging.Log(logging.ERROR, "ENV: "+n+"="+v)
	}
	httpclient.InitClient()

	http.HandleFunc("/", handlers.GetTag)
	end_err := http.ListenAndServe(":8080", nil)
	logging.Log(logging.ERROR, end_err.Error())
}
