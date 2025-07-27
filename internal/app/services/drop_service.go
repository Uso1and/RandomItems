package services

import (
	"context"
	"log"
	"math/rand/v2"
	"sort"
	"time"

	"RandomItems/internal/domain/models"
	"RandomItems/internal/domain/repositories"
)

type DropService struct {
	itemRepo repositories.ItemRepositoryInterface
	dropRepo repositories.DropRepositoryInterface
	userRepo repositories.UserRepositoryInterface
}

func NewDropService(
	itemRepo repositories.ItemRepositoryInterface,
	dropRepo repositories.DropRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
) *DropService {
	return &DropService{
		itemRepo: itemRepo,
		dropRepo: dropRepo,
		userRepo: userRepo,
	}
}

func (s *DropService) GenerateDrop(c context.Context, userID int) (*models.Item, error) {
	user, err := s.userRepo.GetUser(c, userID)
	if err != nil {
		return nil, err
	}

	items, err := s.itemRepo.GetItems(c)
	if err != nil {
		return nil, err
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].MinPity > items[j].MinPity
	})

	var guaranteedItem *models.Item
	for _, item := range items {
		if item.MinPity > 0 && user.PityCounter >= item.MinPity {
			guaranteedItem = item
			break
		}
	}

	if guaranteedItem != nil {
		log.Printf("GUARANTEED DROP - UserID: %d, ItemID: %d", userID, guaranteedItem.ID)
		if err := s.dropRepo.UpdateUserPityCounter(c, userID, 0); err != nil {
			return nil, err
		}
		dropEvent := &models.DropEvent{
			UserID:       userID,
			ItemID:       guaranteedItem.ID,
			DroppedAt:    time.Now(),
			IsGuaranteed: true,
		}
		if err := s.dropRepo.CreateDropEvent(c, dropEvent); err != nil {
			return nil, err
		}
		return guaranteedItem, nil
	}

	baseSuccessRate := 0.3
	pityBonus := float64(user.PityCounter) * 0.03
	successThreshold := baseSuccessRate + pityBonus
	if successThreshold > 0.8 {
		successThreshold = 0.8
	}

	r := rand.Float64()
	if r > successThreshold {

		newPity := user.PityCounter + 1
		if newPity > 50 {
			newPity = 50
		}
		if err := s.dropRepo.UpdateUserPityCounter(c, userID, newPity); err != nil {
			return nil, err
		}
		return nil, nil
	}

	var availableItems []*models.Item
	totalWeight := 0.0

	for _, item := range items {
		weight := item.BaseChance

		if item.MinPity > 0 {
			pityProgress := float64(user.PityCounter) / float64(item.MinPity)
			if pityProgress > 0.3 {
				weight = item.BaseChance * (1 + pityProgress*2)
			}
		}

		if weight > 0 {
			availableItems = append(availableItems, &models.Item{
				ID:         item.ID,
				Name:       item.Name,
				Rarity:     item.Rarity,
				BaseChance: weight,
				MinPity:    item.MinPity,
			})
			totalWeight += weight
		}
	}

	if len(availableItems) == 0 {
		return nil, nil
	}

	r = rand.Float64() * totalWeight
	cumulative := 0.0
	var selectedItem *models.Item

	for _, item := range availableItems {
		cumulative += item.BaseChance
		if r <= cumulative {
			selectedItem = item
			break
		}
	}

	if selectedItem != nil {
		if err := s.dropRepo.UpdateUserPityCounter(c, userID, 0); err != nil {
			return nil, err
		}
		dropEvent := &models.DropEvent{
			UserID:       userID,
			ItemID:       selectedItem.ID,
			DroppedAt:    time.Now(),
			IsGuaranteed: false,
		}
		if err := s.dropRepo.CreateDropEvent(c, dropEvent); err != nil {
			return nil, err
		}
		return selectedItem, nil
	}

	return nil, nil
}
