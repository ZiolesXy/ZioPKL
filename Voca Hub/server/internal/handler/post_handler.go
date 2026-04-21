package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"server/internal/domain/dto"
	"server/internal/helper"
	"server/internal/service"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) Create(c *gin.Context) {
	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user := helper.MustCurrentUser(c)
	post, err := h.postService.Create(user.ID, req.Content)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusCreated, "post created", post)
}

func (h *PostHandler) List(c *gin.Context) {
	posts, err := h.postService.ListAll()
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "posts fetched", helper.WrapListIfNeeded(posts))
}

func (h *PostHandler) ListMine(c *gin.Context) {
	user := helper.MustCurrentUser(c)
	posts, err := h.postService.ListMine(user.ID)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "user posts fetched", helper.WrapListIfNeeded(posts))
}

func (h *PostHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	post, err := h.postService.GetByID(uint(id))
	if err != nil {
		helper.Error(c, http.StatusNotFound, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "post fetched", post)
}

func (h *PostHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user := helper.MustCurrentUser(c)
	post, err := h.postService.Update(uint(id), user, req.Content)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		} else if err.Error() == "post not found" {
			status = http.StatusNotFound
		}
		helper.Error(c, status, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "post updated", post)
}

func (h *PostHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	user := helper.MustCurrentUser(c)
	if err := h.postService.Delete(uint(id), user); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		} else if err.Error() == "post not found" {
			status = http.StatusNotFound
		}
		helper.Error(c, status, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "post deleted", nil)
}
