package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AweSomeApi struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Nome       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func CotacaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	modoCambioParam := r.URL.Query().Get("cambio")
	if modoCambioParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cotacao, error := PegaCotacao(modoCambioParam)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacao)

	//fmt.Fprintln(w, "Chama a API de cotacao!")
}

func PegaCotacao(modoCambio string) (*map[string]AweSomeApi, error) {

	fmt.Printf("modoCambio: %s%s", "https://economia.awesomeapi.com.br/json/last/", modoCambio)

	resp, error := http.Get("https://economia.awesomeapi.com.br/json/last/" + modoCambio)
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}

	// Usando mapa, pois o cambio pode mudar
	var resultado map[string]AweSomeApi

	error = json.Unmarshal(body, &resultado)
	if error != nil {
		return nil, error
	}
	return &resultado, nil
}
