package handlers

import (
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "H-Bem-vindo servidor Go consumidor cotação dolar!")
}
