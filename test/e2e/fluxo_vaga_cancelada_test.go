package e2e_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFluxoVagaCancelada(t *testing.T) {
	// 1. Registrar um cliente
	email := "cliente_vaga_cancel@test.com"
	telefone := "841112266"
	senha := "password123"
	registerUser("CLIENTE", email, telefone, senha)

	// 2. Login
	token := loginUser(telefone, senha)
	assert.NotEmpty(t, token)

	// 3. Criar uma vaga
	fields := map[string]string{
		"titulo":             "Vaga de Teste Cancelável",
		"descricao":          "Esta vaga será cancelada",
		"localizacao":        "Maputo Periférico",
		"latitude":           "-25.9692",
		"longitude":          "32.5732",
		"preco":              "2000.0",
		"urgente":            "false",
		"telefone_pagamento": "841234588",
	}

	w := doMultipartRequest("POST", "/vagas", fields, token)
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	vagaID := uint(resp["id"].(float64))
	assert.NotZero(t, vagaID)

	// 3.5. Simular confirmação de pagamento para a vaga
	doRequest("GET", fmt.Sprintf("/pagamento/simular-query?ref=VAG-%d", vagaID), nil, nil)

	// 4. Verificar que a vaga foi criada COM status DISPONIVEL
	w = doRequest("GET", fmt.Sprintf("/vagas/%d", vagaID), nil, map[string]string{
		"Authorization": "Bearer " + token,
	})
	assert.Equal(t, http.StatusOK, w.Code)

	var getVagaResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &getVagaResp)
	initialStatus := getVagaResp["status"].(string)
	assert.NotEqual(t, "CANCELADO", initialStatus, "Vaga não deve estar cancelada logo após criação")

	// 5. Cancelar a vaga
	w = doRequest("POST", fmt.Sprintf("/vagas/cancelar/%d", vagaID), nil, map[string]string{
		"Authorization": "Bearer " + token,
	})
	assert.Equal(t, http.StatusOK, w.Code)

	var cancelResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &cancelResp)
	assert.Equal(t, "Vaga cancelada com sucesso", cancelResp["message"])

	// 6. Verificar que o status da vaga é agora CANCELADO
	w = doRequest("GET", fmt.Sprintf("/vagas/%d", vagaID), nil, map[string]string{
		"Authorization": "Bearer " + token,
	})
	assert.Equal(t, http.StatusOK, w.Code)

	var getVagaCanceledResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &getVagaCanceledResp)
	assert.Equal(t, "CANCELADO", getVagaCanceledResp["status"])

	// 7. Verificar que a vaga não aparece mais na listagem de vagas ativas
	w = doRequest("GET", "/vagas/cliente", nil, map[string]string{
		"Authorization": "Bearer " + token,
	})
	assert.Equal(t, http.StatusOK, w.Code)

	var listResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &listResp)
	data := listResp["data"].([]interface{})

	// Verificar que a vaga não aparece na listagem (ou está marcada como deleted)
	for _, item := range data {
		v := item.(map[string]interface{})
		if uint(v["id"].(float64)) == vagaID {
			status := v["status"].(string)
			assert.Equal(t, "CANCELADO", status, "Vaga cancelada não deve aparecer como ativa")
			break
		}
	}
}
