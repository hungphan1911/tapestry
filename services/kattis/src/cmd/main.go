package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"kattis/src/internal/api"
)

func main() {
	fmt.Println("Serving server")

	err := http.ListenAndServe(":8080", api.NewRouter())
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed\n")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
