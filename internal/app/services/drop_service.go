package services

import (
	"RandomItems/internal/domain/models"
	"RandomItems/internal/domain/repositories"
	"context"
	"math/rand/v2"
	"sort"
	"time"
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

	// 1. Проверка гарантированного дропа (выбираем предмет с самым высоким min_pity)
	var guaranteedItem *models.Item
	currentPity := user.PityCounter

	// Сортируем предметы по min_pity в порядке убывания
	sort.Slice(items, func(i, j int) bool {
		return items[i].MinPity > items[j].MinPity
	})

	for _, item := range items {
		if item.MinPity > 0 && currentPity >= item.MinPity {
			guaranteedItem = item
			break
		}
	}

	if guaranteedItem != nil {
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

	// 2. Рассчитываем шанс успешного дропа с учётом pity_counter
	baseSuccessRate := 0.7
	pityBonus := float64(user.PityCounter) * 0.02
	successThreshold := baseSuccessRate + pityBonus
	if successThreshold > 0.95 {
		successThreshold = 0.95
	}

	// 3. Проверяем, выпал ли предмет в этой попытке
	r := rand.Float64()
	if r > successThreshold {
		// Неудача - увеличиваем pity_counter
		newPity := user.PityCounter + 1
		if err := s.dropRepo.UpdateUserPityCounter(c, userID, newPity); err != nil {
			return nil, err
		}
		return nil, nil
	}

	// 4. Выбираем предмет из доступных (min_pity == 0)
	var availableItems []*models.Item
	for _, item := range items {
		if item.MinPity == 0 {
			availableItems = append(availableItems, item)
		}
	}

	if len(availableItems) == 0 {
		// Если нет доступных предметов, увеличиваем pity_counter
		newPity := user.PityCounter + 1
		if err := s.dropRepo.UpdateUserPityCounter(c, userID, newPity); err != nil {
			return nil, err
		}
		return nil, nil
	}

	// 5. Выбираем случайный предмет с учётом их базовых шансов
	totalWeight := 0.0
	for _, item := range availableItems {
		totalWeight += item.BaseChance
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

	// 6. Обновляем pity_counter и записываем дроп
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

	// На всякий случай, если что-то пошло не так
	newPity := user.PityCounter + 1
	if err := s.dropRepo.UpdateUserPityCounter(c, userID, newPity); err != nil {
		return nil, err
	}
	return nil, nil
}
