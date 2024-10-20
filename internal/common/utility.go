package common

import (
	"encoding/json"
	"log/slog"
	"math"
	"net/http"
	"strconv"
)

func Ternary(condition bool, valueIfTrue, valueIfFalse interface{}) interface{} {
	if condition {
		return valueIfTrue
	}
	return valueIfFalse
}

func HandleErrorRes(writer http.ResponseWriter, errMap map[string]string) {
	HandleRes(writer, errMap, http.StatusBadRequest)
}

func HandleUnAuthRes(writer http.ResponseWriter, errMap map[string]string) {
	HandleRes(writer, errMap, http.StatusUnauthorized)
}

func HandleRes(writer http.ResponseWriter, errMap map[string]string, status int) {
	writer.WriteHeader(status)
	if errMap != nil && len(errMap) > 0 {
		jsonData, _ := json.Marshal(errMap)
		writer.Write(jsonData)
	}
}

func SendOkRes(writer http.ResponseWriter, errMap map[string]string) {
	HandleRes(writer, errMap, http.StatusOK)
}

func ReadReqPayload(writer http.ResponseWriter, request *http.Request, payload interface{}) bool {
	err := json.NewDecoder(request.Body).Decode(payload)
	if err != nil {
		slog.Error(err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		return false
	}
	return true
}

func TwoDigitPrecision(val float64) float64 {
	return math.Floor(val*100) / 100
}

func FetchPathVariable(writer http.ResponseWriter, request *http.Request, pathVariable string) (int, error) {
	variable, err := strconv.Atoi(request.PathValue(pathVariable))

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		slog.Debug("cannot read " + pathVariable + " from " + request.PathValue(pathVariable))
	}

	return variable, nil
}
