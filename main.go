package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Chave secreta usada para assinar o token JWT
var jwtKey = []byte("minha_chave_secreta")

// Mutex para sincronização de acesso às estruturas de dados
var mutex sync.Mutex

// Variável global para gerar IDs sequenciais
var nextID = 1

// Estrutura para representar um livro
type livro struct {
	ID        int     `json:"id"`        // ID do livro
	Title     string  `json:"title"`     // Título do livro
	Author    string  `json:"author"`    // Autor do livro
	MainGenre string  `json:"main_genre"` // Gênero principal
	Rating    float64 `json:"rating"`    // Avaliação do livro
}

// Lista de livros (simulação de banco de dados)
var livros []livro

// Função para carregar dados do CSV
func carregarDadosCSV(caminho string) error {
	file, err := os.Open(caminho)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("erro ao ler o arquivo CSV: %v", err)
	}

	// Processar registros, ignorando o cabeçalho
	for i, record := range records {
		if i == 0 {
			continue // Ignora o cabeçalho
		}
		rating, _ := strconv.ParseFloat(record[4], 64) // Converte o rating para float64
		id, _ := strconv.Atoi(record[0])              // Converte o ID de string para int

		// Atualiza o próximo ID se necessário
		if id >= nextID {
			nextID = id + 1
		}

		livros = append(livros, livro{
			ID:        id,
			Title:     record[1],
			Author:    record[2],
			MainGenre: record[3],
			Rating:    rating,
		})
	}
	return nil
}

// Função para salvar os livros no arquivo CSV
func salvarDadosCSV(caminho string) error {
	file, err := os.Create(caminho)
	if err != nil {
		return fmt.Errorf("erro ao criar o arquivo CSV: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escrever cabeçalho
	writer.Write([]string{"ID", "Title", "Author", "Main Genre", "Rating"})

	// Escrever livros
	for _, livro := range livros {
		record := []string{
			strconv.Itoa(livro.ID),    // Usando o ID sequencial
			livro.Title,
			livro.Author,
			livro.MainGenre,
			fmt.Sprintf("%.1f", livro.Rating),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("erro ao escrever o livro no arquivo CSV: %v", err)
		}
	}
	return nil
}

// Função para gerar o token JWT
func gerarToken(usuarioID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  usuarioID,
		"role": role,
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Função para validar o token JWT
func validarToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inválido")
		}
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("token inválido")
}

// Middleware para autenticar JWT
func autenticarJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
			c.Abort()
			return
		}
		tokenString := authHeader[len("Bearer "):]
		claims, err := validarToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

// Endpoint de login
var usuarios = map[string]struct {
	Senha string
	Role  string
}{
	"admin": {"senha123", "admin"},
	"user1": {"senha123", "user"},
}

func login(c *gin.Context) {
	var credenciais struct {
		UsuarioID string `json:"usuario_id"`
		Senha     string `json:"senha"`
	}
	if err := c.ShouldBindJSON(&credenciais); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}
	usuario, existe := usuarios[credenciais.UsuarioID]
	if !existe || usuario.Senha != credenciais.Senha {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}
	token, err := gerarToken(credenciais.UsuarioID, usuario.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Endpoint para obter livro por ID
func getLivroPorID(c *gin.Context) {
	idStr := c.Param("id") // Obtém o parâmetro ID da URL
	id, err := strconv.Atoi(idStr) // Converte o ID de string para int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verifica se o livro com o ID existe
	for _, livro := range livros {
		if livro.ID == id {
			c.IndentedJSON(http.StatusOK, livro)
			return
		}
	}

	// Se o livro não for encontrado
	c.JSON(http.StatusNotFound, gin.H{"error": "Livro não encontrado"})
}



// Middleware para verificar se o usuário é admin
func autenticarAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Acesso não autorizado"})
			c.Abort()
			return
		}

		// Verificar se o usuário tem o papel de admin
		if claims.(jwt.MapClaims)["role"] != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Endpoint para adicionar um livro (somente admin)
func adicionarLivro(c *gin.Context) {
	var novoLivro livro
	if err := c.ShouldBindJSON(&novoLivro); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Atribuir um ID sequencial ao novo livro
	novoLivro.ID = nextID
	nextID++ // Incrementa o próximo ID

	// Adicionar o novo livro
	livros = append(livros, novoLivro)
	if err := salvarDadosCSV("dados.csv"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar os dados"})
		return
	}

	c.JSON(http.StatusCreated, novoLivro)
}

// Endpoint para remover um livro (somente admin)
func removerLivro(c *gin.Context) {
	id := c.Param("id")
	for i, livro := range livros {
		if strconv.Itoa(livro.ID) == id {
			// Remover livro
			livros = append(livros[:i], livros[i+1:]...)
			if err := salvarDadosCSV("dados.csv"); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar os dados"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Livro removido com sucesso"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Livro não encontrado"})
}

func main() {
	// Carregar dados do CSV
	if err := carregarDadosCSV("dados.csv"); err != nil {
		fmt.Println("Erro ao carregar os dados do CSV:", err)
		return
	}

	router := gin.Default()

	// Rota de login
	router.POST("/login", login)

	// Grupo de rotas protegidas
	protected := router.Group("/api")
	protected.Use(autenticarJWT())
	{
		protected.GET("/livros/:id", getLivroPorID) // Lista de livros
		protected.POST("/livros", autenticarAdmin(), adicionarLivro) // Adicionar livro (somente admin)
		protected.DELETE("/livros/:id", autenticarAdmin(), removerLivro) // Remover livro (somente admin)
	}

	// Iniciar servidor
	router.Run(":8080")
}
