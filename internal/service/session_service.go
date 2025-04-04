package service

import (
	"errors"
	"flash-cards/backend/internal/domain"
	"flash-cards/backend/internal/repository"
	"time"
)

var (
	ErrSessionNotFound = errors.New("sessão não encontrada")
	ErrUnauthorized    = errors.New("usuário não autorizado")
	ErrSessionClosed   = errors.New("sessão está fechada")
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

func (s *SessionService) CreateSession(req domain.CreateSessionRequest) (domain.CreateSessionResponse, error) {
	session, code := s.sessionRepo.CreateSession()

	// Criar o usuário owner
	owner := domain.User{
		Name:      req.OwnerName,
		Role:      domain.UserRoleOwner,
		SessionID: session.ID,
		JoinedAt:  time.Now(),
	}

	// Adicionar o owner à sessão através do repositório para gerar o ID
	owner, err := s.sessionRepo.AddUserToSession(code, owner)
	if err != nil {
		return domain.CreateSessionResponse{}, err
	}

	// Atualizar o ownerID da sessão
	session.OwnerID = owner.ID
	session.State = domain.SessionStateOpen
	session.Users = append(session.Users, owner)

	// Atualizar a sessão no repositório
	err = s.sessionRepo.UpdateSession(session)
	if err != nil {
		return domain.CreateSessionResponse{}, err
	}

	return domain.CreateSessionResponse{
		Session: session,
		Code:    code,
	}, nil
}

func (s *SessionService) JoinSession(code string, req domain.JoinSessionRequest) (domain.User, error) {
	session, err := s.sessionRepo.GetSessionByCode(code)
	if err != nil {
		return domain.User{}, ErrSessionNotFound
	}

	if session.State == domain.SessionStateClosed {
		return domain.User{}, ErrSessionClosed
	}

	// Criar novo usuário como convidado
	user := domain.User{
		Name:      req.UserName,
		Role:      domain.UserRoleGuest,
		SessionID: session.ID,
	}

	return s.sessionRepo.AddUserToSession(code, user)
}

func (s *SessionService) UpdateSessionState(code string, userID string, req domain.UpdateSessionStateRequest) error {
	session, err := s.sessionRepo.GetSessionByCode(code)
	if err != nil {
		return ErrSessionNotFound
	}

	if !session.IsOwner(userID) {
		return ErrUnauthorized
	}

	session.UpdateState(req.State)
	return s.sessionRepo.UpdateSession(session)
}

func (s *SessionService) GetSessionByCode(code string) (domain.Session, error) {
	return s.sessionRepo.GetSessionByCode(code)
}

func (s *SessionService) CreateCardInSession(code string, userID string, card domain.Card) (domain.Card, error) {
	session, err := s.sessionRepo.GetSessionByCode(code)
	if err != nil {
		return domain.Card{}, ErrSessionNotFound
	}

	if !session.IsOwner(userID) {
		return domain.Card{}, ErrUnauthorized
	}

	if session.State == domain.SessionStateClosed {
		return domain.Card{}, ErrSessionClosed
	}

	card = s.cardRepo.Create(card)
	err = s.sessionRepo.AddCardToSession(session.ID, card)
	if err != nil {
		return domain.Card{}, err
	}
	return card, nil
}

func (s *SessionService) GetSession(sessionID string) (domain.Session, error) {
	return s.sessionRepo.GetSession(sessionID)
}

func (s *SessionService) LeaveSession(code string, userID string) error {
	session, err := s.sessionRepo.GetSessionByCode(code)
	if err != nil {
		return ErrSessionNotFound
	}

	if session.IsOwner(userID) {
		session.State = domain.SessionStateClosed
	}

	session.RemoveUser(userID)
	return s.sessionRepo.UpdateSession(session)
}

func (s *SessionService) ResetSessionVotes(sessionCode string) ([]domain.Card, error) {
	session, err := s.sessionRepo.GetSessionByCode(sessionCode)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	// Resetar votos de todos os cards da sessão
	for i := range session.Cards {
		session.Cards[i].Votes = []int{}
		session.Cards[i].Result.Average = 0
		session.Cards[i].Result.Distribution = make(map[int]int)
		session.Cards[i].Closed = false
	}

	// Atualizar a sessão no repositório
	err = s.sessionRepo.UpdateSession(session)
	if err != nil {
		return nil, err
	}

	return session.Cards, nil
}
