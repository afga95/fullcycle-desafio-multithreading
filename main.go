package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type BrasilCEP struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {

	c1 := make(chan BrasilCEP)
	c2 := make(chan ViaCEP)
	//var brasilCEP BrasilCEP
	//var viaCEP ViaCEP

	// Faz a requisição para a API do Brasil CEP
	go BuscaBrasilCEP(c1, "89010025")
	// Faz a requisição para a API do Via CEP
	go BuscaViaCEP(c2, "89010025")

	// Espera a resposta de uma das APIs
	select {
	case brasilCEP := <-c1:
		log.Println("Brasil CEP:", brasilCEP)
	case viaCEP := <-c2:
		log.Println("Via CEP:", viaCEP)
	case <-time.After(time.Second):
		log.Println("Timeout")
	}
}

func BuscaBrasilCEP(c1 chan<- BrasilCEP, cep string) {
	//time.Sleep(time.Second * 2)
	// Faz a requisição para a API do Brasil CEP
	resp, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Lê a resposta da API e converte para um slice de bytes
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Converte o slice de bytes para uma struct
	var data BrasilCEP
	err = json.Unmarshal(res, &data)
	if err != nil {
		log.Fatal(err)
	}

	// Retorna a struct com o CEP
	c1 <- data
}

func BuscaViaCEP(c2 chan<- ViaCEP, cep string) {
	//time.Sleep(time.Second * 2)
	// Faz a requisição para a API do Via CEP
	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Lê a resposta da API e converte para um slice de bytes
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Converte o slice de bytes para uma struct
	var data ViaCEP
	err = json.Unmarshal(res, &data)
	if err != nil {
		log.Fatal(err)
	}

	// Retorna a struct com o CEP
	c2 <- data

}
