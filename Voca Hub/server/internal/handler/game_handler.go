package handler

import (
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"server/internal/domain/dto"
	"server/internal/domain/models"
	"server/internal/helper"
	"server/internal/service"
)

type GameHandler struct {
	gameService *service.GameService
}

func NewGameHandler(gameService *service.GameService) *GameHandler {
	return &GameHandler{gameService: gameService}
}

func (h *GameHandler) UploadGame(c *gin.Context) {
	var req dto.UploadGameRequest
	if err := c.ShouldBind(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "file is required")
		return
	}
	thumbnail, err := c.FormFile("thumbnail")
	if err != nil {
		thumbnail = nil
	}

	user := helper.MustCurrentUser(c)
	if gin.Mode() == gin.DebugMode {
		log.Printf("upload requested by user id=%d clerk_id=%s role=%s", user.ID, user.ClerkID, user.Role)
	}
	game, err := h.gameService.UploadGame(user.ID, req.Title, req.Description, req.CategoryIDs, req.DifficultyID, file, thumbnail)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	helper.Success(c, http.StatusCreated, "game uploaded", presentGame(game, h.gameService))
}

func (h *GameHandler) UpdateGame(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req dto.UpdateGameRequest
	if err := c.ShouldBind(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		file = nil
	}
	thumbnail, err := c.FormFile("thumbnail")
	if err != nil {
		thumbnail = nil
	}

	user := helper.MustCurrentUser(c)
	game, err := h.gameService.UpdateGame(uint(id), user, req.Title, req.Description, req.CategoryIDs, req.DifficultyID, file, thumbnail)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "game not found" {
			status = http.StatusNotFound
		} else if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		helper.Error(c, status, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "game updated", presentGame(game, h.gameService))
}

func (h *GameHandler) ListApprovedGames(c *gin.Context) {
	games, err := h.gameService.ListApprovedGames()
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "approved games fetched", helper.WrapListIfNeeded(presentGames(games, h.gameService)))
}

func (h *GameHandler) ListMyGames(c *gin.Context) {
	user := helper.MustCurrentUser(c)

	games, err := h.gameService.ListGamesByDeveloperID(user.ID)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "user games fetched", helper.WrapListIfNeeded(presentGames(games, h.gameService)))
}

func (h *GameHandler) GetApprovedGame(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	game, err := h.gameService.GetApprovedGame(uint(id))
	if err != nil {
		helper.Error(c, http.StatusNotFound, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "game fetched", presentGame(game, h.gameService))
}

func (h *GameHandler) PlayGame(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	url, err := h.gameService.PlayGame(uint(id))
	if err != nil {
		helper.Error(c, http.StatusNotFound, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "game ready", gin.H{"file_url": url})
}

func (h *GameHandler) ServeGameFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.servePlayAssetFallback(c)
		return
	}

	filePath := strings.TrimPrefix(c.Param("filepath"), "/")
	reader, contentType, err := h.gameService.OpenGameAsset(uint(id), filePath)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	defer reader.Close()

	contentType = detectContentType(filePath, contentType)
	if contentType != "" {
		c.Header("Content-Type", contentType)
	}
	c.Header("Cache-Control", "no-store")

	if isHTMLAsset(filePath, contentType) {
		body, err := io.ReadAll(reader)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Data(http.StatusOK, contentType, []byte(injectBaseHref(string(body), buildPlayBaseHref(uint(id)))))
		return
	}

	if _, err := io.Copy(c.Writer, reader); err != nil {
		c.Status(http.StatusInternalServerError)
	}
}

func (h *GameHandler) servePlayAssetFallback(c *gin.Context) {
	refererGameID, ok := gameIDFromReferer(c.GetHeader("Referer"))
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Requests like /play/Assets/img.png come from resolving ../Assets/img.png
	// against <base href="/play/:id/">. Rebuild the intended in-game asset path.
	filePath := strings.TrimPrefix(path.Join(c.Param("id"), c.Param("filepath")), "/")
	reader, contentType, err := h.gameService.OpenGameAsset(refererGameID, filePath)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	defer reader.Close()

	contentType = detectContentType(filePath, contentType)
	if contentType != "" {
		c.Header("Content-Type", contentType)
	}
	c.Header("Cache-Control", "no-store")

	if _, err := io.Copy(c.Writer, reader); err != nil {
		c.Status(http.StatusInternalServerError)
	}
}

func (h *GameHandler) ServeRootAssetFallback(c *gin.Context) {
	filePath := strings.TrimPrefix(c.Param("filepath"), "/")
	if filePath == "" || isReservedRootPath(filePath) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	gameID, ok := gameIDFromReferer(c.GetHeader("Referer"))
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	reader, contentType, err := h.gameService.OpenGameAsset(gameID, filePath)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	defer reader.Close()

	contentType = detectContentType(filePath, contentType)
	if contentType != "" {
		c.Header("Content-Type", contentType)
	}
	c.Header("Cache-Control", "no-store")

	if _, err := io.Copy(c.Writer, reader); err != nil {
		c.Status(http.StatusInternalServerError)
	}
}

func (h *GameHandler) ServeThumbnail(c *gin.Context) {
	filePath := strings.TrimPrefix(c.Param("filepath"), "/")
	reader, contentType, err := h.gameService.OpenThumbnailAsset(filePath)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	defer reader.Close()

	contentType = detectContentType(filePath, contentType)
	if contentType != "" {
		c.Header("Content-Type", contentType)
	}
	c.Header("Cache-Control", "public, max-age=3600")

	if _, err := io.Copy(c.Writer, reader); err != nil {
		c.Status(http.StatusInternalServerError)
	}
}

func isHTMLAsset(path string, contentType string) bool {
	if strings.HasPrefix(contentType, "text/html") {
		return true
	}

	ext := strings.ToLower(filepath.Ext(path))
	return ext == "" || ext == ".html" || ext == ".htm"
}

func detectContentType(path string, current string) string {
	if strings.TrimSpace(current) != "" {
		return current
	}

	switch strings.ToLower(filepath.Ext(path)) {
	case ".html", ".htm":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js", ".mjs":
		return "text/javascript; charset=utf-8"
	case ".json":
		return "application/json"
	}

	if guessed := mime.TypeByExtension(strings.ToLower(filepath.Ext(path))); guessed != "" {
		return guessed
	}

	if filepath.Ext(path) == "" {
		return "text/plain; charset=utf-8"
	}

	return "application/octet-stream"
}

func buildPlayBaseHref(id uint) string {
	return "/play/" + strconv.FormatUint(uint64(id), 10) + "/"
}

var playRefererPattern = regexp.MustCompile(`^/play/(\d+)(?:/|$)`)

func gameIDFromReferer(raw string) (uint, bool) {
	if strings.TrimSpace(raw) == "" {
		return 0, false
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return 0, false
	}

	matches := playRefererPattern.FindStringSubmatch(parsed.Path)
	if len(matches) != 2 {
		return 0, false
	}

	id, err := strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		return 0, false
	}

	return uint(id), true
}

func isReservedRootPath(path string) bool {
	cleanPath := strings.Trim(strings.ToLower(filepath.ToSlash(path)), "/")
	return cleanPath == "health" || cleanPath == "api" || strings.HasPrefix(cleanPath, "api/") || cleanPath == "play" || strings.HasPrefix(cleanPath, "play/")
}

func injectBaseHref(html string, baseHref string) string {
	baseTag := `<base href="` + baseHref + `">`
	lowerHTML := strings.ToLower(html)

	if strings.Contains(lowerHTML, "<base ") {
		return html
	}

	if idx := strings.Index(lowerHTML, "<head"); idx >= 0 {
		headEnd := strings.Index(lowerHTML[idx:], ">")
		if headEnd >= 0 {
			insertAt := idx + headEnd + 1
			return html[:insertAt] + baseTag + html[insertAt:]
		}
	}

	if idx := strings.Index(lowerHTML, "<html"); idx >= 0 {
		htmlEnd := strings.Index(lowerHTML[idx:], ">")
		if htmlEnd >= 0 {
			insertAt := idx + htmlEnd + 1
			return html[:insertAt] + "<head>" + baseTag + "</head>" + html[insertAt:]
		}
	}

	if idx := strings.Index(lowerHTML, "<head>"); idx >= 0 {
		insertAt := idx + len("<head>")
		return html[:insertAt] + baseTag + html[insertAt:]
	}

	return "<head>" + baseTag + "</head>" + html
}

func (h *GameHandler) ListCategories(c *gin.Context) {
	categories, err := h.gameService.ListCategories()
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "categories fetched", helper.WrapListIfNeeded(categories))
}

func (h *GameHandler) ListDifficulties(c *gin.Context) {
	difficulties, err := h.gameService.ListDifficulties()
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "difficulties fetched", helper.WrapListIfNeeded(difficulties))
}

func (h *GameHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.gameService.CreateCategory(req.Name)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	helper.Success(c, http.StatusCreated, "category created", category)
}

func (h *GameHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.gameService.UpdateCategory(uint(id), req.Name)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "category not found" {
			status = http.StatusNotFound
		}
		helper.Error(c, status, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "category updated", category)
}

func (h *GameHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.gameService.DeleteCategory(uint(id)); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "category not found" {
			status = http.StatusNotFound
		}
		helper.Error(c, status, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "category deleted", nil)
}

type gameResponse struct {
	ID           uint              `json:"id"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	FileURL      string            `json:"file_url"`
	ThumbnailURL string            `json:"thumbnail_url"`
	DeveloperID  uint              `json:"developer_id"`
	Status       string            `json:"status"`
	CreatedAt    string            `json:"upload_at"`
	Developer    models.User       `json:"developer"`
	Difficulty   models.Difficulty `json:"difficulty"`
	Categories   []models.Category `json:"categories"`
}

func presentGames(games []models.Game, gameService *service.GameService) []gameResponse {
	result := make([]gameResponse, 0, len(games))
	for _, game := range games {
		result = append(result, presentGame(&game, gameService))
	}
	return result
}

func presentGame(game *models.Game, gameService *service.GameService) gameResponse {
	return gameResponse{
		ID:           game.ID,
		Title:        game.Title,
		Description:  game.Description,
		FileURL:      game.FileURL,
		ThumbnailURL: gameService.BuildThumbnailURL(game.ThumbnailPath),
		DeveloperID:  game.DeveloperID,
		Status:       game.Status,
		CreatedAt:    game.CreatedAt.Format(time.RFC3339),
		Developer:    game.Developer,
		Difficulty:   game.Difficulty,
		Categories:   game.Categories,
	}
}
