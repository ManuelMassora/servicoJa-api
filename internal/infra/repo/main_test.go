package repo_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/db"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

// TestDB é a conexão compartilhada por todos os testes de integração deste pacote.
var TestDB *gorm.DB

// TestMain é o ponto de entrada para os testes de integração do pacote repo_test.
// Ele sobe um container Postgres, roda as migrations e seed de roles uma única vez,
// executa todos os testes e por fim encerra o container.
func TestMain(m *testing.M) {
	ctx := context.Background()

	// ── Variáveis de ambiente para os testes ──────────────────────────────────
	os.Setenv("APP_ENV", "test")

	// ── Container Postgres ────────────────────────────────────────────────────
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("servicoja_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("TestMain: falha ao iniciar container postgres: %v", err)
	}

	// Garante que o container será encerrado ao final dos testes.
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Printf("TestMain: falha ao encerrar container: %v", err)
		}
	}()

	// ── Conexão + Migrations ──────────────────────────────────────────────────
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("TestMain: falha ao obter connection string: %v", err)
	}

	// InitDB já executa: createEnums, autoMigrate e insertInitialRoles.
	gormDB, err := db.InitDB(connStr)
	if err != nil {
		log.Fatalf("TestMain: falha ao inicializar banco de dados: %v", err)
	}

	TestDB = gormDB

	// ── Executa todos os testes do pacote ─────────────────────────────────────
	code := m.Run()
	os.Exit(code)
}