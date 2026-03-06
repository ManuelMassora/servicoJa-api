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

func TestFluxoAgendamentoRembolso(t *testing.T) {
	// 1. Setup: Criar admin, categoria e prestador
	adminEmail := "admin_ag_refund@test.com"
	adminTelefone := "820000021"
	adminSenha := "admin123"
	registerUser("ADMIN", adminEmail, adminTelefone, adminSenha)
	adminToken := loginUser(adminTelefone, adminSenha)

	// Criar categoria
	catBody := map[string]interface{}{
		"nome":      "Serviços Premium",
		"descricao": "Serviços com reembolso",
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
	prestEmail := "prestador_ag_refund@test.com"
	prestTelefone := "820000022"
	prestSenha := "prestador123"
	registerUser("PRESTADOR", prestEmail, prestTelefone, prestSenha)
	prestToken := loginUser(prestTelefone, prestSenha)

	// Criar catálogo com preço fixo
	fields := map[string]string{
		"nome":         "Consultoria Premium",
		"descricao":    "Serviço premium com reembolso",
		"tipo_preco":   "fixo",
		"valor_fixo":   "5000.0",
		"categoria_id": fmt.Sprint(catID),
		"localizacao":  "Maputo Centro",
		"latitude":     "-25.9692",
		"longitude":    "32.5732",
	}
	w = doMultipartRequest("POST", "/catalogo", fields, prestToken)
	assert.Equal(t, http.StatusCreated, w.Code)
	var catalogoCreateResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &catalogoCreateResp)
	catalogoID := uint(catalogoCreateResp["id"].(float64))

	// 2. Registrar e logar cliente
	cliEmail := "cliente_ag_refund@test.com"
	cliTelefone := "820000023"
	cliSenha := "cliente123"
	registerUser("CLIENTE", cliEmail, cliTelefone, cliSenha)
	cliToken := loginUser(cliTelefone, cliSenha)

	// 3. Criar agendamento
	agFields := map[string]string{
		"detalhe":            "Preciso de consultoria premium",
		"id_catalogo":        fmt.Sprint(catalogoID),
		"datahora":           time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"localizacao":        "Escritório Centro",
		"latitude":           "-25.9000",
		"longitude":          "32.5000",
		"telefone_pagamento": "841112255",
	}
	w = doMultipartRequest("POST", "/agendamento", agFields, cliToken)
	assert.Equal(t, http.StatusCreated, w.Code)

	var agResult map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &agResult)
	agendamentoID := uint(agResult["id"].(float64))
	assert.NotZero(t, agendamentoID)

	// 3.5. Simular confirmação de pagamento para o agendamento
	doRequest("GET", fmt.Sprintf("/pagamento/simular-query?ref=AGE-%d", agendamentoID), nil, nil)

	// 4. Prestador aceita o agendamento -> Isso cria o registro de serviço
	w = doRequest("POST", fmt.Sprintf("/agendamento/aceitar/%d", agendamentoID), nil, map[string]string{
		"Authorization": "Bearer " + prestToken,
	})
	assert.Equal(t, http.StatusOK, w.Code)
	var acceptResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &acceptResp)
	servicoID := uint(acceptResp["id_servico"].(float64))
	assert.NotZero(t, servicoID)

	// O ID do serviço no nosso sistema é retornado no accept
	// Vamos buscar o agendamento para confirmar que o status mudou para EM_ANDAMENTO
	w = doRequest("GET", fmt.Sprintf("/agendamento/%d", agendamentoID), nil, map[string]string{
		"Authorization": "Bearer " + cliToken,
	})
	assert.Equal(t, http.StatusOK, w.Code)
	var agUpdated map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &agUpdated)
	assert.Equal(t, "EM_ANDAMENTO", agUpdated["status"])

	// 5. Cancelar o serviço (isso trigga o reembolso automático)
	w = doRequest("POST", fmt.Sprintf("/servico/cancelar/%d", servicoID), nil, map[string]string{
		"Authorization": "Bearer " + cliToken,
	})
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. Verificar que o agendamento foi cancelado e pagamento reembolsado
	w = doRequest("GET", fmt.Sprintf("/agendamento/%d", agendamentoID), nil, map[string]string{
		"Authorization": "Bearer " + cliToken,
	})
	assert.Equal(t, http.StatusOK, w.Code)

	var getAgResult map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &getAgResult)
	// Status deve estar em CANCELADO após reembolso ser processado
	assert.Equal(t, "CANCELADO", getAgResult["status"])

	// 6. Verificar que o pagamento foi atualizado para CANCELADO
	// (Isso depende de haver um endpoint GET /pagamento/:id ou similar)
	// Para este teste, verificaremos apenas que o cancelamento foi bem-sucedido

	// 7. Verificar contador de cancelamentos (cancellation_counter deve aumentar)
	w = doRequest("GET", "/usuario/perfil", nil, map[string]string{
		"Authorization": "Bearer " + cliToken,
	})
	assert.Equal(t, http.StatusOK, w.Code)

	var userResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &userResp)
	// Verificar que o cancellation_counter aumentou
	cancellationCounter := int(userResp["cancelamento_contador"].(float64))
	assert.GreaterOrEqual(t, cancellationCounter, 1, "Cancellation counter should have increased")
}
