package utils

import (
	"fmt"

	"github.com/umarkotak/animapu-api/config"
)

func AnimensionImgProxy(target string) string {
	host := config.Get().AnimapuOnlineHost
	// host := "http://localhost:6001"
	return fmt.Sprintf(
		`%s/animes/animension/image_proxy/%s`, host, target,
	)
}
