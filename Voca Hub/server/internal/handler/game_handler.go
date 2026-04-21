package handler

import (
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/domain/dto"
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

	user := helper.MustCurrentUser(c)
	if gin.Mode() == gin.DebugMode {
		log.Printf("upload requested by user id=%d clerk_id=%s role=%s", user.ID, user.ClerkID, user.Role)
	}
	game, err := h.gameService.UploadGame(user.ID, req.Title, req.Description, file)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	helper.Success(c, http.StatusCreated, "game uploaded", game)
}

func (h *GameHandler) ListApprovedGames(c *gin.Context) {
	games, err := h.gameService.ListApprovedGames()
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "approved games fetched", games)
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
	helper.Success(c, http.StatusOK, "game fetched", game)
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
		c.AbortWithStatus(http.StatusBadRequest)
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
