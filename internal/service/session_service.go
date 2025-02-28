package service

import (
	"flash-cards/backend/internal/domain"
	"flash-cards/backend/internal/repository"
)

type SessionService struct {
	sessionRepo *repository.SessionRepository
	cardRepo    *repository.CardRepository
}

func NewSessionService(sessionRepo *repository.SessionRepository, cardRepo *repository.CardRepository) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
		cardRepo:    cardRepo,
	}
}

func (s *SessionService) CreateSession() domain.CreateSessionResponse {
	session, code := s.sessionRepo.CreateSession()
	return domain.CreateSessionResponse{
		Session: session,
		Code:    code,
	}
}

func (s *SessionService) JoinSession(req domain.JoinSessionRequest) (domain.User, error) {
	return s.sessionRepo.AddUserToSession(req.Code, req.UserName)
}

func (s *SessionService) GetSessionByCode(code string) (domain.Session, error) {
	return s.sessionRepo.GetSessionByCode(code)
}

func (s *SessionService) CreateCardInSession(sessionID string, card domain.Card) (domain.Card, error) {
	card = s.cardRepo.Create(card)
	err := s.sessionRepo.AddCardToSession(sessionID, card)
	if err != nil {
		return domain.Card{}, err
	}
	return card, nil
}

func (s *SessionService) GetSession(sessionID string) (domain.Session, error) {
	return s.sessionRepo.GetSession(sessionID)
}
