package responseHandler

import (
	"bitcoin-wallet/models"
	"encoding/json"
	"net/http"
)

func ReturnJsonMessage(w http.ResponseWriter, status string, message string, obj ...interface{}) {
	w.Header().Set("Content-Type", "application/json")

	response := models.Response{
		Status: status,
	}
	if message != "" {
		response.Message = message
	}

	if len(obj) > 0 {
		response.Object = obj
	}

	json.NewEncoder(w).Encode(response)
}
