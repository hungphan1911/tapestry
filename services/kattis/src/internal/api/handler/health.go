package handler

import (
	"fmt"
	"io"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got a health check request\n")
	io.WriteString(w, "OK")
}