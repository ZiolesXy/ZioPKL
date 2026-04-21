package helper

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"

	"server/internal/domain/models"
)

const (
	ContextUserKey   = "current_user"
	ContextClaimsKey = "clerk_claims"
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
	MinIOUseSSL          bool
	ClerkJWKSURL         string
	ClerkIssuer          string
	ClerkAudience        string
	ClerkSecretKey       string
	ClerkAPIURL          string
	StorageBaseURL       string
}

type ClerkClaims struct {
	ClerkID string `json:"sub"`
	Email   string `json:"email"`
	jwt.RegisteredClaims
}

type ClerkVerifier struct {
	jwks     keyfunc.Keyfunc
	issuer   string
	audience string
}

type ClerkClient struct {
	baseURL    string
	secretKey  string
	httpClient *http.Client
}

type ClerkUserProfile struct {
	ID                    string `json:"id"`
	Username              string `json:"username"`
	PrimaryEmailAddressID string `json:"primary_email_address_id"`
	EmailAddresses        []struct {
		ID           string `json:"id"`
		EmailAddress string `json:"email_address"`
	} `json:"email_addresses"`
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
		MinIOUseSSL:          getEnvAsBool("MINIO_USE_SSL", false),
		ClerkJWKSURL:         getEnv("CLERK_JWKS_URL", ""),
		ClerkIssuer:          getEnv("CLERK_ISSUER", ""),
		ClerkAudience:        getEnv("CLERK_AUDIENCE", ""),
		ClerkSecretKey:       getEnv("CLERK_SECRET_KEY", ""),
		ClerkAPIURL:          getEnv("CLERK_API_URL", "https://api.clerk.com"),
		StorageBaseURL:       getEnv("STORAGE_BASE_URL", ""),
	}
}

func NewClerkVerifier(cfg Config) (*ClerkVerifier, error) {
	if strings.TrimSpace(cfg.ClerkJWKSURL) == "" {
		return nil, errors.New("missing CLERK_JWKS_URL")
	}

	jwks, err := keyfunc.NewDefaultCtx(context.Background(), []string{cfg.ClerkJWKSURL})
	if err != nil {
		return nil, err
	}

	return &ClerkVerifier{
		jwks:     jwks,
		issuer:   cfg.ClerkIssuer,
		audience: cfg.ClerkAudience,
	}, nil
}

func NewClerkClient(cfg Config) *ClerkClient {
	return &ClerkClient{
		baseURL:    strings.TrimRight(strings.TrimSpace(cfg.ClerkAPIURL), "/"),
		secretKey:  strings.TrimSpace(cfg.ClerkSecretKey),
		httpClient: HTTPClient(),
	}
}

func (v *ClerkVerifier) VerifyToken(tokenString string) (*ClerkClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ClerkClaims{}, v.jwks.Keyfunc)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*ClerkClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid claims")
	}
	if v.issuer != "" && claims.Issuer != v.issuer {
		return nil, errors.New("invalid issuer")
	}
	if v.audience != "" && !containsAudience(claims.Audience, v.audience) {
		return nil, errors.New("invalid audience")
	}
	if strings.TrimSpace(claims.ClerkID) == "" {
		return nil, errors.New("missing token claims")
	}
	return claims, nil
}

func (c *ClerkClient) FetchPrimaryEmail(clerkUserID string) (string, error) {
	profile, err := c.FetchUserProfile(clerkUserID)
	if err != nil {
		return "", err
	}
	if profile == nil {
		return "", nil
	}

	for _, item := range profile.EmailAddresses {
		if item.ID == profile.PrimaryEmailAddressID && strings.TrimSpace(item.EmailAddress) != "" {
			return strings.TrimSpace(item.EmailAddress), nil
		}
	}
	for _, item := range profile.EmailAddresses {
		if strings.TrimSpace(item.EmailAddress) != "" {
			return strings.TrimSpace(item.EmailAddress), nil
		}
	}

	return "", nil
}

func (c *ClerkClient) FetchUsername(clerkUserID string) (string, error) {
	profile, err := c.FetchUserProfile(clerkUserID)
	if err != nil {
		return "", err
	}
	if profile == nil {
		return "", nil
	}
	return strings.TrimSpace(profile.Username), nil
}

func (c *ClerkClient) FetchUserProfile(clerkUserID string) (*ClerkUserProfile, error) {
	if strings.TrimSpace(clerkUserID) == "" {
		return nil, errors.New("missing clerk user id")
	}
	if c == nil || c.secretKey == "" {
		return nil, nil
	}

	endpoint := fmt.Sprintf("%s/v1/users/%s", c.baseURL, url.PathEscape(clerkUserID))
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("clerk api returned status %d", resp.StatusCode)
	}

	var payload ClerkUserProfile
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func (c *ClerkClient) FetchUserIDByEmail(email string) (string, error) {
	if strings.TrimSpace(email) == "" {
		return "", errors.New("missing email")
	}
	if c == nil || c.secretKey == "" {
		return "", errors.New("missing CLERK_SECRET_KEY")
	}

	endpoint := fmt.Sprintf("%s/v1/users?email_address=%s&limit=1", c.baseURL, url.QueryEscape(strings.TrimSpace(email)))
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("clerk api returned status %d", resp.StatusCode)
	}

	var payload []struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	if len(payload) == 0 || strings.TrimSpace(payload[0].ID) == "" {
		return "", fmt.Errorf("clerk user not found for email %s", email)
	}

	return strings.TrimSpace(payload[0].ID), nil
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

func containsAudience(audiences jwt.ClaimStrings, audience string) bool {
	for _, value := range audiences {
		if value == audience {
			return true
		}
	}
	return false
}

func HTTPClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second}
}
