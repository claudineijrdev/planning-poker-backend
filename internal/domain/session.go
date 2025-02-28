package domain

import (
	"time"
)

type Session struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"createdAt"`
	Cards     []Card    `json:"cards"`
	Users     []User    `json:"users"`
}

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	JoinedAt  time.Time `json:"joinedAt"`
	SessionID string    `json:"sessionId"`
}

type JoinSessionRequest struct {
	Code     string `json:"code"`
	UserName string `json:"userName"`
}

type CreateSessionResponse struct {
	Session Session `json:"session"`
	Code    string  `json:"code"`
}
