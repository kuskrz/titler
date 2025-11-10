package responses

import (
	"encoding/json"
	"net/http"

	"kus/krzysztof/titler/logging"
)

type ResponseData struct {
	Status   string
	Tag      string
	TagValue string
}

func (r *ResponseData) SendResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(r); err != nil {
		logging.Log(logging.ERROR, "cannot parse output")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
