package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
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
	vars := mux.Vars(r)
	modoCambioParam, ok := vars["cambio"]
	if !ok || modoCambioParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cotacao, err := PegaCotacao(modoCambioParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacao)

}

func PegaCotacao(modoCambio string) (*map[string]AweSomeApi, error) {

	resp, err := http.Get("https://economia.awesomeapi.com.br/json/last/" + modoCambio)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Usando mapa, pois o cambio pode mudar
	var resultado map[string]AweSomeApi

	err = json.Unmarshal(body, &resultado)
	if err != nil {
		return nil, err
	}
	return &resultado, nil
}
