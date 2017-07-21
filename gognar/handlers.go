package gognar

import (
	"encoding/json"
	"io"
	"net/http"
)

func ReadJSON(reader io.ReadCloser, model interface{}) error {
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(model); err != nil {
		return err
	}
	defer reader.Close()
	return nil
}

// 200 RESPONSES
func ResponseJson(w http.ResponseWriter, code int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func Send(w http.ResponseWriter, response interface{}) {
	ResponseJson(w, http.StatusOK, response)
}

func Created(w http.ResponseWriter, response interface{}) {
	ResponseJson(w, http.StatusCreated, response)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// 400 RESPONSES
func NotFound(w http.ResponseWriter, err error) {
	Abort(w, http.StatusNotFound, "not_found", err.Error())
}

func BadRequest(w http.ResponseWriter, err error) {
	Abort(w, http.StatusBadRequest, "bad_request", err.Error())
}

// 500 RESPONSES
func InternalServerError(w http.ResponseWriter, err error) {
	Abort(w, http.StatusInternalServerError, "internal_server_error", err.Error())
}

func Abort(w http.ResponseWriter, status int, err string, message string) {
	ResponseJson(w, status, ResponseError{
		Status:  status,
		Errors:  err,
		Message: message,
	})
}
