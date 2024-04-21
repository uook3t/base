package repo

import (
	"context"
	"sync"

	"github.com/uook3t/base/logger"
	"github.com/uook3t/base/redis"
)

type AlipanDriver interface {
	GetAuthQrCode(ctx context.Context, scopeIDs []ScopeID) (*GetQrCodeResp, error)
	WaitAuthQrCodeStatus(ctx context.Context) error
	FetchAccessTokenByCode(ctx context.Context) error
	RefreshAccessToken(ctx context.Context) error
	GetDownloadUrl(ctx context.Context, driveID DriveID, fileID FileID) (*GetDownloadURLResponse, error)
}

type AlipanDriverImpl struct {
	codeStatus   QrCodeStatus
	sid          string
	authCode     string
	accessToken  string
	refreshToken string

	config       *OpenConfig
	redisCli     *redis.RedisCli
	statusChan   chan *GetQrCodeStatusResp
	refreshMutex sync.Mutex
}

func NewAlipanDriverImpl(redisCli *redis.RedisCli) *AlipanDriverImpl {
	return &AlipanDriverImpl{
		codeStatus: QrCodeStatusWaiting,
		statusChan: make(chan *GetQrCodeStatusResp),
		redisCli:   redisCli,
	}
}

func (a *AlipanDriverImpl) WithConfig(c *OpenConfig) *AlipanDriverImpl {
	a.config = c
	return a
}

func (a *AlipanDriverImpl) SetSid(sid string) {
	a.sid = sid
}

func (a *AlipanDriverImpl) SetToken(accessToken, refreshToken string) {
	a.accessToken = accessToken
	a.refreshToken = refreshToken
}

func (a *AlipanDriverImpl) SetAuthCodeResp(res *GetQrCodeStatusResp) bool {
	a.codeStatus = res.Status
	if len(res.AuthCode) != 0 {
		a.authCode = res.AuthCode
	}
	return len(res.AuthCode) != 0 || IsStopQrCodeStatus(res.Status)
}

func (a *AlipanDriverImpl) GetAuthCode() string {
	return a.authCode
}

func (a *AlipanDriverImpl) GetAccessToken() string {
	return a.accessToken
}

func (a *AlipanDriverImpl) TryWithCacheAccessToken(ctx context.Context) bool {
	at, err := a.getAccessTokenFromCache(ctx)
	if err != nil {
		logger.Ctx(ctx).WithError(err).Warnf("get access_token from cache failed")
		return false
	}
	a.accessToken = at
	return len(at) > 0
}

var _ AlipanDriver = &AlipanDriverImpl{}
