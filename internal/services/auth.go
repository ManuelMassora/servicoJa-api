package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/ManuelMassora/servicoJa-api/internal/middleware"
	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/alexedwards/argon2id"
)

type AuthRequest struct {
	Telefone 	string		`json:"telefone" binding:"required"`
	Senha 		string		`json:"senha" binding:"required"`
}

type AuthResponse struct {
	ID 			uint		`json:"id"`
	Role 		string		`json:"role"`
	Token 		string		`json:"token"`
}

type AuthUSer struct {
	r          model.UsuarioRepo
	jwtService middleware.JwtService
}

func NewAuthUser(r model.UsuarioRepo, jwt *middleware.JwtService) *AuthUSer {
	return &AuthUSer{r: r, jwtService: *jwt}
}

func (uc *AuthUSer) Authenticate(ctx context.Context, request AuthRequest) (*AuthResponse, error) {
	user, err := uc.r.BuscarPorTelefone(ctx, request.Telefone)
	if err != nil {
		if err.Error() == "usuário não encontrado" {
			return nil, errors.New("credenciais inválidas")
		}
		return nil, err
	}
	match, err := argon2id.ComparePasswordAndHash(request.Senha, user.Senha)
	if err != nil {
		switch err {
		case argon2id.ErrInvalidHash:
			return nil, errors.New("o hash não está no formato correto")
		case argon2id.ErrIncompatibleVariant:
			return nil, errors.New("credenciais inválidas")
		case argon2id.ErrIncompatibleVersion:
			return nil, errors.New("versão incompatível do argon2")
		default:
			return nil, fmt.Errorf("erro interno de autenticação: %w", err)
		}
	}
	if !match {
		return nil, errors.New("credenciais inválidas")
	}
	token, err := uc.jwtService.GenateToken(uint(user.ID), string(user.RolePermissao.Role))
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar token de autenticacao")
	}
	response := &AuthResponse{
		ID:    user.ID,
		Role:  string(user.RolePermissao.Role),
		Token: token,
	}
	return response, nil
}

func (uc *AuthUSer) BuscarTelefone(ctx context.Context, telefone string) (*model.Usuario, error) {
	user, err := uc.r.BuscarPorTelefone(ctx, telefone)
	if err != nil {
		return nil, fmt.Errorf("usuário não encontrado")
	}
	if user == nil {
		return nil, errors.New("usuário não encontrado")
	}
	return user, nil
}