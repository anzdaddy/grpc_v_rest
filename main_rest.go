package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/husobee/vestigo"
	"github.com/sirupsen/logrus"
)

func mainREST(addr string, creds tlsCreds) *http.Server {
	r := vestigo.NewRouter()
	r.Post("/info", SetInfo)

	server := &http.Server{Addr: addr, Handler: r}
	go func() {
		if err := server.ListenAndServeTLS(creds.certFile, creds.keyFile); err != nil {
			logrus.Fatal(err)
		}
	}()
	return server
}

// SetInfo - Rest HTTP Handler
func SetInfo(w http.ResponseWriter, r *http.Request) {
	var (
		input    apiInput
		response apiResponse
	)

	// decode input
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	r.Body.Close()

	// validate input
	if err := validate(input); err != nil {
		response.Success = false
		response.Reason = err.Error()
		respBytes, _ := json.Marshal(response)

		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write(respBytes); err != nil {
			logrus.Error(err)
		}
		return
	}
	response.Success = true
	respBytes, _ := json.Marshal(response)

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(respBytes); err != nil {
		logrus.Error(err)
	}
}

type apiResponse struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason,omitempty"`
}

type apiInput struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Height int    `json:"height"`
}

// Validate - implementation of Validatable
func (ai apiInput) Validate() error {
	var err validationErrors
	if ai.Name == "" {
		err = append(err, errors.New("Name must be present"))
	}
	if ai.Age <= 0 {
		err = append(err, errors.New("Age must be real"))
	}
	if ai.Height <= 0 {
		err = append(err, errors.New("Height must be real"))
	}
	if len(err) == 0 {
		return nil
	}
	logrus.Infof("input: %#v", ai)
	return err
}
