package repo_test

import (
	"context"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/infra/repo"
	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GaleriaTestSuite struct {
	suite.Suite
	repo          model.GaleriaRepo
	usuarioRepo   model.UsuarioRepo
	prestadorRepo model.PrestadorRepo
	ctx           context.Context
}

func (s *GaleriaTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.repo = repo.NewGaleriaRepo(TestDB)
	s.usuarioRepo = repo.NewUsuarioRepository(TestDB)
	s.prestadorRepo = repo.NewPrestadorRepository(TestDB)
}

func (s *GaleriaTestSuite) SetupTest() {
	TestDB.Exec("TRUNCATE TABLE imagem RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE galeria RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE prestadors RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE usuarios RESTART IDENTITY CASCADE")
}

func TestGaleriaTestSuite(t *testing.T) {
	suite.Run(t, new(GaleriaTestSuite))
}

func (s *GaleriaTestSuite) createPrestador() uint {
	u := &model.Usuario{Nome: "Prestador", Telefone: "84111", RolePermissaoID: 2}
	s.usuarioRepo.Criar(s.ctx, u)
	p := &model.Prestador{IDUsuario: u.ID, Nome: u.Nome, Telefone: u.Telefone, Localizacao: "Maputo"}
	s.prestadorRepo.Criar(s.ctx, p)
	return p.IDUsuario
}

func (s *GaleriaTestSuite) TestGaleria_CreateAndAddImage() {
	pID := s.createPrestador()
	galeria := &model.Galeria{
		PrestadorID: pID,
	}

	created, err := s.repo.Create(s.ctx, galeria)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), created)
	assert.NotZero(s.T(), created.ID)

	img := &model.Imagem{
		GaleriaID: created.ID,
		URL:       "http://image.com/art.png",
	}

	err = s.repo.AddImage(s.ctx, img)
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), img.ID)
}

func (s *GaleriaTestSuite) TestGaleria_FindByPrestadorID() {
	pID := s.createPrestador()
	g := &model.Galeria{PrestadorID: pID}
	created, _ := s.repo.Create(s.ctx, g)

	found, err := s.repo.FindByPrestadorID(s.ctx, pID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), found)
	assert.Equal(s.T(), created.ID, found.ID)
}

func (s *GaleriaTestSuite) TestGaleria_CountImages() {
	pID := s.createPrestador()
	g := &model.Galeria{PrestadorID: pID}
	created, _ := s.repo.Create(s.ctx, g)

	img1 := &model.Imagem{GaleriaID: created.ID, URL: "url1"}
	img2 := &model.Imagem{GaleriaID: created.ID, URL: "url2"}
	s.repo.AddImage(s.ctx, img1)
	s.repo.AddImage(s.ctx, img2)

	count, err := s.repo.CountImages(s.ctx, created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(2), count)
}
