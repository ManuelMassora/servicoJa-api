package e2e_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFluxoVagaRembolso(t *testing.T) {
	// 1. Registrar um cliente
	email := "cliente_vaga_refund@test.com"
	telefone := "841112277"
	senha := "password123"
	registerUser("CLIENTE", email, telefone, senha)

	// 2. Login
	token := loginUser(telefone, senha)
	assert.NotEmpty(t, token)

	// 3. Criar uma vaga com preço premium (para triggar reembolso)
	fields := map[string]string{
		"titulo":             "Vaga Premium com Reembolso",
		"descricao":          "Vaga que será cancelada com reembolso M-Pesa",
		"localizacao":        "Maputo Centro",
		"latitude":           "-25.9692",
		"longitude":          "32.5732",
		"preco":              "10000.0",
		"urgente":            "true",
		"telefone_pagamento": "841112299",
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

	// 4. Cancelar a VAGA via /vagas/cancelar
	// Isso trigga o reembolso automático via M-Pesa B2C
	w = doRequest("POST", fmt.Sprintf("/vagas/cancelar/%d", vagaID), nil, map[string]string{
		"Authorization": "Bearer " + token,
	})
	assert.Equal(t, http.StatusOK, w.Code)

	// 5. Verificar que a vaga foi cancelada
	w = doRequest("GET", fmt.Sprintf("/vagas/%d", vagaID), nil, map[string]string{
		"Authorization": "Bearer " + token,
	})
	assert.Equal(t, http.StatusOK, w.Code)

	var getVagaResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &getVagaResp)
	assert.Equal(t, "CANCELADO", getVagaResp["status"])

	// 6. Verificar que o perfil do usuário mostra contador de cancelamentos aumentado
	w = doRequest("GET", "/usuario/perfil", nil, map[string]string{
		"Authorization": "Bearer " + token,
	})
	assert.Equal(t, http.StatusOK, w.Code)

	var userResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &userResp)
	cancellationCounter := int(userResp["cancelamento_contador"].(float64))
	assert.GreaterOrEqual(t, cancellationCounter, 1, "Cancellation counter should have increased")

	// 7. Verificar que a vaga não aparece mais nas vagas ativas do cliente
	w = doRequest("GET", "/vagas/cliente", nil, map[string]string{
		"Authorization": "Bearer " + token,
	})
	assert.Equal(t, http.StatusOK, w.Code)

	var listResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &listResp)
	data := listResp["data"].([]interface{})

	// Verificar que a vaga não aparece como ativa na listagem
	found := false
	for _, item := range data {
		v := item.(map[string]interface{})
		if uint(v["id"].(float64)) == vagaID {
			status := v["status"].(string)
			if status != "CANCELADO" {
				found = true
			}
		}
	}
	assert.False(t, found, "Vaga cancelada com reembolso não deve aparecer nas vagas ativas")

	// 8. Simulação: Verificar que a transação de reembolso foi criada
	// (Isso depende de haver um endpoint GET /transacoes/:id ou similar)
	// Para este teste, apenas verificamos que o cancelamento foi bem-sucedido
}
