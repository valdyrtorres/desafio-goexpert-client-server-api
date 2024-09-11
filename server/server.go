package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/usuario/desafio-go/handlers"
)

func main() {
	// Criar o ServeMux (multiplexador) para roteamento de URLs
	muxDesafio := mux.NewRouter()

	// Registrar a rota "/" para o handler HomeHandler
	muxDesafio.HandleFunc("/", handlers.HomeHandler)

	muxDesafio.HandleFunc("/cotacao/{cambio}", handlers.CotacaoHandler)

	muxDesafio.HandleFunc("/cotacao/full/{cambio}", handlers.CotacaoFullHandler)

	// Iniciar o servidor na porta 8080
	fmt.Println("Servidor rodando na porta 8080...")
	http.ListenAndServe(":8080", muxDesafio)
}
