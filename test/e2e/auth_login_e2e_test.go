package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain/entities"
	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setupTestDatabase cria uma conexão com o banco de dados de teste
func setupTestDatabase(tb testing.TB) (*mongo.Database, func()) {
	// Usar variáveis de ambiente diretamente em vez de configs.NewConfig()
	connectionURI := os.Getenv("MONGODB_CONNECTION_URI")
	if connectionURI == "" {
		connectionURI = "mongodb://admin:password123@localhost:27017/"
	}

	databaseName := os.Getenv("MONGODB_DATABASE_NAME")
	if databaseName == "" {
		databaseName = "aetheris_test"
	}

	// Conectar ao MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionURI))
	require.NoError(tb, err)

	// Verificar conexão
	err = client.Ping(context.Background(), nil)
	require.NoError(tb, err)

	db := client.Database(databaseName)

	// Função de cleanup
	cleanup := func() {
		// Limpar todas as coleções de teste
		collections := []string{"users", "otps", "clients", "authorization_codes", "refresh_tokens"}
		for _, collectionName := range collections {
			_, err := db.Collection(collectionName).DeleteMany(context.Background(), bson.M{})
			if err != nil {
				tb.Logf("Erro ao limpar coleção %s: %v", collectionName, err)
			}
		}
		client.Disconnect(context.Background())
	}

	return db, cleanup
}

// createTestUser cria um usuário de teste no banco de dados
func createTestUser(tb testing.TB, db *mongo.Database, email string) primitive.ObjectID {
	user := entities.User{
		ID:        primitive.NewObjectID(),
		FirstName: "Test",
		LastName:  "User",
		Email:     email,
		CreatedAt: time.Now(),
	}

	_, err := db.Collection("users").InsertOne(context.Background(), user)
	require.NoError(tb, err)

	return user.ID
}

// deleteTestUser remove um usuário de teste do banco de dados
func deleteTestUser(tb testing.TB, db *mongo.Database, userID primitive.ObjectID) {
	_, err := db.Collection("users").DeleteOne(context.Background(), bson.M{"_id": userID})
	require.NoError(tb, err)
}

// getAPIBaseURL retorna a URL base da API
func getAPIBaseURL() string {
	apiURL := os.Getenv("API_BASE_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}
	return apiURL
}

// waitForAppReady aguarda a aplicação estar pronta
func waitForAppReady(tb testing.TB, maxRetries int) {
	// Aguardar um tempo fixo para a aplicação inicializar
	time.Sleep(5 * time.Second)
	tb.Log("Aguardando aplicação inicializar...")
}

func TestLoginE2E(t *testing.T) {
	waitForAppReady(t, 10)
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	t.Run("should return success when valid email is provided", func(t *testing.T) {
		// Arrange
		email := fmt.Sprintf("test_user_%d@example.com", time.Now().Unix())
		userID := createTestUser(t, db, email)
		defer deleteTestUser(t, db, userID)

		// Preparar payload
		payload := models.LoginPayload{
			Email: email,
		}
		jsonPayload, err := json.Marshal(payload)
		require.NoError(t, err)

		// Fazer requisição HTTP real
		apiURL := getAPIBaseURL()
		resp, err := http.Post(
			apiURL+"/api/v1/auth/login",
			"application/json",
			bytes.NewReader(jsonPayload),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verificar resposta
		var response models.SendVerificationCodeResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.NotEmpty(t, response.OTPToken)
		assert.True(t, response.ExpiresAt.After(time.Now()))

		// Verificar se o cookie foi definido
		cookies := resp.Cookies()
		assert.NotEmpty(t, cookies, "Esperado que pelo menos um cookie seja definido")
	})

	t.Run("should return not found when user does not exist", func(t *testing.T) {
		// Arrange
		email := "nonexistent@example.com"

		// Preparar payload
		payload := models.LoginPayload{
			Email: email,
		}
		jsonPayload, err := json.Marshal(payload)
		require.NoError(t, err)

		// Fazer requisição HTTP real
		apiURL := getAPIBaseURL()
		resp, err := http.Post(
			apiURL+"/api/v1/auth/login",
			"application/json",
			bytes.NewReader(jsonPayload),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return bad request when email is invalid", func(t *testing.T) {
		// Arrange
		payload := models.LoginPayload{
			Email: "invalid-email",
		}
		jsonPayload, err := json.Marshal(payload)
		require.NoError(t, err)

		// Fazer requisição HTTP real
		apiURL := getAPIBaseURL()
		resp, err := http.Post(
			apiURL+"/api/v1/auth/login",
			"application/json",
			bytes.NewReader(jsonPayload),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("should return bad request when email is empty", func(t *testing.T) {
		// Arrange
		payload := models.LoginPayload{
			Email: "",
		}
		jsonPayload, err := json.Marshal(payload)
		require.NoError(t, err)

		// Fazer requisição HTTP real
		apiURL := getAPIBaseURL()
		resp, err := http.Post(
			apiURL+"/api/v1/auth/login",
			"application/json",
			bytes.NewReader(jsonPayload),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("should return bad request when payload is invalid JSON", func(t *testing.T) {
		// Fazer requisição HTTP real com JSON inválido
		apiURL := getAPIBaseURL()
		resp, err := http.Post(
			apiURL+"/api/v1/auth/login",
			"application/json",
			bytes.NewReader([]byte("invalid json")),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return method not allowed for GET request", func(t *testing.T) {
		// Fazer requisição HTTP GET
		apiURL := getAPIBaseURL()
		resp, err := http.Get(apiURL + "/api/v1/auth/login")
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})
}

// TestLoginE2EWithRepository testa usando o repositório diretamente
func TestLoginE2EWithRepository(t *testing.T) {
	waitForAppReady(t, 10)
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	// Criar repositório
	userRepo := repositories.NewUserRepository(db)

	t.Run("should create user and then login successfully", func(t *testing.T) {
		// Arrange - Criar usuário via repositório
		email := fmt.Sprintf("repo_test_%d@example.com", time.Now().Unix())
		user := &entities.User{
			FirstName: "Repository",
			LastName:  "Test",
			Email:     email,
		}

		err := userRepo.Create(context.Background(), user)
		require.NoError(t, err)
		defer deleteTestUser(t, db, user.ID)

		// Verificar se usuário foi criado
		foundUser, err := userRepo.FindByEmail(context.Background(), email)
		require.NoError(t, err)
		assert.Equal(t, email, foundUser.Email)

		// Act - Fazer login via HTTP
		payload := models.LoginPayload{
			Email: email,
		}
		jsonPayload, err := json.Marshal(payload)
		require.NoError(t, err)

		apiURL := getAPIBaseURL()
		resp, err := http.Post(
			apiURL+"/api/v1/auth/login",
			"application/json",
			bytes.NewReader(jsonPayload),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response models.SendVerificationCodeResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.NotEmpty(t, response.OTPToken)
		assert.True(t, response.ExpiresAt.After(time.Now()))
	})
}

// Benchmark para testar performance do endpoint
func BenchmarkLoginE2E(b *testing.B) {
	db, cleanup := setupTestDatabase(b)
	defer cleanup()

	email := fmt.Sprintf("benchmark_%d@example.com", time.Now().Unix())
	userID := createTestUser(b, db, email)
	defer deleteTestUser(b, db, userID)

	payload := models.LoginPayload{
		Email: email,
	}
	jsonPayload, _ := json.Marshal(payload)
	apiURL := getAPIBaseURL()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Post(
			apiURL+"/api/v1/auth/login",
			"application/json",
			bytes.NewReader(jsonPayload),
		)
		if err == nil {
			resp.Body.Close()
		}
	}
}
