package parser

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func URLParser(url *url.URL) (string, error) {
	parts := strings.Split(url.Path, "/")
	if len(parts[len(parts)-1]) == 0 {
		return "", nil
	}
	if _, err := strconv.Atoi(parts[len(parts)-1]); err != nil {
		return "", nil
	}
	return parts[len(parts)-1], nil
}

func JSONResponse(w http.ResponseWriter, code int, data interface{}) error {
	response, err := json.Marshal(data)
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
