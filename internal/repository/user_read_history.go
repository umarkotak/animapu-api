package repository

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func RecordUserReadHistory(ctx context.Context, user models.User, manga models.Manga) error {
	var err error

	oneUser := usersRef.Child(user.Uid)

	if oneUser == nil {
		user.ReadHistories = []models.Manga{manga}
		user.ReadHistoriesMap = map[string]models.Manga{manga.GetFbUniqueKey(): manga}

		err = oneUser.Set(ctx, user)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return err
		}

		return nil
	}

	err = oneUser.Get(ctx, &user)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	err = oneUser.Child("read_histories_map").Child(manga.GetFbUniqueKey()).Set(ctx, manga)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	tempReadHistories := []models.Manga{}
	for _, historyManga := range user.ReadHistories {
		if historyManga.GetFbUniqueKey() != manga.GetFbUniqueKey() {
			tempReadHistories = append(tempReadHistories, historyManga)
		}
	}
	user.ReadHistories = tempReadHistories

	user.ReadHistories = append([]models.Manga{manga}, user.ReadHistories...)
	err = oneUser.Child("read_histories").Set(ctx, user.ReadHistories)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	return nil
}

func GetUserReadHistories(ctx context.Context, user models.User) ([]models.Manga, map[string]models.Manga, error) {
	mangaHistories := []models.Manga{}
	mangaHistoriesMap := map[string]models.Manga{}

	oneUser := usersRef.Child(user.Uid)
	if oneUser == nil {
		return mangaHistories, mangaHistoriesMap, nil
	}

	readHistoriesRef := oneUser.Child("read_histories")
	if readHistoriesRef == nil {
		return mangaHistories, mangaHistoriesMap, nil
	}

	err := readHistoriesRef.Get(ctx, &mangaHistories)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangaHistories, mangaHistoriesMap, err
	}

	// readHistoriesMapRef := oneUser.Child("read_histories_map")
	// if readHistoriesMapRef == nil {
	// 	return mangaHistories, mangaHistoriesMap, nil
	// }

	// err = readHistoriesMapRef.Get(ctx, &mangaHistoriesMap)
	// if err != nil {
	// 	logrus.WithContext(ctx).Error(err)
	// 	return mangaHistories, mangaHistoriesMap, err
	// }

	return mangaHistories, mangaHistoriesMap, nil
}
