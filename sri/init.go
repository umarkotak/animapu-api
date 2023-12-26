package sri

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/config"
	"google.golang.org/api/option"
)

type Sri struct {
	FirebaseApp    *firebase.App
	FirebaseClient *db.Client
	DataGpsLog     *db.Ref
}

func New() (Sri, error) {
	var err error
	sriServer := Sri{}
	ctx := context.Background()

	// Init FB app
	firebaseConfig := &firebase.Config{
		DatabaseURL: config.Get().SriFirebaseUrl,
	}
	opts := []option.ClientOption{
		option.WithCredentialsJSON([]byte(config.Get().SriFirebaseGoogleServiceAccount)),
	}
	sriServer.FirebaseApp, err = firebase.NewApp(ctx, firebaseConfig, opts...)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	// Init FB client
	sriServer.FirebaseClient, err = sriServer.FirebaseApp.Database(ctx)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	// Init FB base reference
	sriServer.DataGpsLog = sriServer.FirebaseClient.NewRef("data_gps_log")
	if sriServer.DataGpsLog == nil {
		err = sriServer.DataGpsLog.Set(ctx, []interface{}{})
		if err != nil {
			logrus.WithContext(ctx).Error(err)
		}
	}

	return sriServer, nil
}
