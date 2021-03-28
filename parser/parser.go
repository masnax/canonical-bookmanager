package parser

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

func URLParser(url *url.URL) []string {
	parts := strings.Split(url.Path, "/")
	if len(parts) < 3 {
		return nil
	}
	return parts
}

func JSONResponse(w http.ResponseWriter, code int, data interface{}) error {
	out := map[string]interface{}{
		"status-code": code,
		"status":      http.StatusText(code),
		"data":        data,
	}
	response, err := json.Marshal(out)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	return nil
}

func ErrorResponse(w http.ResponseWriter, code int, msg string) error {
	return JSONResponse(w, code, map[string]string{"error": msg})
}
