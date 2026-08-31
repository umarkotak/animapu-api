package history_service

import (
	"context"

	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/history_repository"
)

func GetHistories(ctx context.Context, user models.User, pagination models.Pagination) ([]contract.History, error) {
	histories, err := history_repository.GetByUserID(ctx, user.ID, pagination)
	if err != nil {
		return []contract.History{}, err
	}

	result := make([]contract.History, 0, len(histories))
	for _, history := range histories {
		result = append(result, contract.History{
			MediaType:    history.MediaType,
			Source:       history.Source,
			SourceID:     history.SourceID,
			Title:        history.Title,
			CoverURLs:    history.CoverURLs,
			LatestNumber: history.LatestNumber,
			Progress:     history.Progress,
			LastLink:     history.LastLink,
			UpdatedAt:    history.UpdatedAt,
		})
	}

	return result, nil
}
