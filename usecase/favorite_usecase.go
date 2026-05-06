package usecase

import (
	"context"
	"time"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
)

type favoriteUsecase struct {
	favoriteRepository domain.FavoriteRepository
	contextTimeout     time.Duration
}

func NewFavoriteUsecase(favoriteRepository domain.FavoriteRepository, timeout time.Duration) domain.FavoriteUsecase {
	return &favoriteUsecase{
		favoriteRepository: favoriteRepository,
		contextTimeout:     timeout,
	}
}

func (fu *favoriteUsecase) Create(c context.Context, favorite *domain.Favorite) error {
	ctx, cancel := context.WithTimeout(c, fu.contextTimeout)
	defer cancel()
	return fu.favoriteRepository.Create(ctx, favorite)
}

func (fu *favoriteUsecase) FetchByUserID(c context.Context, userID string) ([]domain.Favorite, error) {
	ctx, cancel := context.WithTimeout(c, fu.contextTimeout)
	defer cancel()
	return fu.favoriteRepository.FetchByUserID(ctx, userID)
}

func (fu *favoriteUsecase) Delete(c context.Context, userID string, placeID string) error {
	ctx, cancel := context.WithTimeout(c, fu.contextTimeout)
	defer cancel()
	return fu.favoriteRepository.Delete(ctx, userID, placeID)
}

func (fu *favoriteUsecase) IsFavorite(c context.Context, userID string, placeID string) (bool, error) {
	ctx, cancel := context.WithTimeout(c, fu.contextTimeout)
	defer cancel()
	return fu.favoriteRepository.IsFavorite(ctx, userID, placeID)
}
