package repository

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/config"
	"github.com/umarkotak/animapu-api/internal/models"
	"google.golang.org/api/option"
)

var (
	firebaseApp       *firebase.App
	firebaseClient    *db.Client
	animapuLiteApiRef *db.Ref
	popularMangaRef   *db.Ref
)

func Initialize() {
	var err error
	ctx := context.Background()

	// Init FB app
	firebaseConfig := &firebase.Config{
		DatabaseURL: config.Get().AnimapuFirebaseUrl,
	}
	opts := []option.ClientOption{
		option.WithCredentialsJSON([]byte(config.Get().AnimapuGoogleServiceAccount)),
	}
	firebaseApp, err = firebase.NewApp(ctx, firebaseConfig, opts...)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	// Init FB client
	firebaseClient, err = firebaseApp.Database(ctx)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	// Init FB base reference
	animapuLiteApiRef = firebaseClient.NewRef("animapu-lite-api")
	if animapuLiteApiRef == nil {
		err = animapuLiteApiRef.Set(ctx, map[string]interface{}{
			"animapu-lite-api": map[string]interface{}{},
		})
		if err != nil {
			logrus.WithContext(ctx).Error(err)
		}
	}

	popularMangaRef = animapuLiteApiRef.Child("popular_manga")
	if popularMangaRef == nil {
		err = popularMangaRef.Set(ctx, map[string]interface{}{})
		if err != nil {
			logrus.WithContext(ctx).Error(err)
		}
	}

	for _, mangaSource := range models.MangaSources {
		mangaSourceRef := animapuLiteApiRef.Child(mangaSource.ID)
		if mangaSourceRef == nil {
			err = mangaSourceRef.Set(ctx, map[string]interface{}{
				mangaSource.ID: map[string]interface{}{},
				"available":    true,
			})
			if err != nil {
				logrus.WithContext(ctx).Error(err)
			}
		}
	}
}
