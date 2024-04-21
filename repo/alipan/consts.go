package repo

import (
	"errors"
	"net/url"
	"time"
)

const (
	AlinpanOpenDomain = "openapi.alipan.com"

	PathAccessToken      = "/oauth/access_token"
	PathAuthQrCode       = "/oauth/authorize/qrcode"
	PathAuthQrCodeStatus = "/oauth/qrcode/%s/status"

	PathGetDownloadURL = "/adrive/v1.0/openFile/getDownloadUrl"

	KeyAccessToken  = "drivebox:alipanopen:access_token"
	KeyRefreshToken = "drivebox:alipanopen:refresh_token"

	AccessTokenExpire  = time.Hour * 1
	RefreshTokenExpire = time.Hour * 24 * 90
)

var (
	ErrRefreshTokenNotExist = errors.New("refresh token not exist")
)

func OpenUrl(path string, queries map[string]string) string {
	u := url.URL{
		Scheme: "https",
		Host:   AlinpanOpenDomain,
		Path:   path,
	}

	values := url.Values{}
	for k, v := range queries {
		values.Add(k, v)
	}
	u.RawQuery = values.Encode()

	// u.Fragment = "example-ParseQuery"
	return u.String()
}
