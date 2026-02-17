package service

import (
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/model"
	"github.com/ysmmc/backend/internal/repository"
)

type FavoriteService struct {
	favoriteRepo *repository.FavoriteRepository
	modelRepo    *repository.ModelRepository
}

func NewFavoriteService() *FavoriteService {
	return &FavoriteService{
		favoriteRepo: repository.NewFavoriteRepository(),
		modelRepo:    repository.NewModelRepository(),
	}
}

func (s *FavoriteService) Add(userID, modelID uuid.UUID) error {
	if s.favoriteRepo.Exists(userID, modelID) {
		return nil
	}

	favorite := &model.Favorite{
		UserID:  userID,
		ModelID: modelID,
	}

	return s.favoriteRepo.Create(favorite)
}

func (s *FavoriteService) Remove(userID, modelID uuid.UUID) error {
	return s.favoriteRepo.Delete(userID, modelID)
}

func (s *FavoriteService) IsFavorited(userID, modelID uuid.UUID) bool {
	return s.favoriteRepo.Exists(userID, modelID)
}

func (s *FavoriteService) ListByUserID(userID uuid.UUID, page, pageSize int) ([]model.Favorite, int64, error) {
	return s.favoriteRepo.ListByUserID(userID, page, pageSize)
}

func (s *FavoriteService) CountByModel(modelID uuid.UUID) (int64, error) {
	return s.favoriteRepo.CountByModel(modelID)
}
