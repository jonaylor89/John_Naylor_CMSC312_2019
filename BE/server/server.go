
package server

import (
	"fmt"
	"net/http"

	// "github.com/jonaylor89/John_Naylor_CMSC312_2019/BE/memory"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

// StartServer : state http server
func StartServer() {
	http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}