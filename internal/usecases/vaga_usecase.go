package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
)

type VagaUseCase struct {
	vagaRepo        model.VagaRepo
	anexoImagemRepo model.AnexoImagemRepo
}

func NewVagaUseCase(vagaRepo model.VagaRepo, anexoImagemRepo model.AnexoImagemRepo) *VagaUseCase {
	return &VagaUseCase{vagaRepo: vagaRepo, anexoImagemRepo: anexoImagemRepo}
}

type VagaRequest struct {
	Titulo      string   `json:"titulo" form:"titulo" binding:"required"`
	Descricao   string   `json:"descricao" form:"descricao" binding:"required"`
	Localizacao string   `json:"localizacao" form:"localizacao" binding:"required"`
	Latitude    float64  `json:"latitude" form:"latitude" binding:"required"`
	Longitude   float64  `json:"longitude" form:"longitude" binding:"required"`
	Preco       float64  `json:"preco" form:"preco" binding:"required,gte=0"`
	Urgente     bool     `json:"urgente" form:"urgente"`
	Anexos      []string `binding:"-"`
}

type VagaResponse struct {
	ID          uint    `json:"id"`
	Titulo      string  `json:"titulo"`
	Descricao   string  `json:"descricao"`
	Localizacao string  `json:"localizacao"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Status      string 	`json:"status"`
	Preco       float64 `json:"preco"`
	Urgente     bool    `json:"urgente"`
	Cliente    	string  `json:"cliente"`
	DataCriacao string  `json:"data_criacao"`
	Anexos      []string `json:"anexos"`
}

func(uc *VagaUseCase) CriarVaga(ctx context.Context, req VagaRequest, idUsuario uint) error {
	vaga := &model.Vaga{
		Titulo:      req.Titulo,
		Descricao:   req.Descricao,
		Localizacao: req.Localizacao,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Preco:       req.Preco,
		Status:      model.StatusDisponivel,
		IDCliente:   idUsuario,
		Urgente:     req.Urgente,
	}
	if err := uc.vagaRepo.Criar(ctx, vaga); err != nil {
		return err
	}

	for _, anexoURL := range req.Anexos {
		anexo := &model.AnexoImagem{
			URL:    anexoURL,
			VagaID: &vaga.ID,
		}
		if err := uc.anexoImagemRepo.Create(ctx, anexo); err != nil {
			// In a real application, you might want to handle the rollback of the vaga creation
			return err
		}
	}

	return nil
}

func(uc *VagaUseCase) CancelarVaga(ctx context.Context, id, idUsuario uint) error {
	vaga, err := uc.vagaRepo.BuscarPorID(ctx, id)
	if err != nil {
		return err
	}
	if vaga.IDCliente != idUsuario {
		return errors.New("vaga não pertence ao cliente")
	}
	vaga.DeletedAt.Time = time.Now()
	vaga.DeletedAt.Valid = true
	vaga.Status = model.StatusCancelado
	return uc.vagaRepo.Salvar(ctx, vaga)
}

func(uc *VagaUseCase) ListarVagasDisponiveis(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]VagaResponse, error) {
	vagas, err := uc.vagaRepo.ListarDisponiveis(ctx, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	var vagasIDs []uint
	for _, vaga := range vagas {
		vagasIDs = append(vagasIDs, vaga.ID)
	}

	anexos, err := uc.anexoImagemRepo.FindByVagaIDs(ctx, vagasIDs)
	if err != nil {
		return nil, err
	}
	anexosPorVagaMap := make(map[uint][]string)
	for _, anexo := range anexos {
		anexosPorVagaMap[*anexo.VagaID] = append(anexosPorVagaMap[*anexo.VagaID], anexo.URL)
	}
	var resp []VagaResponse
	for _, vaga := range vagas {
		clienteNome := ""
		if vaga.Cliente != nil {
			clienteNome = vaga.Cliente.Usuario.Nome
		}
		urls := anexosPorVagaMap[vaga.ID]
		resp = append(resp, VagaResponse{
			ID:          vaga.ID,
			Titulo:      vaga.Titulo,
			Descricao:   vaga.Descricao,
			Localizacao: vaga.Localizacao,
			Latitude:    vaga.Latitude,
			Longitude:   vaga.Longitude,
			Status:      string(vaga.Status),
			Preco:       vaga.Preco,
			Urgente:     vaga.Urgente,
			Cliente:     clienteNome,
			DataCriacao: vaga.CreatedAt.Format("2006-01-02 15:04:05"),
			Anexos: urls,
		})
	}
	return resp, nil
}

func(uc *VagaUseCase) ListarPorCliente(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]VagaResponse, error) {
	vagas, err := uc.vagaRepo.ListarPorCliente(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	var vagasIDs []uint
	for _, vaga := range vagas {
		vagasIDs = append(vagasIDs, vaga.ID)
	}

	anexos, err := uc.anexoImagemRepo.FindByVagaIDs(ctx, vagasIDs)
	if err != nil {
		return nil, err
	}
	anexosPorVagaMap := make(map[uint][]string)
	for _, anexo := range anexos {
		anexosPorVagaMap[*anexo.VagaID] = append(anexosPorVagaMap[*anexo.VagaID], anexo.URL)
	}
	var resp []VagaResponse
	for _, vaga := range vagas {
		clienteNome := ""
		if vaga.Cliente != nil {
			clienteNome = vaga.Cliente.Usuario.Nome
		}
		urls := anexosPorVagaMap[vaga.ID]
		resp = append(resp, VagaResponse{
			ID:          vaga.ID,
			Titulo:      vaga.Titulo,
			Descricao:   vaga.Descricao,
			Localizacao: vaga.Localizacao,
			Latitude:    vaga.Latitude,
			Longitude:   vaga.Longitude,
			Status:      string(vaga.Status),
			Preco:       vaga.Preco,
			Urgente:     vaga.Urgente,
			Cliente:     clienteNome,
			DataCriacao: vaga.CreatedAt.Format("2006-01-02 15:04:05"),
			Anexos: urls,
		})
	}
	return resp, nil
}

func(uc *VagaUseCase) ListarPorLocalizacao(ctx context.Context, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]VagaResponse, error) {
	vagas, err := uc.vagaRepo.FindByLocation(ctx, latitude, longitude, radius, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	var vagasIDs []uint
	for _, vaga := range vagas {
		vagasIDs = append(vagasIDs, vaga.ID)
	}

	anexos, err := uc.anexoImagemRepo.FindByVagaIDs(ctx, vagasIDs)
	if err != nil {
		return nil, err
	}
	anexosPorVagaMap := make(map[uint][]string)
	for _, anexo := range anexos {
		anexosPorVagaMap[*anexo.VagaID] = append(anexosPorVagaMap[*anexo.VagaID], anexo.URL)
	}
	var resp []VagaResponse
	for _, vaga := range vagas {
		clienteNome := ""
		if vaga.Cliente != nil {
			clienteNome = vaga.Cliente.Usuario.Nome
		}
		urls := anexosPorVagaMap[vaga.ID]
		resp = append(resp, VagaResponse{
			ID:          vaga.ID,
			Titulo:      vaga.Titulo,
			Descricao:   vaga.Descricao,
			Localizacao: vaga.Localizacao,
			Latitude:    vaga.Latitude,
			Longitude:   vaga.Longitude,
			Status:      string(vaga.Status),
			Preco:       vaga.Preco,
			Urgente:     vaga.Urgente,
			Cliente:     clienteNome,
			DataCriacao: vaga.CreatedAt.Format("2006-01-02 15:04:05"),
			Anexos: urls,
		})
	}
	return resp, nil
}