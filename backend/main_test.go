package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/login", login)
	protected := router.Group("/api")
	protected.Use(autenticarJWT())
	{
		protected.GET("/livros/:id", getLivroPorID)
		protected.POST("/livros", autenticarAdmin(), adicionarLivro)
		protected.DELETE("/livros/:id", autenticarAdmin(), removerLivro)
		protected.POST("/recommend", recomendarLivros)
	}
	return router
}

// ==========================
// TESTES DE FUNÇÕES AUXILIARES
// ==========================

func TestGerarToken(t *testing.T) {
	token, err := gerarToken("user1", "user")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidarTokenInvalido(t *testing.T) {
	_, err := validarToken("token_invalido")
	assert.Error(t, err)
}

func TestCarregarDadosCSV(t *testing.T) {
	// Cria um arquivo CSV de teste
	file, _ := os.CreateTemp("", "livros_test.csv")
	defer os.Remove(file.Name())
	file.WriteString("ID,Title,Author,Main Genre,Rating\n1,Book 1,Author 1,Fiction,4.5\n")
	file.Close()

	err := carregarDadosCSV(file.Name())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(livros))
	assert.Equal(t, "Book 1", livros[0].Title)
}

func TestSalvarDadosCSV(t *testing.T) {
	// Configura dados para salvar
	livros = []livro{
		{ID: 1, Title: "Book 1", Author: "Author 1", MainGenre: "Fiction", Rating: 4.5},
	}

	file, _ := os.CreateTemp("", "livros_test.csv")
	defer os.Remove(file.Name())

	err := salvarDadosCSV(file.Name())
	assert.NoError(t, err)

	// Verifica se o arquivo contém os dados corretos
	content, _ := os.ReadFile(file.Name())
	assert.Contains(t, string(content), "Book 1")
}

// ==========================
// TESTES DE ENDPOINTS HTTP
// ==========================

func TestLoginEndpoint(t *testing.T) {
	router := setupRouter()

	body := `{"usuario_id":"admin","senha":"senha123"}`
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "token")
}

func TestGetLivroPorID(t *testing.T) {
	router := setupRouter()

	// Configura livros para o teste
	livros = []livro{
		{ID: 1, Title: "Test Book", Author: "Author", MainGenre: "Fiction", Rating: 4.5},
	}

	req, _ := http.NewRequest("GET", "/api/livros/1", nil)
	token, _ := gerarToken("user1", "user")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Test Book")
}

func TestAdicionarLivro(t *testing.T) {
	router := setupRouter()

	body := `{"title":"New Book","author":"New Author","main_genre":"Fiction","rating":4.8}`
	req, _ := http.NewRequest("POST", "/api/livros", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	token, _ := gerarToken("admin", "admin")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "New Book")
}

func TestRemoverLivro(t *testing.T) {
	router := setupRouter()

	livros = []livro{
		{ID: 1, Title: "Test Book", Author: "Author", MainGenre: "Fiction", Rating: 4.5},
	}

	req, _ := http.NewRequest("DELETE", "/api/livros/1", nil)
	token, _ := gerarToken("admin", "admin")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Livro removido com sucesso")
}

func TestRecomendarLivros(t *testing.T) {
	router := setupRouter()

	livros = []livro{
		{ID: 1, Title: "Book A", Author: "Author A", MainGenre: "Biographies", Rating: 4.5},
		{ID: 2, Title: "Book B", Author: "Author B", MainGenre: "Engineering", Rating: 4.7},
		{ID: 3, Title: "Book C", Author: "Author C", MainGenre: "Politics", Rating: 4.3},
	}

	body := `{"generos":["Biographies","Engineering","Politics"]}`
	req, _ := http.NewRequest("POST", "/api/recommend", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	token, _ := gerarToken("user1", "user")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Book A")
	assert.Contains(t, resp.Body.String(), "Book B")
	assert.Contains(t, resp.Body.String(), "Book C")
}