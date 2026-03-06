package db

import (
	"fmt"
	"log"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter sql.DB do GORM: %w", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := createEnums(db); err != nil {
		return nil, fmt.Errorf("erro criando ENUMs: %w", err)
	}

	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("erro no AutoMigrate: %w", err)
	}

	if err := insertInitialRoles(db); err != nil {
		return nil, fmt.Errorf("erro ao adicionar Roles iniciais: %w", err)
	}

	log.Println("Conexão com o banco de dados e migrações realizadas com sucesso!")
	return db, nil
}

func createEnums(db *gorm.DB) error {
	const enumSQL = `
	DO $$
	BEGIN
		-- Status Enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status') THEN
			CREATE TYPE status AS ENUM ('PENDENTE', 'EM_ANDAMENTO', 'CONCLUIDO', 'CANCELADO');
		END IF;

		-- TipoMovimento Enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tipo_movimento') THEN
			CREATE TYPE tipo_movimento AS ENUM ('CREDITO', 'DEBITO');
		END IF;

		-- MetodoPagamento Enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'metodo_pagamento') THEN
			CREATE TYPE metodo_pagamento AS ENUM ('M_PESA', 'CARTEIRA', 'OUTRO');
		END IF;

		-- Role Enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role') THEN
			CREATE TYPE role AS ENUM ('CLIENTE', 'PRESTADOR', 'ADMIN');
		END IF;
	END
	$$;
	`
	if err := db.Exec(enumSQL).Error; err != nil {
		return err
	}
	return nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// Nível 0 — sem dependências
		&model.RolePermissao{},
		&model.CategoriaPrestador{},
		&model.Categoria{},
		// Nível 1 — depende de RolePermissao
		&model.Usuario{},
		// Nível 2 — depende de Usuario
		&model.Cliente{},
		// Nível 3 — depende de Usuario, CategoriaPrestador
		&model.Prestador{},
		// Nível 4 — depende de Cliente, Prestador
		&model.Vaga{},
		// Nível 5 — depende de Prestador, Categoria
		&model.Catalogo{},
		// Nível 6 — depende de Vaga, Prestador
		&model.Proposta{},
		// Nível 7 — depende de Cliente, Catalogo
		&model.Agendamento{},
		// Nível 8 — depende de Agendamento, Vaga, Cliente, Prestador
		&model.Servico{},
		// Nível 9 — depende de Servico
		&model.Avaliacao{},
		&model.Chat{},
		&model.Mensagem{},
		&model.Notificacao{},
		&model.Pagamento{},
		&model.Transacao{},
		// Nível 10 — depende de Agendamento, Vaga, Catalogo
		&model.AnexoImagem{},
		&model.Galeria{},
		&model.Imagem{},
	)
}

func insertInitialRoles(db *gorm.DB) error {
	const insertRole = `
		INSERT INTO role_permissaos (role) VALUES 
			('CLIENTE'),
			('PRESTADOR'),
			('ADMIN')
		ON CONFLICT (role) DO NOTHING;
	`
	if err := db.Exec(insertRole).Error; err != nil {
		return err
	}
	return nil
}
