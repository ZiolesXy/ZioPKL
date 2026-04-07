package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"voca-store/internal/domain/models"
	"voca-store/internal/domain/repository"
	"voca-store/internal/websocket"

	"github.com/google/uuid"
)

type ChatService interface {
	CreateChatRequest(ctx context.Context, userID uint, message *string) (*models.ChatSession, error)
	AcceptChatRequest(ctx context.Context, sessionUID string, adminID uint) (*models.ChatSession, error)
	GetPendingChatRequests(ctx context.Context) ([]models.ChatSession, error)
	GetUserActiveSession(ctx context.Context, userID uint) (*models.ChatSession, error)
	GetSessionByUID(ctx context.Context, uid string, requesterID uint, role string) (*models.ChatSession, error)
	GetChatHistory(ctx context.Context, sessionID uint, limit, offset int) ([]models.ChatMessage, error)
	SendMessage(ctx context.Context, sessionUID, content, msgType string, mediaURL *string, senderID uint) (*models.ChatMessage, error)
	MarkMessagesRead(ctx context.Context, sessionUID string, userID uint) error
	CloseChatSession(ctx context.Context, sessionUID string, closerID uint) error
}

type chatService struct {
	repo repository.ChatRepository
	hub  *websocket.Hub
}

func NewChatService(repo repository.ChatRepository, hub *websocket.Hub) ChatService {
	return &chatService{repo: repo, hub: hub}
}

func (s *chatService) CreateChatRequest(ctx context.Context, userID uint, message *string) (*models.ChatSession, error) {
	existing, _ := s.repo.GetUserActiveSession(ctx, userID)
	if existing != nil && existing.Status == models.ChatSessionActive {
		return nil, errors.New("you already have an active chat session")
	}

	session := &models.ChatSession{
		UID:     uuid.New().String(),
		UserID:  userID,
		Status:  models.ChatSessionActive,
	}

	if message != nil && *message != "" {
		session.LastMessage = message
		now := time.Now()
		session.LastMessageAt = &now
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create chat session: %w", err)
	}

	go s.hub.NotifyAdminsNewRequest(*session)

	return session, nil
}

func (s *chatService) AcceptChatRequest(ctx context.Context, sessionUID string, adminID uint) (*models.ChatSession, error) {
	session, err := s.repo.GetSessionByUID(ctx, sessionUID)
	if err != nil {
		return nil, errors.New("chat session not found")
	}

	if session.Status != models.ChatSessionPending && session.Status != models.ChatSessionActive {
		return nil, fmt.Errorf("cannot accept session with status: %s", session.Status)
	}

	if session.AdminID != nil {
		return nil, errors.New("chat session already has an assigned admin")
	}

	if err := s.repo.UpdateSessionStatus(ctx, session.ID, models.ChatSessionActive, &adminID); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	session, err = s.repo.GetSessionByID(ctx, session.ID)
	if err != nil {
		return nil, err
	}

	go s.hub.NotifySessionAccepted(session.UID, session.Admin.Name)

	return session, nil
}

func (s *chatService) GetPendingChatRequests(ctx context.Context) ([]models.ChatSession, error) {
	return s.repo.GetPendingSessions(ctx)
}

func (s *chatService) GetUserActiveSession(ctx context.Context, userID uint) (*models.ChatSession, error) {
	return s.repo.GetUserActiveSession(ctx, userID)
}

func (s *chatService) GetSessionByUID(ctx context.Context, uid string, requesterID uint, role string) (*models.ChatSession, error) {
	session, err := s.repo.GetSessionByUID(ctx, uid)
	if err != nil {
		return nil, errors.New("chat session not found")
	}

	if session.UserID != requesterID && role != "Admin" {
		if session.AdminID == nil || *session.AdminID != requesterID {
			return nil, errors.New("access denied")
		}
	}

	return session, nil
}

func (s *chatService) GetChatHistory(ctx context.Context, sessionID uint, limit, offset int) ([]models.ChatMessage, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.GetMessagesBySession(ctx, sessionID, limit, offset)
}

func (s *chatService) SendMessage(ctx context.Context, sessionUID, content, msgType string, mediaURL *string, senderID uint) (*models.ChatMessage, error) {
	session, err := s.repo.GetSessionByUID(ctx, sessionUID)
	if err != nil {
		return nil, errors.New("chat session not found")
	}

	if session.Status != models.ChatSessionActive {
		return nil, errors.New("cannot send message to inactive session")
	}

	if session.UserID != senderID && (session.AdminID == nil || *session.AdminID != senderID) {
		return nil, errors.New("access denied")
	}

	chatMsg := &models.ChatMessage{
		SessionID:   session.ID,
		SenderID:    senderID,
		Content:     content,
		MessageType: models.MessageType(msgType),
	}
	if mediaURL != nil && *mediaURL != "" {
		chatMsg.MediaURL = mediaURL
	}

	if err := s.repo.CreateMessage(ctx, chatMsg); err != nil {
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	if err := s.repo.UpdateLastMessage(ctx, session.ID, content); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to update last message: %v\n", err)
	}

	return chatMsg, nil
}

func (s *chatService) MarkMessagesRead(ctx context.Context, sessionUID string, userID uint) error {
	session, err := s.repo.GetSessionByUID(ctx, sessionUID)
	if err != nil {
		return errors.New("chat session not found")
	}

	if session.UserID != userID && (session.AdminID == nil || *session.AdminID != userID) {
		return errors.New("access denied")
	}

	return s.repo.MarkMessagesAsRead(ctx, session.ID, userID)
}

func (s *chatService) CloseChatSession(ctx context.Context, sessionUID string, closerID uint) error {
	session, err := s.repo.GetSessionByUID(ctx, sessionUID)
	if err != nil {
		return errors.New("chat session not found")
	}

	if session.UserID != closerID && (session.AdminID == nil || *session.AdminID != closerID) {
		return errors.New("access denied")
	}

	if err := s.repo.CloseSession(ctx, session.ID); err != nil {
		return fmt.Errorf("failed to close session: %w", err)
	}

	return nil
}