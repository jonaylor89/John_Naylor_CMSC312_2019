
package server

import (
	"fmt"
	"net/http"

	"github.com/jonaylor89/John_Naylor_CMSC312_2019/BE/kernel"
)

var kern *kernel.Kernel

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func repl(w http.ResponseWriter, r *http.Request) { 
	fmt.Fprint(w, kern)
}

// StartServer : state http server
func StartServer(k *kernel.Kernel) {

	kern = k

	http.HandleFunc("/", handler)
	http.HandleFunc("/repl", repl)
    http.ListenAndServe(":8080", nil)
}