package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

const (
	RATE_LIMIT_QUOTA = 300 // requests per REDIS_RATE_LIMIT_RESET
)

type FilterRouter struct{}

func (h *FilterRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	currentQuota := pb.RedisUpdateQuota(now.Minute(), remoteIpAddress(r))
	if currentQuota >= RATE_LIMIT_QUOTA {
		retryAfter := 59 - now.Second()

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Retry-After", strconv.Itoa(retryAfter))
		w.WriteHeader(429)
		w.Write([]byte(fmt.Sprintf("{\"limit\":\"%d\",\"retryAfter\":\"%d\"}", RATE_LIMIT_QUOTA, retryAfter)))
	} else {
		v1Router.ServeHTTP(w, r)
	}
}

func remoteIpAddress(r *http.Request) string {
	ipAddress := r.Header.Get("X-Real-Ip")
	if len(ipAddress) == 0 {
		ipAddress = strings.Split(r.RemoteAddr, ":")[0]
	}

	return ipAddress
}

func shelterIdFromRequest(r *http.Request) (string, error) {
	shelterId := parameterFromRequest("shelterId", r)
	if err := validateShelterId(shelterId); err != nil {
		return "", err
	}

	return shelterId, nil
}

func updateIdFromRequest(shelterId string, r *http.Request) (string, error) {
	updateId := parameterFromRequest("updateId", r)
	if err := validateUpdateId(shelterId, updateId); err != nil {
		return "", err
	}

	return updateId, nil
}

func shelterIdUpdateIdFromRequest(r *http.Request) (string, string, error) {
	shelterId, err := shelterIdFromRequest(r)
	if err != nil {
		return "", "", err
	}

	updateId, err := updateIdFromRequest(shelterId, r)
	if err != nil {
		return "", "", err
	}

	return shelterId, updateId, nil
}

func animalIdFromRequest(shelterId, updateId string, r *http.Request) (string, error) {
	animalId := parameterFromRequest("animalId", r)
	if err := validateAnimalId(shelterId, updateId, animalId); err != nil {
		return "", err
	}

	return animalId, nil
}

func shelterIdUpdateIdAnimalIdFromRequest(r *http.Request) (string, string, string, error) {
	shelterId, err := shelterIdFromRequest(r)
	if err != nil {
		return "", "", "", err
	}

	updateId, err := updateIdFromRequest(shelterId, r)
	if err != nil {
		return "", "", "", err
	}

	animalId, err := animalIdFromRequest(shelterId, updateId, r)
	if err != nil {
		return "", "", "", err
	}

	return shelterId, updateId, animalId, nil
}

func parameterFromRequest(name string, r *http.Request) string {
	param := mux.Vars(r)[name]
	if len(param) > 0 {
		return strings.ToLower(param)
	}

	return ""
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func response(w http.ResponseWriter, i interface{}) {
	w.Header().Add("Cache-Control", "no-cache")
	responseStatus(w, http.StatusOK, i)
}

func cachedResponse(w http.ResponseWriter, i interface{}) {
	w.Header().Add("Cache-Control", "max-age=2592000")
	responseStatus(w, http.StatusOK, i)
}

func responseCreated(w http.ResponseWriter, i interface{}) {
	responseStatus(w, http.StatusCreated, i)
}

func responseNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func responseStatus(w http.ResponseWriter, status int, i interface{}) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(i); err != nil {
		var bytes bytes.Buffer
		bytes.WriteString("{ \"error\": \"Response serialization failed! Error: '")
		bytes.WriteString(err.Error())
		bytes.WriteString("'\" }")
		w.Write(bytes.Bytes())
		return
	}
}

func unmarshalRequestBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(&v)
}

func conflict(w http.ResponseWriter, err error) {
	responseStatus(w, http.StatusConflict, ErrorResponse{err.Error()})
}

func unauthorized(w http.ResponseWriter, err error) {
	responseStatus(w, http.StatusUnauthorized, ErrorResponse{err.Error()})
}

func badRequest(w http.ResponseWriter, err error) {
	responseStatus(w, http.StatusBadRequest, ErrorResponse{err.Error()})
}

func internalServerError(w http.ResponseWriter, err error) {
	responseStatus(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
}
