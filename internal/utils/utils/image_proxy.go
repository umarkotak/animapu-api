package utils

import "fmt"

func AnimensionImgProxy(target string) string {
	host := "https://api.shadow-animapu-1.site"
	// host := "http://localhost:6001"
	return fmt.Sprintf(`%s/animes/animension/image_proxy/%s`, host, target)
}
