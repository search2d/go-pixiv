package pixiv

import "net/http"

var DefaultDownloadHeaders = map[string]string{
	"User-Agent":      "PixivAndroidApp/5.0.64 (Android 6.0; Google Nexus 5X - 6.0.0 - API 23 - 1080x1920)",
	"Referer":         "https://app-api.pixiv.net/",
	"Accept-Encoding": "identity",
}

func SetDownloadHeaders(req *http.Request) {
	for k, v := range DefaultDownloadHeaders {
		req.Header.Set(k, v)
	}
}
