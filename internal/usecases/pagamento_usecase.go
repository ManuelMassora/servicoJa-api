package usecases

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
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
	mpesaB2CURL   = "https://api.sandbox.vm.co.mz:18352/ipg/v1x/b2cPayment/"
	mpesaQueryURL = "https://api.sandbox.vm.co.mz:18352/ipg/v1x/queryTransactionStatus/"

	// Response codes
	mpesaSuccessCode = "INS-0"

	// Delays
	statusCheckDelay = 20 * time.Second
)

type PagamentoUseCase struct {
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
) *PagamentoUseCase {
	return &PagamentoUseCase{
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

func (uc *PagamentoUseCase) IniciarPagamentoC2B(ctx context.Context, idPagamento uint, telefone string) error {
	pagamentos, err := uc.repo.ListarPorUsuario(ctx, 0, map[string]interface{}{"id": idPagamento}, "", "", 1, 0)
	if err != nil || len(pagamentos) == 0 {
		return fmt.Errorf("pagamento não encontrado")
	}
	p := &pagamentos[0]

	// Limpar telefone (remover + se houver)
	telefone = strings.TrimPrefix(telefone, "+")

	amount := fmt.Sprintf("%.2f", p.Valor)
	thirdPartyRef := fmt.Sprintf("REF%d", p.ID)

	var targetID uint
	if p.IDVaga != nil {
		targetID = *p.IDVaga
	} else if p.IDAgendamento != nil {
		targetID = *p.IDAgendamento
	}

	payload := uc.construirPayloadC2B(amount, telefone, thirdPartyRef, targetID)
	resp, err := uc.enviarPagamento(mpesaC2BURL, payload)
	if err != nil {
		return err
	}

	log.Printf("Pagamento C2B iniciado: %s", string(resp))

	// Agendar verificação em background
	go uc.agendarVerificacaoStatus(context.Background(), thirdPartyRef)

	return nil
}

func (uc *PagamentoUseCase) ConfirmarPagamentoC2B(ctx context.Context, idPagamento uint) error {
	pagamento, err := uc.repo.ListarPorUsuario(ctx, 0, map[string]interface{}{"id": idPagamento}, "", "", 1, 0)
	if err != nil || len(pagamento) == 0 {
		return fmt.Errorf("pagamento não encontrado")
	}
	p := &pagamento[0]

	if p.Status != model.StatusPendente {
		return nil // Já processado
	}

	err = uc.repo.AtualizarStatus(ctx, p.ID, model.StatusConcluido)
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

func (uc *PagamentoUseCase) ProcessarCancelamentoComReembolso(ctx context.Context, idServico uint, idUsuarioCancelou uint) error {
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

func (uc *PagamentoUseCase) ProcessarPagamentoPrestador(ctx context.Context, idServico uint) error {
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
func (uc *PagamentoUseCase) construirPayloadB2C(amount, customerMSISDN, thirdPartyRef string, servicoID uint) []byte {
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
func (uc *PagamentoUseCase) construirPayloadC2B(amount, customerMSISDN, thirdPartyRef string, AgendamentoVagaID uint) []byte {
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
func (uc *PagamentoUseCase) enviarPagamento(url string, payload []byte) ([]byte, error) {
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
func (uc *PagamentoUseCase) agendarVerificacaoStatus(ctx context.Context, referencia string) {
	time.Sleep(statusCheckDelay)

	if err := uc.verificarStatusPagamento(ctx, referencia); err != nil {
		log.Printf("Erro ao verificar status do pagamento %s: %v", referencia, err)
	}
}

// verificarStatusPagamento queries the payment status from M-Pesa
func (uc *PagamentoUseCase) verificarStatusPagamento(ctx context.Context, referencia string) error {
	token, err := uc.generateBearerToken()
	if err != nil {
		return fmt.Errorf("erro ao gerar token: %w", err)
	}

	payload := uc.construirPayloadQuery(referencia)

	response, err := uc.mpesaGateway.SendPost(mpesaQueryURL, token, payload)
	if err != nil {
		return fmt.Errorf("erro na query: %w", err)
	}

	respStr := string(response)
	log.Printf("Resultado da query para %s: %s\n", referencia, respStr)

	// Se a query indicar sucesso, confirmar o pagamento no sistema
	// Na prática, deve-se parsear o JSON e verificar "output_ResponseCode" == mpesaSuccessCode
	if strings.Contains(respStr, mpesaSuccessCode) {
		// Extrair ID do pagamento da referência (ex: REF123 -> 123)
		var idPagamento uint
		_, err := fmt.Sscanf(referencia, "REF%d", &idPagamento)
		if err == nil {
			return uc.ConfirmarPagamentoC2B(ctx, idPagamento)
		}
	}

	return nil
}

// construirPayloadQuery builds the query status payload
func (uc *PagamentoUseCase) construirPayloadQuery(referencia string) []byte {
	return []byte(fmt.Sprintf(`{
		"input_ThirdPartyReference": "%s",
		"input_ServiceProviderCode": "%s"
	}`, referencia, uc.mpesaServiceProviderCode))
}

func (uc *PagamentoUseCase) generateBearerToken() (string, error) {
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
func (uc *PagamentoUseCase) parsePublicKey(keyData string) (*rsa.PublicKey, error) {
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
