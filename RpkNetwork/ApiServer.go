package RpkNetwork

import (
	"fmt"
	"net/http"
)

func StartApiServer(port string) {
	http.HandleFunc("/", apiHandler)
	fmt.Println("API server listening on port", port)
	http.ListenAndServe(":"+port, nil)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from the RESTful API!")
}
