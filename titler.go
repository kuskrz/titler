package main

import (
	"kus/krzysztof/titler/environment"
	"kus/krzysztof/titler/handlers"
	"kus/krzysztof/titler/httpclient"
	"kus/krzysztof/titler/logging"
	"net/http"
	"sort"
)

func main() {
	logging.InitLogging()
	environment.InitEnv()
	logging.Log(logging.ERROR, "MNM: "+"Start \xF0\x9F\xAA\xBF")
	logEnv(&environment.EnvVars)
	httpclient.InitClient()

	// request -> basicAuthMux -> mux -> registered HandleFunc
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.GetTag)

	basicAuthMux := basicAuthMiddleware(mux, environment.EnvVars["BASIC_USER"], environment.EnvVars["BASIC_PASS"])
	end_err := http.ListenAndServe(":8080", basicAuthMux)
	logging.Log(logging.ERROR, end_err.Error())
}

func basicAuthMiddleware(next http.Handler, username, password string) http.Handler {
	// HandlerFunc is a function type that implements ServeHTTP thus it is a http.Handler
	// HandlerFunc implements ServeHTTP function as a call to function which is its argument
	// HandlerFunc function is called when ServeHTTP is called
	// the function checks user and pass and calls wrapped handler (mux in this case)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func logEnv(env *map[string]string) {
	names := make([]string, 0, len(*env))
	for n := range *env {
		names = append(names, n)
	}
	sort.Strings(names)

	for _, n := range names {
		logging.Log(logging.ERROR, "ENV: "+n+"="+(*env)[n])
	}
}
