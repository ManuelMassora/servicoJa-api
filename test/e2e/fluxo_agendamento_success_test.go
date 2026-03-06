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

func TestFluxoAgendamentoSuccess(t *testing.T) {
	// 1. Criar ADMIN para criar categoria
	adminEmail := "admin@test.com"
	adminTelefone := "820000001"
	adminSenha := "admin123"
	registerUser("ADMIN", adminEmail, adminTelefone, adminSenha)
	adminToken := loginUser(adminTelefone, adminSenha)

	// 2. Criar Categoria
	catBody := map[string]interface{}{
		"nome":      "Limpeza Doméstica",
		"descricao": "Serviços de limpeza para casas",
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

	// 3. Registrar e Logar Prestador
	prestEmail := "prestador@test.com"
	prestTelefone := "820000002"
	prestSenha := "prestador123"
	registerUser("PRESTADOR", prestEmail, prestTelefone, prestSenha)
	prestToken := loginUser(prestTelefone, prestSenha)

	// 4. Criar Catálogo (Serviço do Prestador)
	fields := map[string]string{
		"nome":         "Faxina Completa",
		"descricao":    "Faxina profunda em residências",
		"tipo_preco":   "fixo",
		"valor_fixo":   "2500.0",
		"categoria_id": fmt.Sprint(catID),
		"localizacao":  "Cidade de Maputo",
		"latitude":     "-25.9692",
		"longitude":    "32.5732",
	}

	w = doMultipartRequest("POST", "/catalogo", fields, prestToken)
	assert.Equal(t, http.StatusCreated, w.Code)
	var catalogoCreateResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &catalogoCreateResp)
	catalogoID := uint(catalogoCreateResp["id"].(float64))

	// 5. Registrar e Logar Cliente
	cliEmail := "cliente_agend@test.com"
	cliTelefone := "820000003"
	cliSenha := "cliente123"
	registerUser("CLIENTE", cliEmail, cliTelefone, cliSenha)
	cliToken := loginUser(cliTelefone, cliSenha)

	// 6. Criar Agendamento
	agFields := map[string]string{
		"detalhe":            "Preciso de uma faxina para amanhã de manhã",
		"id_catalogo":        fmt.Sprint(catalogoID),
		"datahora":           time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"localizacao":        "Minha Casa",
		"latitude":           "-25.9000",
		"longitude":          "32.5000",
		"telefone_pagamento": "841112233",
	}

	w = doMultipartRequest("POST", "/agendamento", agFields, cliToken)
	assert.Equal(t, http.StatusCreated, w.Code)

	var agResult map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &agResult)
	assert.Equal(t, "Agendamento criado com sucesso!", agResult["message"])
	assert.NotNil(t, agResult["id"])
}
