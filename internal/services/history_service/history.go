package history_service

import (
	"context"
	"sort"

	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/history_repository"
	"github.com/umarkotak/animapu-api/internal/repository/user_repository"
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

func GetUserActivities(ctx context.Context, pagination models.Pagination) (contract.UserActivityData, error) {
	histories, err := history_repository.GetRecentActivities(ctx, pagination)
	if err != nil {
		return contract.UserActivityData{}, err
	}

	userIDs := make([]int64, 0, len(histories))
	for _, history := range histories {
		userIDs = append(userIDs, history.UserID)
	}

	users, err := user_repository.GetByIDs(ctx, userIDs)
	if err != nil {
		return contract.UserActivityData{}, err
	}

	result := contract.UserActivityData{Users: make([]contract.UserActivity, 0, len(users))}
	userIndex := make(map[int64]int, len(users))
	for _, user := range users {
		userIndex[user.ID] = len(result.Users)
		result.Users = append(result.Users, contract.UserActivity{
			VisitorID: user.VisitorId,
			Email:     user.Email.String,
			Histories: []contract.History{},
		})
	}

	for _, history := range histories {
		index, ok := userIndex[history.UserID]
		if !ok {
			continue
		}
		result.Users[index].Histories = append(result.Users[index].Histories, contract.History{
			MediaType: history.MediaType, Source: history.Source, SourceID: history.SourceID,
			Title: history.Title, CoverURLs: history.CoverURLs, LatestNumber: history.LatestNumber,
			Progress: history.Progress, LastLink: history.LastLink, UpdatedAt: history.UpdatedAt,
		})
	}

	sort.Slice(result.Users, func(i, j int) bool {
		return result.Users[i].Histories[0].UpdatedAt.After(result.Users[j].Histories[0].UpdatedAt)
	})
	return result, nil
}
