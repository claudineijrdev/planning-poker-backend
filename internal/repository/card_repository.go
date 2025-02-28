package repository

import (
	"fmt"
	"sync"

	"flash-cards/backend/internal/domain"

	"github.com/google/uuid"
)

type CardRepository struct {
	cards []domain.Card
	mutex sync.RWMutex
}

func NewCardRepository() *CardRepository {
	return &CardRepository{
		cards: make([]domain.Card, 0),
	}
}

func (r *CardRepository) GetAll() []domain.Card {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.cards
}

func (r *CardRepository) Create(card domain.Card) domain.Card {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	card.ID = uuid.New().String()
	r.cards = append(r.cards, card)
	return card
}

func (r *CardRepository) AddVote(cardID string, vote domain.Vote) (domain.Card, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i := range r.cards {
		if r.cards[i].ID == cardID {
			r.cards[i].Votes = append(r.cards[i].Votes, vote.Score)
			r.updateCardResults(&r.cards[i])
			return r.cards[i], nil
		}
	}
	return domain.Card{}, fmt.Errorf("card not found")
}

func (r *CardRepository) CloseVoting(cardID string) (domain.Card, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i := range r.cards {
		if r.cards[i].ID == cardID {
			r.cards[i].Closed = true
			return r.cards[i], nil
		}
	}
	return domain.Card{}, fmt.Errorf("card not found")
}

func (r *CardRepository) ResetAllVotes() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i := range r.cards {
		r.cards[i].Votes = []int{}
		r.cards[i].Result.Average = 0
		r.cards[i].Result.Distribution = make(map[int]int)
		r.cards[i].Closed = false
	}
}

func (r *CardRepository) updateCardResults(card *domain.Card) {
	sum := 0
	for _, v := range card.Votes {
		sum += v
	}
	card.Result.Average = float64(sum) / float64(len(card.Votes))
	card.Result.Distribution = make(map[int]int)
	for _, v := range card.Votes {
		card.Result.Distribution[v]++
	}
}
