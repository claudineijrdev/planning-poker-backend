package domain

import (
	"time"
)

type SessionState string

const (
	SessionStateOpen   SessionState = "OPEN"
	SessionStateClosed SessionState = "CLOSED"
)

type UserRole string

const (
	UserRoleOwner UserRole = "OWNER"
	UserRoleGuest UserRole = "GUEST"
)

type Session struct {
	ID        string       `json:"id"`
	Code      string       `json:"code"`
	CreatedAt time.Time    `json:"createdAt"`
	State     SessionState `json:"state"`
	OwnerID   string       `json:"ownerId"`
	Cards     []Card       `json:"cards"`
	Users     []User       `json:"users"`
}

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Role      UserRole  `json:"role"`
	JoinedAt  time.Time `json:"joinedAt"`
	SessionID string    `json:"sessionId"`
}

type JoinSessionRequest struct {
	Code     string `json:"code"`
	UserName string `json:"userName"`
}

type CreateSessionRequest struct {
	OwnerName string `json:"ownerName"`
}

type CreateSessionResponse struct {
	Session Session `json:"session"`
	Code    string  `json:"code"`
}

type UpdateSessionStateRequest struct {
	State SessionState `json:"state"`
}

// Métodos auxiliares para Session
func (s *Session) IsOwner(userID string) bool {
	return s.OwnerID == userID
}

func (s *Session) GetUser(userID string) *User {
	for _, user := range s.Users {
		if user.ID == userID {
			return &user
		}
	}
	return nil
}

func (s *Session) AddUser(user User) {
	s.Users = append(s.Users, user)
}

func (s *Session) RemoveUser(userID string) {
	for i, user := range s.Users {
		if user.ID == userID {
			s.Users = append(s.Users[:i], s.Users[i+1:]...)
			return
		}
	}
}

func (s *Session) UpdateState(newState SessionState) {
	s.State = newState
}
