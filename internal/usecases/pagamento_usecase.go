package usecases

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"strings"
	"time"

	gatewaympesa "github.com/ManuelMassora/servicoJa-api/internal/infra/gateway_mpesa"
	"github.com/ManuelMassora/servicoJa-api/internal/model"
)

const (
	// M-Pesa API URLs
	mpesaC2BURL   = "https://api.sandbox.vm.co.mz:18352/ipg/v1x/c2bPayment/singleStage/"
	mpesaB2CURL   = "https://api.sandbox.vm.co.mz:18345/ipg/v1x/b2cPayment/"
	mpesaQueryURL = "https://api.sandbox.vm.co.mz:18353/ipg/v1x/queryTransactionStatus/"

	// Response codes
	mpesaSuccessCode = "INS-0"

	// Delays
	statusCheckDelay = 20 * time.Second
)

type PagamentoUseCase interface {
	IniciarPagamentoC2B(ctx context.Context, idPagamento uint, telefone string) error
	ConfirmarPagamentoC2B(ctx context.Context, referencia string) error
	ProcessarCancelamentoComReembolso(ctx context.Context, idServico uint, idUsuarioCancelou uint) error
	ProcessarReembolsoVaga(ctx context.Context, idVaga uint, idUsuarioCancelou uint) error
	ProcessarPagamentoPrestador(ctx context.Context, idServico uint) error
	ProcessarCallbackMpesa(ctx context.Context, payload gatewaympesa.MpesaCallbackPayload) error
	ProcessarQuerySimulada(ctx context.Context, referencia string) error
}

type pagamentoUseCaseImpl struct {
	repo                     model.PagamentoRepo
	servicoRepo              model.ServicoRepo
	vagaRepo                 model.VagaRepo
	agendamentoRepo          model.AgendamentoRepo
	usuarioRepo              model.UsuarioRepo
	notificacaoRepo          model.NotificacaoRepo
	mpesaGateway             *gatewaympesa.MpesaGateway
	mpesaAppKey              string
	mpesaAppPub              string
	mpesaServiceProviderCode string
}

func NewPagamentoUseCase(
	repo model.PagamentoRepo,
	servicoRepo model.ServicoRepo,
	vagaRepo model.VagaRepo,
	agendamentoRepo model.AgendamentoRepo,
	usuarioRepo model.UsuarioRepo,
	notificacaoRepo model.NotificacaoRepo,
	mpesaGateway *gatewaympesa.MpesaGateway,
	mpesaAppKey, mpesaAppPub, mpesaServiceProviderCode string,
) PagamentoUseCase {
	return &pagamentoUseCaseImpl{
		repo:                     repo,
		servicoRepo:              servicoRepo,
		vagaRepo:                 vagaRepo,
		agendamentoRepo:          agendamentoRepo,
		usuarioRepo:              usuarioRepo,
		notificacaoRepo:          notificacaoRepo,
		mpesaGateway:             mpesaGateway,
		mpesaAppKey:              mpesaAppKey,
		mpesaAppPub:              mpesaAppPub,
		mpesaServiceProviderCode: mpesaServiceProviderCode,
	}
}

func (uc *pagamentoUseCaseImpl) IniciarPagamentoC2B(ctx context.Context, idPagamento uint, telefone string) error {
	p, err := uc.repo.BuscarPorID(ctx, idPagamento)
	if err != nil {
		log.Printf("erro ao buscar pagamento: %v", err)
		return fmt.Errorf("pagamento não encontrado")
	}

	// Limpar telefone (remover + se houver)
	telefone = strings.TrimPrefix(telefone, "+")

	amount := fmt.Sprintf("%.2f", p.Valor)
	log.Println(p.Referencia)

	var targetID uint
	if p.IDVaga != nil {
		targetID = *p.IDVaga
	} else if p.IDAgendamento != nil {
		targetID = *p.IDAgendamento
	}

	payload := uc.construirPayloadC2B(amount, telefone, p.Referencia, targetID)
	resp, err := uc.enviarPagamento(mpesaC2BURL, payload)
	if err != nil {
		log.Printf("erro ao enviar pagamento: %v", err)
		return err
	}

	log.Printf("Pagamento C2B iniciado: %s", string(resp))

	// Extrair ConversationID da resposta
	var c2bResp gatewaympesa.MpesaCallbackPayload
	if err := json.Unmarshal(resp, &c2bResp); err != nil {
		log.Printf("erro ao parsear resposta C2B: %v", err)
		return nil // Pagamento foi enviado, mas não conseguimos agendar verificação
	}

	// Se por algum motivo o status já vier como Completed (ex: algumas APIs específicas)
	if c2bResp.ResponseTransactionStatus == "Completed" || c2bResp.ResponseCode == mpesaSuccessCode {
		log.Printf("Pagamento %s confirmed na resposta síncrona", p.Referencia)
		go uc.ConfirmarPagamentoC2B(context.Background(), p.Referencia)
		return nil
	}

	// Agendar verificação em background se não estiver concluído ainda
	go uc.agendarVerificacaoStatus(context.Background(), p.Referencia, c2bResp.ConversationID)

	return nil
}

func (uc *pagamentoUseCaseImpl) ConfirmarPagamentoC2B(ctx context.Context, referencia string) error {
	p, err := uc.repo.BuscarPorReferencia(ctx, referencia)
	if err != nil {
		return fmt.Errorf("pagamento com referência %s não encontrado", referencia)
	}

	if p.Status != model.StatusPendente {
		return nil // Já processado
	}

	err = uc.repo.AtualizarStatusPorReferencia(ctx, p.Referencia, model.StatusConcluido)
	if err != nil {
		return err
	}

	if p.IDVaga != nil {
		vaga, err := uc.vagaRepo.BuscarPorID(ctx, *p.IDVaga)
		if err != nil {
			return err
		}
		vaga.Status = model.StatusDisponivel
		return uc.vagaRepo.Salvar(ctx, vaga)
	}

	if p.IDAgendamento != nil {
		ag, err := uc.agendamentoRepo.BuscarPorID(ctx, *p.IDAgendamento)
		if err != nil {
			return err
		}
		ag.Status = "PENDENTE" // Agora o prestador pode ver

		// Enviar notificação ao prestador agora que está pago
		_ = uc.notificacaoRepo.Enviar(ctx, &model.Notificacao{
			IDUsuario: ag.Catalogo.Prestador.IDUsuario,
			Titulo:    "Novo Agendamento Confirmado",
			Mensagem:  "Um novo agendamento foi pago e aguarda sua aprovação para o serviço: " + ag.Catalogo.Nome,
		})

		return uc.agendamentoRepo.AtualizarStatus(ctx, ag.ID, "PENDENTE")
	}

	return nil
}

func (uc *pagamentoUseCaseImpl) ProcessarCancelamentoComReembolso(ctx context.Context, idServico uint, idUsuarioCancelou uint) error {
	servico, err := uc.servicoRepo.BuscarPorID(ctx, idServico)
	if err != nil {
		return err
	}

	// Incrementar contador de cancelamentos
	count, err := uc.usuarioRepo.IncrementarCancelamentos(ctx, idUsuarioCancelou)
	if err == nil && count >= 5 {
		// Suspender por 24 horas se atingir 5 cancelamentos
		suspensao := time.Now().Add(24 * time.Hour)
		_ = uc.usuarioRepo.SuspenderUsuario(ctx, idUsuarioCancelou, suspensao)
	}

	// Buscar pagamento original (C2B)
	pagamento, err := uc.repo.BuscarPorServico(ctx, idServico)
	if err != nil {
		// Se não achar por serviço, tenta por agendamento ou vaga
		if servico.IDAgendamento != nil {
			pagamento, _ = uc.repo.BuscarPorAgendamento(ctx, *servico.IDAgendamento)
		} else if servico.IDVaga != nil {
			pagamento, _ = uc.repo.BuscarPorVaga(ctx, *servico.IDVaga)
		}
	}

	if pagamento == nil {
		return fmt.Errorf("pagamento não encontrado para reembolso")
	}

	// Buscar informações do cliente para o MSISDN
	cliente, err := uc.usuarioRepo.BuscarPorID(ctx, servico.IDCliente)
	if err != nil {
		return fmt.Errorf("erro ao buscar cliente para reembolso: %w", err)
	}
	telefone := strings.TrimPrefix(cliente.Telefone, "+")

	amount := fmt.Sprintf("%.2f", pagamento.Valor)
	thirdPartyRef := fmt.Sprintf("REF_REFUND_%d", servico.ID)

	payload := uc.construirPayloadB2C(amount, telefone, thirdPartyRef, servico.ID)
	resp, err := uc.enviarPagamento(mpesaB2CURL, payload)
	if err != nil {
		return fmt.Errorf("erro ao processar reembolso M-Pesa: %w", err)
	}

	log.Printf("Reembolso B2C iniciado: %s", string(resp))

	err = uc.repo.AtualizarStatus(ctx, pagamento.ID, model.StatusCancelado)
	return err
}

func (uc *pagamentoUseCaseImpl) ProcessarReembolsoVaga(ctx context.Context, idVaga uint, idUsuarioCancelou uint) error {
	// Incrementar contador de cancelamentos
	count, err := uc.usuarioRepo.IncrementarCancelamentos(ctx, idUsuarioCancelou)
	if err == nil && count >= 5 {
		suspensao := time.Now().Add(24 * time.Hour)
		_ = uc.usuarioRepo.SuspenderUsuario(ctx, idUsuarioCancelou, suspensao)
	}

	// Buscar pagamento original (C2B) vinculado à vaga
	pagamento, err := uc.repo.BuscarPorVaga(ctx, idVaga)
	if err != nil || pagamento == nil {
		return fmt.Errorf("pagamento não encontrado para reembolso de vaga")
	}

	if pagamento.Status == model.StatusCancelado {
		return nil // Já reembolsado
	}

	vaga, err := uc.vagaRepo.BuscarPorID(ctx, idVaga)
	if err != nil {
		return fmt.Errorf("erro ao buscar vaga para reembolso: %w", err)
	}

	// Buscar informações do cliente para o MSISDN
	cliente, err := uc.usuarioRepo.BuscarPorID(ctx, vaga.IDCliente)
	if err != nil {
		return fmt.Errorf("erro ao buscar cliente para reembolso: %w", err)
	}
	telefone := strings.TrimPrefix(cliente.Telefone, "+")

	amount := fmt.Sprintf("%.2f", pagamento.Valor)
	thirdPartyRef := fmt.Sprintf("REF_VAGA_REFUND_%d", idVaga)

	// Usamos 0 como servicoID no payload, ou algum identificador de que é Vaga
	payload := uc.construirPayloadB2C(amount, telefone, thirdPartyRef, 0)
	resp, err := uc.enviarPagamento(mpesaB2CURL, payload)
	if err != nil {
		log.Printf("Aviso: Falha no envio M-Pesa B2C (vaga %d): %v", idVaga, err)
		// Em ambiente de teste, podemos querer continuar marcando como cancelado
	} else {
		log.Printf("Reembolso Vaga B2C iniciado: %s", string(resp))
	}

	return uc.repo.AtualizarStatus(ctx, pagamento.ID, model.StatusCancelado)
}

func (uc *pagamentoUseCaseImpl) ProcessarPagamentoPrestador(ctx context.Context, idServico uint) error {
	servico, err := uc.servicoRepo.BuscarPorID(ctx, idServico)
	if err != nil {
		return err
	}

	pagamento, err := uc.repo.BuscarPorServico(ctx, idServico)
	if err != nil {
		if servico.IDAgendamento != nil {
			pagamento, _ = uc.repo.BuscarPorAgendamento(ctx, *servico.IDAgendamento)
		} else if servico.IDVaga != nil {
			pagamento, _ = uc.repo.BuscarPorVaga(ctx, *servico.IDVaga)
		}
	}

	if pagamento == nil {
		return fmt.Errorf("pagamento não encontrado para repasse")
	}

	// Buscar informações do prestador para o MSISDN
	prestador, err := uc.usuarioRepo.BuscarPorID(ctx, servico.IDPrestador)
	if err != nil {
		return fmt.Errorf("erro ao buscar prestador para pagamento: %w", err)
	}
	telefone := strings.TrimPrefix(prestador.Telefone, "+")

	// Aqui poderíamos descontar comissão da plataforma
	amount := fmt.Sprintf("%.2f", pagamento.Valor)
	thirdPartyRef := fmt.Sprintf("REF_PAYOUT_%d", servico.ID)

	payload := uc.construirPayloadB2C(amount, telefone, thirdPartyRef, servico.ID)
	resp, err := uc.enviarPagamento(mpesaB2CURL, payload)
	if err != nil {
		return fmt.Errorf("erro ao processar pagamento ao prestador M-Pesa: %w", err)
	}

	log.Printf("Pagamento B2C ao prestador iniciado: %s", string(resp))

	return nil
}

// construirPayloadB2C builds the B2C payment payload
func (uc *pagamentoUseCaseImpl) construirPayloadB2C(amount, customerMSISDN, thirdPartyRef string, servicoID uint) []byte {
	return []byte(fmt.Sprintf(`{
		"input_Amount": "%s",
		"input_CustomerMSISDN": "%s",
		"input_ThirdPartyReference": "%s",
		"input_ServiceProviderCode": "%s",
		"input_TransactionReference": "%s%d",
		"input_TransactionDesc": "Pagamento B2C do servico %d"
	}`, amount, customerMSISDN, thirdPartyRef, uc.mpesaServiceProviderCode, "B2C", servicoID, servicoID))
}

// construirPayloadC2B builds the C2B payment payload
func (uc *pagamentoUseCaseImpl) construirPayloadC2B(amount, customerMSISDN, thirdPartyRef string, AgendamentoVagaID uint) []byte {
	return []byte(fmt.Sprintf(`{
		"input_Amount": "%s",
		"input_CustomerMSISDN": "%s",
		"input_ThirdPartyReference": "%s",
		"input_ServiceProviderCode": "%s",
		"input_TransactionReference": "%s%d",
		"input_TransactionDesc": "Pagamento de agendamento ou vaga %d"
	}`, amount, customerMSISDN, thirdPartyRef, uc.mpesaServiceProviderCode, "C2B", AgendamentoVagaID, AgendamentoVagaID))
}

// enviarPagamento sends payment request to M-Pesa
func (uc *pagamentoUseCaseImpl) enviarPagamento(url string, payload []byte) ([]byte, error) {
	token, err := uc.generateBearerToken()
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar token: %w", err)
	}

	response, err := uc.mpesaGateway.SendPayment(url, token, payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar pagamento M-Pesa: %w", err)
	}

	return response, nil
}

// agendarVerificacaoStatus schedules a status check after a delay
func (uc *pagamentoUseCaseImpl) agendarVerificacaoStatus(ctx context.Context, referencia, conversationID string) {
	time.Sleep(statusCheckDelay)

	if err := uc.verificarStatusPagamento(ctx, referencia, conversationID); err != nil {
		log.Printf("Erro ao verificar status do pagamento %s: %v", referencia, err)
	}
}

// verificarStatusPagamento queries the payment status from M-Pesa
func (uc *pagamentoUseCaseImpl) verificarStatusPagamento(ctx context.Context, referencia, conversationID string) error {
	token, err := uc.generateBearerToken()
	if err != nil {
		return fmt.Errorf("erro ao gerar token: %w", err)
	}

	params := uc.construirQueryParams(referencia, conversationID)

	response, err := uc.mpesaGateway.SendGet(mpesaQueryURL, token, params)
	if err != nil {
		return fmt.Errorf("erro na query: %w", err)
	}

	var queryResp gatewaympesa.MpesaQueryResponse
	if err := json.Unmarshal(response, &queryResp); err != nil {
		log.Printf("Erro ao parsear resposta da query: %v. Raw: %s", err, string(response))
		// Fallback para busca por string se o JSON for inesperado
		if strings.Contains(string(response), mpesaSuccessCode) || strings.Contains(string(response), "Completed") {
			return uc.ConfirmarPagamentoC2B(ctx, referencia)
		}
		return nil
	}

	// Se a query indicar sucesso via ResponseCode or TransactionStatus
	if queryResp.ResponseCode == mpesaSuccessCode ||
		queryResp.TransactionStatus == "Completed" ||
		queryResp.ResponseTransactionStatus == "Completed" {
		return uc.ConfirmarPagamentoC2B(ctx, referencia)
	}

	log.Printf("Query para %s ainda não concluída. Status: %s, Desc: %s",
		referencia, queryResp.TransactionStatus, queryResp.ResponseDesc)

	return nil
}

func (uc *pagamentoUseCaseImpl) ProcessarCallbackMpesa(ctx context.Context, payload gatewaympesa.MpesaCallbackPayload) error {
	log.Printf("Processando callback M-Pesa: Ref=%s, Code=%s, Status=%s",
		payload.ThirdPartyReference, payload.ResponseCode, payload.ResponseTransactionStatus)

	if payload.ResponseCode == mpesaSuccessCode || payload.ResponseTransactionStatus == "Completed" {
		log.Printf("Confirmando pagamento via callback: %s", payload.ThirdPartyReference)
		return uc.ConfirmarPagamentoC2B(ctx, payload.ThirdPartyReference)
	}

	log.Printf("Callback indicou falha ou pendência: %s - %s", payload.ResponseCode, payload.ResponseDesc)
	return nil
}

func (uc *pagamentoUseCaseImpl) ProcessarQuerySimulada(ctx context.Context, referencia string) error {
	log.Printf("Processando query simulada para referência: %s", referencia)
	return uc.ConfirmarPagamentoC2B(ctx, referencia)
}

// construirQueryParams builds the query parameters for transaction status check
func (uc *pagamentoUseCaseImpl) construirQueryParams(referencia, conversationID string) map[string]string {
	return map[string]string{
		"input_ThirdPartyReference": referencia,
		"input_QueryReference":      conversationID,
		"input_ServiceProviderCode": uc.mpesaServiceProviderCode,
	}
}

func (uc *pagamentoUseCaseImpl) generateBearerToken() (string, error) {
	if uc.mpesaAppPub == "" {
		return "", fmt.Errorf("chave pública M-Pesa não configurada")
	}

	// Parse public key
	rsaPub, err := uc.parsePublicKey(uc.mpesaAppPub)
	if err != nil {
		return "", err
	}

	// Encrypt API key
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, []byte(uc.mpesaAppKey))
	if err != nil {
		return "", fmt.Errorf("erro ao criptografar chave: %w", err)
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// parsePublicKey parses the RSA public key from PEM format
func (uc *pagamentoUseCaseImpl) parsePublicKey(keyData string) (*rsa.PublicKey, error) {
	// Add PEM headers if not present
	if !strings.Contains(keyData, "-----BEGIN") {
		keyData = fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----", keyData)
	}

	// Decode PEM
	block, _ := pem.Decode([]byte(keyData))
	if block == nil {
		return nil, fmt.Errorf("falha ao parsear PEM da chave pública: formato inválido")
	}

	// Parse X509 public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("falha ao parsear chave pública X509: %w", err)
	}

	// Assert RSA public key
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("não é uma chave RSA válida")
	}

	return rsaPub, nil
}
