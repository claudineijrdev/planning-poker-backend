package repository

import (
	"fmt"
	"sync"
	"time"

	"flash-cards/backend/internal/domain"
	"flash-cards/backend/internal/random"

	"github.com/google/uuid"
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
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{
		sessions:     make(map[string]domain.Session),
		sessionCodes: make(map[string]string),
		users:        make(map[string]domain.User),
	}
}

func (r *SessionRepository) CreateSession() (domain.Session, string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	code := r.generateUniqueCode()
	session := domain.Session{
		ID:        uuid.New().String(),
		Code:      code,
		CreatedAt: time.Now(),
		State:     domain.SessionStateOpen,
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

func (r *SessionRepository) AddUserToSession(code string, user domain.User) (domain.User, error) {
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

	user.ID = uuid.New().String()
	user.JoinedAt = time.Now()

	session.Users = append(session.Users, user)
	r.sessions[sessionID] = session
	r.users[user.ID] = user

	return user, nil
}

func (r *SessionRepository) UpdateSession(session domain.Session) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.sessions[session.ID]; !exists {
		return fmt.Errorf("session not found")
	}

	r.sessions[session.ID] = session
	return nil
}

func (r *SessionRepository) RemoveUserFromSession(sessionID string, userID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found")
	}

	for i, user := range session.Users {
		if user.ID == userID {
			session.Users = append(session.Users[:i], session.Users[i+1:]...)
			r.sessions[sessionID] = session
			delete(r.users, userID)
			return nil
		}
	}

	return fmt.Errorf("user not found in session")
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
		code[i] = codeChars[random.Intn(len(codeChars))]
	}
	return string(code)
}
