package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Chave secreta usada para assinar o token JWT
var jwtKey = []byte("minha_chave_secreta")

// Mutex para sincronização de acesso às estruturas de dados
var mutex sync.Mutex

// gerarToken cria um token JWT com as claims fornecidas.
func gerarToken(usuarioID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub": usuarioID, // Identificador do usuário
		"role": role,    // Papel do usuário
		"exp": time.Now().Add(1 * time.Hour).Unix(), // Expira em 1 hora
		"iat": time.Now().Unix(),                   // Emitido no momento atual
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// validarToken valida um token JWT e retorna as claims, se válido.
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

// autenticarJWT é um middleware que valida o token JWT presente no cabeçalho da requisição.
func autenticarJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obter o token do cabeçalho Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
			c.Abort()
			return
		}

		// Remover "Bearer " do início do cabeçalho
		tokenString := authHeader[len("Bearer "):] // Assume o prefixo "Bearer "

		// Validar o token
		claims, err := validarToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Adicionar claims ao contexto para uso posterior
		c.Set("claims", claims)
		c.Next()
	}
}

// autenticarRole é um middleware que verifica se o usuário tem o papel necessário.
func autenticarRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsRaw, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Acesso não autorizado"})
			c.Abort()
			return
		}

		claims, ok := claimsRaw.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato inválido de claims"})
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permissão insuficiente"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// usuarios é um mapa que simula um banco de dados de usuários, senhas e papéis.
var usuarios = map[string]struct {
	Senha string
	Role  string
}{
	"admin":      {"senha123", "admin"},
	"joanderson": {"123456", "user"},
	"user2":      {"outrasenha", "user"},
}

// login realiza a autenticação de um usuário e retorna um token JWT se as credenciais forem válidas.
func login(c *gin.Context) {
	var credenciais struct {
		UsuarioID string `json:"usuario_id"` // Identificador do usuário
		Senha     string `json:"senha"`     // Senha do usuário
	}

	// Obter credenciais do corpo da requisição
	if err := c.ShouldBindJSON(&credenciais); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Verificar se o usuário e senha existem
	usuario, existe := usuarios[credenciais.UsuarioID]
	if !existe || usuario.Senha != credenciais.Senha {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	// Gerar token JWT
	token, err := gerarToken(credenciais.UsuarioID, usuario.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// livrorepresenta os dados sobre um livro musical.
type livro struct {
	ID      string  `json:"id"`      // Identificador do livro
	Title   string  `json:"title"`   // Título do livro
	Author  string  `json:"author"`  // Autor do livro
	Price   float64 `json:"price"`   // Preço do livro
	Rating  float64 `json:"rating"`  // Avaliação geral do livro
	Ratings []float64 // Lista de avaliações recebidas
}

// livros é um slice que armazena dados de álbuns para simular um banco de dados.
var livros = []livro{
	{ID: "1", Title: "1984", Author: "George Orwell", Price: 29.90, Ratings: []float64{4.5, 4.8, 5.0}},
	{ID: "2", Title: "To Kill a Mockingbird", Author: "Harper Lee", Price: 35.50, Ratings: []float64{4.8, 4.7, 4.9}},
	{ID: "3", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Price: 27.99, Ratings: []float64{4.4, 4.2, 4.6}},
	{ID: "4", Title: "Moby Dick", Author: "Herman Melville", Price: 50.00, Ratings: []float64{4.1, 4.3, 4.2}},
	{ID: "5", Title: "Pride and Prejudice", Author: "Jane Austen", Price: 39.99, Ratings: []float64{4.6, 4.7, 4.8}},
	{ID: "6", Title: "The Catcher in the Rye", Author: "J.D. Salinger", Price: 32.00, Ratings: []float64{4.3, 4.4, 4.5}},
	{ID: "7", Title: "Brave New World", Author: "Aldous Huxley", Price: 28.99, Ratings: []float64{4.5, 4.6, 4.7}},
	{ID: "8", Title: "The Hobbit", Author: "J.R.R. Tolkien", Price: 45.00, Ratings: []float64{4.9, 5.0, 4.8}},
	{ID: "9", Title: "Crime and Punishment", Author: "Fyodor Dostoevsky", Price: 49.90, Ratings: []float64{4.8, 4.7, 4.9}},
	{ID: "10", Title: "The Alchemist", Author: "Paulo Coelho", Price: 31.50, Ratings: []float64{4.2, 4.4, 4.3}},
	{ID: "11", Title: "War and Peace", Author: "Leo Tolstoy", Price: 59.99, Ratings: []float64{4.7, 4.8, 4.9}},
	{ID: "12", Title: "The Lord of the Rings", Author: "J.R.R. Tolkien", Price: 79.90, Ratings: []float64{5.0, 4.9, 5.0}},
}


// getlivros retorna a lista de álbuns como JSON.
func getlivros(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, livros)
}

// postlivros adiciona um novo livro com base nos dados fornecidos na requisição.
func postlivros(c *gin.Context) {
	var newLivro livro

	// Lê os dados JSON e os vincula à estrutura newLivro
	if err := c.BindJSON(&newLivro); err != nil {
		return
	}

	// Adiciona o novo livro à lista
	mutex.Lock()
	livros = append(livros, newLivro)
	mutex.Unlock()
	c.IndentedJSON(http.StatusCreated, newLivro)
}

// ratelivropermite que o usuário avalie um livro pelo ID.
func rateLivro(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Rating float64 `json:"rating"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	if input.Rating < 1 || input.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A avaliação deve estar entre 1 e 5"})
		return
	}

	// Atualizar o livro com a nova avaliação
	mutex.Lock()
	for i, a := range livros {
		if a.ID == id {
			a.Ratings = append(a.Ratings, input.Rating)
			total := 0.0
			for _, r := range a.Ratings {
				total += r
			}
			a.Rating = total / float64(len(a.Ratings))
			livros[i] = a
			mutex.Unlock()
			c.JSON(http.StatusOK, gin.H{"message": "Avaliação registrada", "livro": a})
			return
		}
	}
	mutex.Unlock()

	c.JSON(http.StatusNotFound, gin.H{"error": "livro não encontrado"})
}

func calcularMedia(id string) (float64, error) {
    mutex.Lock() // Protege o acesso aos dados compartilhados
    defer mutex.Unlock()

    for _, a := range livros {
        if a.ID == id {
            if len(a.Ratings) == 0 {
                return 0.0, nil // Retorna 0 se não houver avaliações
            }
            total := 0.0
            for _, r := range a.Ratings {
                total += r
            }
            return total / float64(len(a.Ratings)), nil
        }
    }
    return 0.0, fmt.Errorf("livro não encontrado")
}


// getLivroRatings retorna a média de avaliações de um livro.
func getLivroRatings(c *gin.Context) {
    id := c.Param("id")

    // Calcular a média de avaliações do livro usando a função calcularMedia
    media, err := calcularMedia(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "livro não encontrado"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "livro_id":       id,
        "average_rating": media,
    })
}

func main() {
	router := gin.Default()

	// Rota pública
	router.POST("/login", login)

	// Grupo de rotas protegidas por JWT
	protected := router.Group("/api")
	protected.Use(autenticarJWT())
	{
		protected.GET("/livros", getlivros) // Aberto para todos os usuários autenticados
		protected.POST("/livros", autenticarRole("admin"), postlivros) // Apenas administradores
		protected.GET("/livros/:id", getLivroRatings) // Média de avaliações de um livro
		protected.POST("/livros/:id/rate", rateLivro) // Avaliar um livro
	}

	// Inicia o servidor na porta

	router.Run("localhost:8080")
}
