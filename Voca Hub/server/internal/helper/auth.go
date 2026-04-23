package helper

import (
	"archive/zip"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"server/internal/domain/models"
)

const (
	ContextUserKey   = "current_user"
	ContextClaimsKey = "auth_claims"
	ContextTokenKey  = "access_token"
)

type Config struct {
	AppPort              string
	CORSOrigins          string
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	DBSSLMode            string
	DBTimeZone           string
	RedisHost            string
	RedisPort            string
	RedisPassword        string
	RedisDB              int
	MinIOEndpoint        string
	MinIOAccessKey       string
	MinIOSecretKey       string
	MinIOBucket          string
	MinIOThumbnailBucket string
	MinIOProfileBucket   string
	MinIOUseSSL          bool
	StorageBaseURL       string
	JWTSecret            string
	AccessTokenTTL       time.Duration
	RefreshTokenTTL      time.Duration
	SeedUserPassword     string
}

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type AuthClaims struct {
	UserID uint      `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

type TokenPair struct {
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	AccessTokenExpiresAt  string `json:"access_token_expires_at"`
	RefreshTokenExpiresAt string `json:"refresh_token_expires_at"`
}

func LoadConfig() Config {
	_ = godotenv.Load()

	return Config{
		AppPort:              getEnv("APP_PORT", "8080"),
		CORSOrigins:          getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173"),
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBUser:               getEnv("DB_USER", "postgres"),
		DBPassword:           getEnv("DB_PASSWORD", "postgres"),
		DBName:               getEnv("DB_NAME", "voca_hub"),
		DBSSLMode:            getEnv("DB_SSLMODE", "disable"),
		DBTimeZone:           getEnv("DB_TIMEZONE", "UTC"),
		RedisHost:            getEnv("REDIS_HOST", "localhost"),
		RedisPort:            getEnv("REDIS_PORT", "6379"),
		RedisPassword:        getEnv("REDIS_PASSWORD", ""),
		RedisDB:              getEnvAsInt("REDIS_DB", 0),
		MinIOEndpoint:        getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey:       getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey:       getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinIOBucket:          getEnv("MINIO_BUCKET", "games"),
		MinIOThumbnailBucket: getEnv("MINIO_THUMBNAIL_BUCKET", "thumbnails"),
		MinIOProfileBucket:   getEnv("MINIO_PROFILE_BUCKET", "profiles"),
		MinIOUseSSL:          getEnvAsBool("MINIO_USE_SSL", false),
		StorageBaseURL:       getEnv("STORAGE_BASE_URL", ""),
		JWTSecret:            getEnv("JWT_SECRET", "change-me"),
		AccessTokenTTL:       getEnvAsDuration("ACCESS_TOKEN_TTL", 15*time.Minute),
		RefreshTokenTTL:      getEnvAsDuration("REFRESH_TOKEN_TTL", 7*24*time.Hour),
		SeedUserPassword:     getEnv("SEED_USER_PASSWORD", "password123"),
	}
}

func NewTokenManager(cfg Config) (*TokenManager, error) {
	if strings.TrimSpace(cfg.JWTSecret) == "" {
		return nil, errors.New("missing JWT_SECRET")
	}

	return &TokenManager{
		secret:          []byte(cfg.JWTSecret),
		accessTokenTTL:  cfg.AccessTokenTTL,
		refreshTokenTTL: cfg.RefreshTokenTTL,
	}, nil
}

func (m *TokenManager) GenerateTokenPair(user *models.User) (*TokenPair, error) {
	now := time.Now().UTC()
	accessExpiresAt := now.Add(m.accessTokenTTL)
	refreshExpiresAt := now.Add(m.refreshTokenTTL)

	accessClaims := AuthClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		Type:   TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(user.ID), 10),
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	refreshClaims := AuthClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		Type:   TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(user.ID), 10),
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        randomTokenID(),
		},
	}

	accessToken, err := m.sign(accessClaims)
	if err != nil {
		return nil, err
	}

	refreshToken, err := m.sign(refreshClaims)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessExpiresAt.Format(time.RFC3339),
		RefreshTokenExpiresAt: refreshExpiresAt.Format(time.RFC3339),
	}, nil
}

func (m *TokenManager) VerifyAccessToken(tokenString string) (*AuthClaims, error) {
	return m.verifyToken(tokenString, TokenTypeAccess)
}

func (m *TokenManager) VerifyRefreshToken(tokenString string) (*AuthClaims, error) {
	return m.verifyToken(tokenString, TokenTypeRefresh)
}

func (m *TokenManager) sign(claims AuthClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *TokenManager) verifyToken(tokenString string, expectedType TokenType) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid claims")
	}
	if claims.Type != expectedType {
		return nil, errors.New("invalid token type")
	}
	if claims.UserID == 0 {
		return nil, errors.New("missing token claims")
	}

	return claims, nil
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func ComparePassword(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func randomTokenID() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return base64.RawURLEncoding.EncodeToString(bytes)
}

func ExtractBearerToken(header string) (string, error) {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
		return "", errors.New("missing bearer token")
	}
	return strings.TrimSpace(parts[1]), nil
}

func MustCurrentUser(ctx interface{ MustGet(string) any }) *models.User {
	return ctx.MustGet(ContextUserKey).(*models.User)
}

func CopyMultipartFile(src multipartFile, destination string) error {
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

type multipartFile interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

func ExtractZip(src string, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		targetPath := filepath.Join(dest, file.Name)
		if !strings.HasPrefix(filepath.Clean(targetPath), filepath.Clean(dest)) {
			return errors.New("invalid zip path")
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
			return err
		}

		in, err := file.Open()
		if err != nil {
			return err
		}

		out, err := os.Create(targetPath)
		if err != nil {
			in.Close()
			return err
		}

		if _, err = io.Copy(out, in); err != nil {
			out.Close()
			in.Close()
			return err
		}

		out.Close()
		in.Close()
	}
	return nil
}

func Slugify(input string) string {
	normalized := strings.ToLower(strings.TrimSpace(input))
	normalized = strings.ReplaceAll(normalized, " ", "-")
	normalized = strings.ReplaceAll(normalized, "_", "-")
	normalized = strings.ReplaceAll(normalized, "/", "-")
	return fmt.Sprintf("%s-%d", normalized, time.Now().Unix())
}

func getEnv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getEnvAsInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	var result int
	_, err := fmt.Sscanf(value, "%d", &result)
	if err != nil {
		return fallback
	}
	return result
}

func getEnvAsBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return strings.EqualFold(value, "true") || value == "1"
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return duration
}

func HTTPClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second}
}
