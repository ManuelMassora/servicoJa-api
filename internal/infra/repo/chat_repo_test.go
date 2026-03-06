package repo_test

import (
	"context"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/infra/repo"
	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ChatTestSuite struct {
	suite.Suite
	repo          model.ChatRepo
	msgRepo       model.MensagemRepo
	usuarioRepo   model.UsuarioRepo
	prestadorRepo model.PrestadorRepo
	clienteRepo   model.ClienteRepo
	servicoRepo   model.ServicoRepo
	ctx           context.Context
}

func (s *ChatTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.repo = repo.NewChatRepository(TestDB)
	s.msgRepo = repo.NewMensagemRepository(TestDB)
	s.usuarioRepo = repo.NewUsuarioRepository(TestDB)
	s.prestadorRepo = repo.NewPrestadorRepository(TestDB)
	s.clienteRepo = repo.NewClienteRepository(TestDB)
	s.servicoRepo = repo.NewServicoRepository(TestDB)
}

func (s *ChatTestSuite) SetupTest() {
	TestDB.Exec("TRUNCATE TABLE mensagems RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE chats RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE servicos RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE clientes RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE prestadors RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE usuarios RESTART IDENTITY CASCADE")
}

func TestChatTestSuite(t *testing.T) {
	suite.Run(t, new(ChatTestSuite))
}

func (s *ChatTestSuite) createDependencies() (uint, uint, uint) {
	// Cliente
	u1 := &model.Usuario{Nome: "Cliente", Telefone: "84111", RolePermissaoID: 1}
	s.usuarioRepo.Criar(s.ctx, u1)
	c := &model.Cliente{IDUsuario: u1.ID, Nome: u1.Nome, Telefone: u1.Telefone}
	s.clienteRepo.Criar(s.ctx, c)

	// Prestador
	u2 := &model.Usuario{Nome: "Prestador", Telefone: "84222", RolePermissaoID: 2}
	s.usuarioRepo.Criar(s.ctx, u2)
	p := &model.Prestador{IDUsuario: u2.ID, Nome: u2.Nome, Telefone: u2.Telefone}
	s.prestadorRepo.Criar(s.ctx, p)

	// Servico
	serv := &model.Servico{
		IDCliente:   c.IDUsuario,
		IDPrestador: p.IDUsuario,
		Preco:       100.0,
		Status:      model.StatusPendente,
		Localizacao: "Maputo",
	}
	s.servicoRepo.Criar(s.ctx, serv)

	return c.IDUsuario, p.IDUsuario, serv.ID
}

func (s *ChatTestSuite) TestChat_CriarChat() {
	cID, pID, sID := s.createDependencies()

	chat := &model.Chat{
		ServicoID:   int64(sID),
		PrestadorID: int64(pID),
		IDCliente:   int64(cID),
	}

	err := s.repo.CriarChat(s.ctx, chat)
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), chat.ID)
}

func (s *ChatTestSuite) TestChat_ListarChatsPorUsuario() {
	cID, pID, sID := s.createDependencies()
	chat := &model.Chat{ServicoID: int64(sID), PrestadorID: int64(pID), IDCliente: int64(cID)}
	s.repo.CriarChat(s.ctx, chat)

	list, err := s.repo.ListarChatsPorUsuario(s.ctx, uint(cID), nil, "", "", 10, 0)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), list, 1)
}

func (s *ChatTestSuite) TestMensagem_EnviarEListar() {
	cID, pID, sID := s.createDependencies()
	chat := &model.Chat{ServicoID: int64(sID), PrestadorID: int64(pID), IDCliente: int64(cID)}
	s.repo.CriarChat(s.ctx, chat)

	msg := &model.Mensagem{
		IDChat:        chat.ID,
		IDRemetente:   uint(cID),
		RemetenteTipo: "CLIENTE",
		Conteudo:      "Olá, tudo bem?",
	}

	err := s.msgRepo.EnviarMensagem(s.ctx, msg)
	assert.NoError(s.T(), err)

	msgs, err := s.msgRepo.ListarMensagens(s.ctx, chat.ID, nil, "created_at", "asc", 10, 0)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), msgs, 1)
	assert.Equal(s.T(), "Olá, tudo bem?", msgs[0].Conteudo)
}
