package usecases

import (
	"context"
	"errors"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/gorm"
)

type UsuarioUseCase struct {
	usuarioRepo model.UsuarioRepo
	clienteRepo model.ClienteRepo
	prestadorRepo model.PrestadorRepo
	galeriaRepo model.GaleriaRepo
}

func NewUsuarioUseCase(
	usuarioRepo model.UsuarioRepo,
	clienteRepo model.ClienteRepo,
	prestadorRepo model.PrestadorRepo,
	galeriaRepo model.GaleriaRepo,
) *UsuarioUseCase {
	return &UsuarioUseCase{
		usuarioRepo: usuarioRepo,
		clienteRepo: clienteRepo,
		prestadorRepo: prestadorRepo,
		galeriaRepo: galeriaRepo,
	}
}

type UsuarioRequest struct {
	Nome      string `json:"nome" binding:"required"`
	Telefone  string `json:"telefone" binding:"required"`
	Senha     string `json:"senha,omitempty" binding:"required"`
	ImagemURL string `json:"imagem_url"`
}

type UsuarioResponse struct {
	ID        uint   `json:"id"`
	Nome      string `json:"nome"`
	Telefone  string `json:"telefone"`
	ImagemURL string `json:"imagem_url"`
}

type PrestadorRequest struct {
	Nome        string  `form:"nome" binding:"required"`
	Telefone    string  `form:"telefone" binding:"required"`
	Senha       string  `form:"senha" binding:"required"`
	ImagemURL   string	`form:"-"`
	Localizacao string  `form:"localizacao" binding:"required"`
	Latitude    float64 `form:"latitude" binding:"required"`
	Longitude   float64 `form:"longitude" binding:"required"`
}

type PrestadorResponse struct {
	ID          uint    `json:"id"`
	Nome        string  `json:"nome"`
	Telefone    string  `json:"telefone"`
	Localizacao string  `json:"localizacao"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Disponivel  bool    `json:"disponivel"`
	ImagemURL   string  `json:"imagem_url"`
	Galeria	 	[]string `json:"galeria"`
}

func (uc *UsuarioUseCase) CriarAdmin(ctx context.Context, request UsuarioRequest) error{
	
	if err := uc.SeTelefoneExiste(ctx, request.Telefone); err != nil {
		return err
	}
	admin, err := model.NewAdmin(request.Nome, request.Telefone, request.Senha)
	if err != nil {
		return err
	}
	return uc.usuarioRepo.Criar(ctx, admin)
}

func (uc *UsuarioUseCase) CriarCliente(ctx context.Context, request UsuarioRequest) error{
	if err := uc.SeTelefoneExiste(ctx, request.Telefone); err != nil {
		return err
	}
	cliente, err := model.NewCliente(request.Nome, request.Telefone, request.Senha, request.ImagemURL)
	if err != nil {
		return err
	} 
	return uc.clienteRepo.Criar(ctx, cliente)
}

func (uc *UsuarioUseCase) CriarPrestador(ctx context.Context, request PrestadorRequest) error{
	if err := uc.SeTelefoneExiste(ctx, request.Telefone); err != nil {
		return err
	}
	prestador, err := model.NewPrestador(
		request.Localizacao,
		request.Latitude,
		request.Longitude,
		request.Nome,
		request.Telefone,
		request.Senha,
		request.ImagemURL)
	if err != nil {
		return err
	}
	return uc.prestadorRepo.Criar(ctx, prestador)
}

func(uc *UsuarioUseCase) ListarTodosUsuarios(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]UsuarioResponse, error) {
	usuarios, err := uc.usuarioRepo.ListarTodos(ctx, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}

	response := make([]UsuarioResponse, 0, len(usuarios))
	for _, u := range usuarios {
		// Since ImagemURL is in Cliente/Prestador, we need to fetch them.
		// This is a simplification. In a real app, you might want to join tables for efficiency.
		var imagemURL string
		switch u.RolePermissao.Role {
			case "CLIENTE":
				cliente, err := uc.clienteRepo.BuscarPorID(ctx, u.ID)
				if err == nil {
					imagemURL = cliente.ImagemURL
				}
			case "PRESTADOR":
				prestador, err := uc.prestadorRepo.BuscarPorID(ctx, u.ID)
				if err == nil {
					imagemURL = prestador.ImagemURL
				}
		}

		response = append(response, UsuarioResponse{
			ID:        u.ID,
			Nome:      u.Nome,
			Telefone:  u.Telefone,
			ImagemURL: imagemURL,
		})
	}
	return response, nil
}

func(uc *UsuarioUseCase) ListarPrestadores(ctx context.Context, filters map[string]interface{}, statusDisponivel interface{}, orderBy string, orderDir string, limit, offset int) ([]PrestadorResponse, error) {
	prestadores, err := uc.prestadorRepo.Listar(ctx, filters, statusDisponivel, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	prestadoresIDs := make([]uint, 0, len(prestadores))
	for _, p := range prestadores {
		prestadoresIDs = append(prestadoresIDs, p.IDUsuario)
	}
	galerias,  err := uc.galeriaRepo.FindByPrestadorIDs(ctx, prestadoresIDs)
	if err != nil {
		return nil, err
	}
	galeriaPrestadorMap := make(map[uint][]string)
	for _, g := range galerias {
		imagensURLs := make([]string, 0, len(g.Imagens))
		for _, i := range g.Imagens {
			imagensURLs = append(imagensURLs, i.URL)
		}
		galeriaPrestadorMap[g.PrestadorID] = imagensURLs
	}
	response := make([]PrestadorResponse, 0, len(prestadores))
	for _, p := range prestadores {
		imagensURLs := galeriaPrestadorMap[p.IDUsuario]
		response = append(response, PrestadorResponse{
			ID:          p.IDUsuario,
			Nome:        p.Nome,
			Telefone:    p.Telefone,
			Localizacao: p.Localizacao,
			Latitude:    p.Latitude,
			Longitude:   p.Longitude,
			Disponivel:  p.StatusDisponivel,
			ImagemURL:   p.ImagemURL,
			Galeria:     imagensURLs,
		})
	}
	return response, nil
}

func (uc *UsuarioUseCase) ListarPrestadoresPorLocalizacao(ctx context.Context, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]PrestadorResponse, error) {
	prestadores, err := uc.prestadorRepo.FindByLocation(ctx, latitude, longitude, radius, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	prestadoresIDs := make([]uint, 0, len(prestadores))
	for _, p := range prestadores {
		prestadoresIDs = append(prestadoresIDs, p.IDUsuario)
	}
	galerias,  err := uc.galeriaRepo.FindByPrestadorIDs(ctx, prestadoresIDs)
	if err != nil {
		return nil, err
	}
	galeriaPrestadorMap := make(map[uint][]string)
	for _, g := range galerias {
		imagensURLs := make([]string, 0, len(g.Imagens))
		for _, i := range g.Imagens {
			imagensURLs = append(imagensURLs, i.URL)
		}
		galeriaPrestadorMap[g.PrestadorID] = imagensURLs
	}
	response := make([]PrestadorResponse, 0, len(prestadores))
	for _, p := range prestadores {
		imagensURLs := galeriaPrestadorMap[p.IDUsuario]
		response = append(response, PrestadorResponse{
			ID:          p.IDUsuario,
			Nome:        p.Nome,
			Telefone:    p.Telefone,
			Localizacao: p.Localizacao,
			Latitude:    p.Latitude,
			Longitude:   p.Longitude,
			Disponivel:  p.StatusDisponivel,
			ImagemURL:   p.ImagemURL,
			Galeria:     imagensURLs,
		})
	}
	return response, nil
}

func(uc *UsuarioUseCase) SeTelefoneExiste(ctx context.Context, telefone string) error {
		if _, err := uc.usuarioRepo.BuscarPorTelefone(ctx, telefone); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return errors.New("ja existe usuario com mesmo contacto")
}