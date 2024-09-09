package handlers

import (
	// "fmt"
	"encoding/json"
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
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	modoCambioParam := r.URL.Query().Get("cotacao")
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

func PegaCotacao(modoCambio string) (*AweSomeApi, error) {

	resp, error := http.Get("https://economia.awesomeapi.com.br/json/last/" + modoCambio)
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}
	var c AweSomeApi
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}
	return &c, nil
}
