package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFluxoVagaSuccess(t *testing.T) {
	// 1. Registrar um Cliente
	email := "cliente_vaga@test.com"
	telefone := "841112233"
	senha := "password123"
	registerUser("CLIENTE", email, telefone, senha)
	// 2. Login
	token := loginUser(telefone, senha)
	assert.NotEmpty(t, token)

	// 3. Criar uma Vaga
	fields := map[string]string{
		"titulo":             "Vaga de Teste E2E",
		"descricao":          "Descrição da vaga de teste",
		"localizacao":        "Maputo",
		"latitude":           "-25.9692",
		"longitude":          "32.5732",
		"preco":              "1500.0",
		"urgente":            "true",
		"telefone_pagamento": "841234567",
	}

	w := doMultipartRequest("POST", "/vagas", fields, token)
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Vaga criada com sucesso", resp["message"])
	vagaID := uint(resp["id"].(float64))
	assert.NotZero(t, vagaID)

	// 5. Listar vagas do cliente para confirmar
	w = doRequest("GET", "/vagas/cliente", nil, map[string]string{"Authorization": "Bearer " + token})
	assert.Equal(t, http.StatusOK, w.Code)

	var listResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &listResp)
	data := listResp["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 1)

	found := false
	for _, item := range data {
		v := item.(map[string]interface{})
		if uint(v["id"].(float64)) == vagaID {
			found = true
			assert.Equal(t, "Vaga de Teste E2E", v["titulo"])
			break
		}
	}
	assert.True(t, found, "Vaga não encontrada na listagem pelo ID retornado")
}
