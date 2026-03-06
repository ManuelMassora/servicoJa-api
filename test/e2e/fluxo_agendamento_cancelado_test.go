package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFluxoAgendamentoCancelado(t *testing.T) {
	// 1. Setup: Criar admin, categoria e prestador
	adminEmail := "admin_ag_cancel@test.com"
	adminTelefone := "820000011"
	adminSenha := "admin123"
	registerUser("ADMIN", adminEmail, adminTelefone, adminSenha)
	adminToken := loginUser(adminTelefone, adminSenha)

	// Criar categoria
	catBody := map[string]interface{}{
		"nome":      "Limpeza Profissional",
		"descricao": "Serviços de limpeza profunda",
	}
	catJson, _ := json.Marshal(catBody)
	w := doRequest("POST", "/categoria", bytes.NewBuffer(catJson), map[string]string{
		"Authorization": "Bearer " + adminToken,
		"Content-Type":  "application/json",
	})
	assert.Equal(t, http.StatusCreated, w.Code)
	var catCreateResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &catCreateResp)
	catID := uint(catCreateResp["id"].(float64))

	// Registrar prestador
	prestEmail := "prestador_ag_cancel@test.com"
	prestTelefone := "820000012"
	prestSenha := "prestador123"
	registerUser("PRESTADOR", prestEmail, prestTelefone, prestSenha)
	prestToken := loginUser(prestTelefone, prestSenha)

	// Criar catálogo
	fields := map[string]string{
		"nome":         "Faxina Cancelável",
		"descricao":    "Serviço que será cancelado",
		"tipo_preco":   "fixo",
		"valor_fixo":   "2500.0",
		"categoria_id": fmt.Sprint(catID),
		"localizacao":  "Maputo",
		"latitude":     "-25.9692",
		"longitude":    "32.5732",
	}
	w = doMultipartRequest("POST", "/catalogo", fields, prestToken)
	assert.Equal(t, http.StatusCreated, w.Code)
	var catalogoCreateResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &catalogoCreateResp)
	catalogoID := uint(catalogoCreateResp["id"].(float64))

	// 2. Registrar e logar cliente
	cliEmail := "cliente_ag_cancel@test.com"
	cliTelefone := "820000013"
	cliSenha := "cliente123"
	registerUser("CLIENTE", cliEmail, cliTelefone, cliSenha)
	cliToken := loginUser(cliTelefone, cliSenha)

	// 3. Criar agendamento
	agFields := map[string]string{
		"detalhe":            "Preciso de uma faxina para cancelar depois",
		"id_catalogo":        fmt.Sprint(catalogoID),
		"datahora":           time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"localizacao":        "Casa de Teste",
		"latitude":           "-25.9000",
		"longitude":          "32.5000",
		"telefone_pagamento": "841112244",
	}
	w = doMultipartRequest("POST", "/agendamento", agFields, cliToken)
	assert.Equal(t, http.StatusCreated, w.Code)

	var agResult map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &agResult)
	agendamentoID := uint(agResult["id"].(float64))
	assert.NotZero(t, agendamentoID)

	// 4. Cancelar agendamento
	w = doRequest("POST", fmt.Sprintf("/agendamento/cancelar/%d", agendamentoID), nil, map[string]string{
		"Authorization": "Bearer " + cliToken,
	})
	assert.Equal(t, http.StatusOK, w.Code)

	// 5. Verificar status do agendamento mudou para CANCELADO
	w = doRequest("GET", fmt.Sprintf("/agendamento/%d", agendamentoID), nil, map[string]string{
		"Authorization": "Bearer " + cliToken,
	})
	assert.Equal(t, http.StatusOK, w.Code)

	var getAgResult map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &getAgResult)
	assert.Equal(t, "CANCELADO", getAgResult["status"])
}
