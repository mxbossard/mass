package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type (
	ServerError struct {
		Status     string
		StatusCode int
		Message    string
		Path       string
	}
)

func (e ServerError) Error() string {
	return fmt.Sprintf("Server error %s on path %s : %s !", e.Status, e.Path, e.Message)
}

func ReadServerError(r *http.Response) error {
	if r.StatusCode == 200 {
		return nil
	}
	se := ServerError{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(body, &se)
	se.Status = r.Status
	se.StatusCode = r.StatusCode
	return se
}
