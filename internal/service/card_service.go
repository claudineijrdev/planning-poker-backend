package service

import (
	"flash-cards/backend/internal/domain"
	"flash-cards/backend/internal/repository"
)

type CardService struct {
	repo *repository.CardRepository
}

func NewCardService(repo *repository.CardRepository) *CardService {
	return &CardService{
		repo: repo,
	}
}

func (s *CardService) GetAllCards() []domain.Card {
	return s.repo.GetAll()
}

func (s *CardService) CreateCard(card domain.Card) domain.Card {
	return s.repo.Create(card)
}

func (s *CardService) AddVote(cardID string, vote domain.Vote) (domain.Card, error) {
	return s.repo.AddVote(cardID, vote)
}

func (s *CardService) CloseVoting(cardID string) (domain.Card, error) {
	return s.repo.CloseVoting(cardID)
}

func (s *CardService) ResetAllVotes() {
	s.repo.ResetAllVotes()
}
