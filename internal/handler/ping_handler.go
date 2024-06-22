package handler

import "net/http"

func PingHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("{\"message\":\"pong\"}"))
}
