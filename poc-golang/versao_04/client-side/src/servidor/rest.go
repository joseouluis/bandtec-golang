package servidor

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func escreverArquivoJSON(arquivo string, jsonAlterado []byte) {
	if erro := ioutil.WriteFile(localArquivosHTMLeJSON+"/"+arquivo, jsonAlterado, 0666); erro != nil {
		log.Println(erro)
	}
}

func lerArquivoJSON(arquivo string) []byte {
	var dadosArquivoJSON []byte
	var erro error

	switch arquivo {
	case "default.json":
		dadosArquivoJSON, erro = ioutil.ReadFile(localArquivosHTMLeJSON + "/" + arquivo)
	case "geral.json.html":
		dadosArquivoJSON, erro = ioutil.ReadFile(localArquivosHTMLeJSON + "/" + arquivo)
	case "categoria.json.html":
		dadosArquivoJSON, erro = ioutil.ReadFile(localArquivosHTMLeJSON + "/" + arquivo)
	}
	verificarErro(erro)

	return dadosArquivoJSON
}

func alterarArquivosJSON(denuncias []DadosDasDenuncias, arquivo string, verificaPorRegiao bool) {

	dadosArquivoJSON := lerArquivoJSON("default.json")
	var continua = true

	for _, item := range denuncias {
		alteracao := []byte(item.Nome)

		if verificaPorRegiao == true {
			alteracao = []byte(item.Regiao)
			continua = false

			if strings.ToUpper(item.Nome) == strings.ToUpper(paginaSelecionada) {
				continua = true
			}
		}

		if continua == true {
			jsonAlterado := bytes.Replace(dadosArquivoJSON, []byte("Categoria"), alteracao, 1)
			escreverArquivoJSON(arquivo, jsonAlterado)

			dadosArquivoJSON = lerArquivoJSON(arquivo)

			jsonAlterado = bytes.Replace(dadosArquivoJSON, []byte("00"), []byte(item.Total), 1)
			escreverArquivoJSON(arquivo, jsonAlterado)

			dadosArquivoJSON = lerArquivoJSON(arquivo)
		}
	}
}

func requisitarDados(url string) []DadosDasDenuncias {
	requisicao, erro := http.Get(url)
	verificarErro(erro)

	corpoDaRequisicao, erro := ioutil.ReadAll(requisicao.Body)
	verificarErro(erro)

	var denuncias []DadosDasDenuncias
	// coverte de json para struct
	json.Unmarshal(corpoDaRequisicao, &denuncias)

	return denuncias
}

func atualizarArquivosJSON() {
	log.Printf("atualiza arquivo JSON")

	denuncias := requisitarDados(urlTodasDenuncias)
	alterarArquivosJSON(denuncias, "geral.json.html", false)

	denunciasPorRegiao := requisitarDados(urlTodasDenunciasPorRegiao)
	alterarArquivosJSON(denunciasPorRegiao, "categoria.json.html", true)

	atualizarArquivosWeb()
}
