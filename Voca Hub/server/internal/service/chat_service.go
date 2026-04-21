package service

import (
	"errors"
	"time"

	"server/internal/domain/dto"
	"server/internal/domain/models"
	domainrepo "server/internal/domain/repository"
)

type ChatService struct {
	messageRepo domainrepo.MessageRepository
	userRepo    domainrepo.UserRepository
}

func NewChatService(messageRepo domainrepo.MessageRepository, userRepo domainrepo.UserRepository) *ChatService {
	return &ChatService{
		messageRepo: messageRepo,
		userRepo:    userRepo,
	}
}

func (s *ChatService) SaveMessage(senderID uint, receiverID uint, content string) (*models.Message, error) {
	if senderID == receiverID {
		return nil, errors.New("cannot send message to yourself")
	}
	receiver, err := s.userRepo.FindByID(receiverID)
	if err != nil {
		return nil, err
	}
	if receiver == nil {
		return nil, errors.New("receiver not found")
	}
	message := &models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
	}
	if err := s.messageRepo.Create(message); err != nil {
		return nil, err
	}
	return message, nil
}

func (s *ChatService) GetHistory(userID uint, otherUserID uint) ([]dto.ChatHistoryResponse, error) {
	other, err := s.userRepo.FindByID(otherUserID)
	if err != nil {
		return nil, err
	}
	if other == nil {
		return nil, errors.New("user not found")
	}
	messages, err := s.messageRepo.ListConversation(userID, otherUserID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.ChatHistoryResponse, 0, len(messages))
	for _, msg := range messages {
		result = append(result, dto.ChatHistoryResponse{
			ID:         msg.ID,
			SenderID:   msg.SenderID,
			ReceiverID: msg.ReceiverID,
			Content:    msg.Content,
			CreatedAt:  msg.CreatedAt.Format(time.RFC3339),
		})
	}
	return result, nil
}
