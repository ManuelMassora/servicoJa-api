package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/config"
	"github.com/ManuelMassora/servicoJa-api/internal/db"
	"github.com/ManuelMassora/servicoJa-api/internal/di"
	"github.com/ManuelMassora/servicoJa-api/internal/middleware"
	"github.com/ManuelMassora/servicoJa-api/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

var (
	testDB     *gorm.DB
	testRouter *gin.Engine
	testAppCfg *config.Config
	testCtx    context.Context
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// ── Variáveis de ambiente para os testes ──────────────────────────────────
	os.Setenv("APP_ENV", "test")

	// ── Container Postgres ────────────────────────────────────────────────────
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("servicoja_e2e_test"),
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

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Printf("TestMain: falha ao encerrar container: %v", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("TestMain: falha ao obter connection string: %v", err)
	}

	gormDB, err := db.InitDB(connStr)
	if err != nil {
		log.Fatalf("TestMain: falha ao inicializar banco de dados: %v", err)
	}

	testDB = gormDB

	// ── Setup do Router ───────────────────────────────────────────────────────
	// Usamos configurações fakes para o teste
	testAppCfg = &config.Config{
		DatabaseDSN:  connStr,
		JwtSecretKey: "test_secret_key",
		ServerPort:   "8080",
		ServerHost:   "localhost",
		SupabaseURL:  "http://localhost",
		SupabaseKey:  "key",
	}

	gin.SetMode(gin.TestMode)
	testRouter = gin.Default()
	jwtService := middleware.NewJWTService(testAppCfg.JwtSecretKey)
	container := di.NewContainer(testDB, testAppCfg, jwtService)
	routes.SetRoutes(testRouter, container, jwtService)

	code := m.Run()
	os.Exit(code)
}

func doRequest(method, url string, body io.Reader, headers map[string]string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, body)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

func registerUser(role, email, telefone, senha string) {
	var url string
	switch role {
	case "CLIENTE":
		url = "/iniciar/cliente"
	case "PRESTADOR":
		url = "/iniciar/prestador"
	case "ADMIN":
		url = "/iniciar/admin"
	}

	fields := map[string]string{
		"nome":     "Test User",
		"email":    email,
		"senha":    senha,
		"telefone": telefone,
	}

	if role == "PRESTADOR" {
		fields["localizacao"] = "Maputo"
		fields["latitude"] = "-25.9692"
		fields["longitude"] = "32.5732"
	}

	if role == "ADMIN" {
		jsonBody, _ := json.Marshal(fields)
		doRequest("POST", url, bytes.NewBuffer(jsonBody), map[string]string{"Content-Type": "application/json"})
	} else {
		doMultipartRequest("POST", url, fields, "")
	}
}

func loginUser(telefone, senha string) string {
	body := map[string]string{
		"telefone": telefone,
		"senha":    senha,
	}
	jsonBody, _ := json.Marshal(body)
	w := doRequest("POST", "/login", bytes.NewBuffer(jsonBody), map[string]string{"Content-Type": "application/json"})

	if w.Code != http.StatusAccepted {
		log.Printf("Login failed for %s: %s", telefone, w.Body.String())
		return ""
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if token, ok := resp["token"].(string); ok {
		return token
	}
	return ""
}

func doMultipartRequest(method, url string, fields map[string]string, token string) *httptest.ResponseRecorder {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, val := range fields {
		_ = writer.WriteField(key, val)
	}
	_ = writer.Close()

	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}
