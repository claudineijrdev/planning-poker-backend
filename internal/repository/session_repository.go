package repository

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"flash-cards/backend/internal/domain"
)

const (
	codeLength = 6
	codeChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type SessionRepository struct {
	sessions     map[string]domain.Session // ID -> Session
	sessionCodes map[string]string         // Code -> ID
	users        map[string]domain.User    // UserID -> User
	mutex        sync.RWMutex
	nextUserID   int
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{
		sessions:     make(map[string]domain.Session),
		sessionCodes: make(map[string]string),
		users:        make(map[string]domain.User),
		nextUserID:   0,
	}
}

func (r *SessionRepository) CreateSession() (domain.Session, string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	code := r.generateUniqueCode()
	session := domain.Session{
		ID:        fmt.Sprintf("session_%d", len(r.sessions)),
		Code:      code,
		CreatedAt: time.Now(),
		Cards:     make([]domain.Card, 0),
		Users:     make([]domain.User, 0),
	}

	r.sessions[session.ID] = session
	r.sessionCodes[code] = session.ID

	return session, code
}

func (r *SessionRepository) GetSessionByCode(code string) (domain.Session, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	sessionID, exists := r.sessionCodes[code]
	if !exists {
		return domain.Session{}, fmt.Errorf("session not found")
	}

	session, exists := r.sessions[sessionID]
	if !exists {
		return domain.Session{}, fmt.Errorf("session not found")
	}

	return session, nil
}

func (r *SessionRepository) AddUserToSession(code string, userName string) (domain.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	sessionID, exists := r.sessionCodes[code]
	if !exists {
		return domain.User{}, fmt.Errorf("session not found")
	}

	session, exists := r.sessions[sessionID]
	if !exists {
		return domain.User{}, fmt.Errorf("session not found")
	}

	user := domain.User{
		ID:        fmt.Sprintf("user_%d", r.nextUserID),
		Name:      userName,
		JoinedAt:  time.Now(),
		SessionID: session.ID,
	}
	r.nextUserID++

	session.Users = append(session.Users, user)
	r.sessions[sessionID] = session
	r.users[user.ID] = user

	return user, nil
}

func (r *SessionRepository) GetSession(sessionID string) (domain.Session, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return domain.Session{}, fmt.Errorf("session not found")
	}

	return session, nil
}

func (r *SessionRepository) AddCardToSession(sessionID string, card domain.Card) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found")
	}

	session.Cards = append(session.Cards, card)
	r.sessions[sessionID] = session

	return nil
}

func (r *SessionRepository) generateUniqueCode() string {
	for {
		code := r.generateCode()
		if _, exists := r.sessionCodes[code]; !exists {
			return code
		}
	}
}

func (r *SessionRepository) generateCode() string {
	code := make([]byte, codeLength)
	for i := range code {
		code[i] = codeChars[rand.Intn(len(codeChars))]
	}
	return string(code)
}
