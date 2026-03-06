package handler

import (
	"context"
	"net/http"

	gatewaympesa "github.com/ManuelMassora/servicoJa-api/internal/infra/gateway_mpesa"
	"github.com/gin-gonic/gin"
)

type PagamentoUseCase interface {
	ProcessarCallbackMpesa(ctx context.Context, payload gatewaympesa.MpesaCallbackPayload) error
	ProcessarQuerySimulada(ctx context.Context, referencia string) error
}

type PagamentoHandler struct {
	uc PagamentoUseCase
}

func NewPagamentoHandler(uc PagamentoUseCase) *PagamentoHandler {
	return &PagamentoHandler{uc: uc}
}

// ReceiveCallback handles the M-Pesa C2B callback
func (h *PagamentoHandler) ReceiveCallback(c *gin.Context) {
	var payload gatewaympesa.MpesaCallbackPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido"})
		return
	}

	err := h.uc.ProcessarCallbackMpesa(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// M-Pesa espera um 200 OK com uma resposta específica as vezes,
	// mas geralmente apenas 200 OK serve para confirmar recebimento.
	c.JSON(http.StatusOK, gin.H{
		"output_ResponseCode": "INS-0",
		"output_ResponseDesc": "Callback received",
	})
}

// SimularQuery triggers a manual status check simulation (useful for sandbox)
func (h *PagamentoHandler) SimularQuery(c *gin.Context) {
	referencia := c.Query("ref")
	if referencia == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "referência (ref) é obrigatória"})
		return
	}

	err := h.uc.ProcessarQuerySimulada(c.Request.Context(), referencia)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pagamento confirmado via simulação de query"})
}
