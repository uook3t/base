package repo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"github.com/uook3t/base/logger"
	"github.com/uook3t/base/utils"
)

func (a *AlipanDriverImpl) GetAuthQrCode(ctx context.Context, scopeIDs []ScopeID) (*GetQrCodeResp, error) {
	httpReq := utils.NewRawHttpRequest(http.MethodPost, OpenUrl(PathAuthQrCode, nil), a.NewGetQrCodeReq(scopeIDs))
	httpResp, err := httpReq.DoSimpleHttp(ctx)
	if err != nil {
		return nil, err
	}
	var resp GetQrCodeResp
	err = sonic.Unmarshal(httpResp.Body, &resp)
	if err != nil {
		return nil, err
	}
	if len(resp.SID) > 0 {
		a.SetSid(resp.SID)
	}
	return &resp, nil
}

func (a *AlipanDriverImpl) WaitAuthQrCodeStatus(ctx context.Context) error {
	if len(a.sid) == 0 {
		return fmt.Errorf("qrcode sid is empty")
	}

	errChan := make(chan error)
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	defer close(errChan)

	go func() {
		for _ = range ticker.C {
			resp, err := a.getAuthQrCodeStatus(ctx)
			if err != nil {
				errChan <- err
				break
			}
			a.statusChan <- resp
			if IsStopQrCodeStatus(resp.Status) {
				close(a.statusChan)
				break
			}
		}
	}()

	timeout := time.After(3 * time.Minute)
	for {
		select {
		case res, ok := <-a.statusChan:
			if !ok {
				logger.Ctx(ctx).Infof("status chan is closed")
				return nil
			}
			logger.Ctx(ctx).Infof("recive status chan msg. status: %s, set to impl", res.Status)
			stop := a.SetAuthCodeResp(res)
			if stop {
				return nil
			}
		case err := <-errChan:
			return err
		case <-timeout:
			return fmt.Errorf("timeout for waiting rq status")
		}
	}
}

func (a *AlipanDriverImpl) getAuthQrCodeStatus(ctx context.Context) (*GetQrCodeStatusResp, error) {
	if len(a.sid) == 0 {
		return nil, fmt.Errorf("qrcode sid is empty")
	}
	httpReq := utils.NewRawHttpRequest(http.MethodGet, OpenUrl(fmt.Sprintf(PathAuthQrCodeStatus, a.sid), nil), nil)
	httpResp, err := httpReq.DoSimpleHttp(ctx)
	if err != nil {
		return nil, err
	}

	var resp GetQrCodeStatusResp
	err = sonic.Unmarshal(httpResp.Body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (a *AlipanDriverImpl) FetchAccessTokenByCode(ctx context.Context) error {
	if len(a.authCode) == 0 {
		return fmt.Errorf("auth code is empty")
	}
	httpReq := utils.NewRawHttpRequest(http.MethodPost, OpenUrl(PathAccessToken, nil), a.NewGetAccessTokenReq(a.authCode))
	httpResp, err := httpReq.DoSimpleHttp(ctx)
	if err != nil {
		return err
	}
	var resp GetAccessTokenResp
	err = sonic.Unmarshal(httpResp.Body, &resp)
	if err != nil {
		return err
	}
	if len(resp.AccessToken) == 0 {
		return fmt.Errorf("access_token is empty")
	}

	logger.Ctx(ctx).Infof("[FetchAccessTokenByCode] success. access_token: %s", resp.AccessToken)
	a.SetToken(resp.AccessToken, resp.RefreshToken)
	err = a.setTokenToCache(ctx, resp.AccessToken, resp.RefreshToken)
	if err != nil {
		logger.Ctx(ctx).WithError(err).Warnf("set refresh toke to cache failed")
	}
	return nil
}

func (a *AlipanDriverImpl) RefreshAccessToken(ctx context.Context) error {
	ok := a.refreshMutex.TryLock()
	if !ok {
		return fmt.Errorf("refresh failed as lock fail")
	}
	defer a.refreshMutex.Unlock()

	if len(a.refreshToken) == 0 {
		rt, err := a.getRefreshTokenFromCache(ctx)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return ErrRefreshTokenNotExist
			}
			logger.Ctx(ctx).WithError(err).Errorf("[RefreshAccessToken] get refresh_toke from cache failed")
			return err
		}
		a.refreshToken = rt
	}

	httpReq := utils.NewRawHttpRequest(http.MethodPost, OpenUrl(PathAccessToken, nil), a.NewRefreshAccessTokenReq(a.refreshToken))
	httpResp, err := httpReq.DoSimpleHttp(ctx)
	if err != nil {
		if httpResp.StatusCode == http.StatusUnauthorized {
			_ = a.delTokenCache(ctx)
		}
		return err
	}
	var resp GetAccessTokenResp
	err = sonic.Unmarshal(httpResp.Body, &resp)
	if err != nil {
		return err
	}
	if len(resp.AccessToken) == 0 {
		return fmt.Errorf("access_token is empty")
	}

	logger.Ctx(ctx).Infof("[RefreshAccessToken] success. new access_token: %s", resp.AccessToken)
	a.SetToken(resp.AccessToken, resp.RefreshToken)
	err = a.setTokenToCache(ctx, resp.AccessToken, resp.RefreshToken)
	if err != nil {
		logger.Ctx(ctx).WithError(err).Warnf("set refresh toke to cache failed")
	}
	return nil
}

func (a *AlipanDriverImpl) getAccessTokenFromCache(ctx context.Context) (string, error) {
	at, err := a.redisCli.Get(ctx, KeyAccessToken)
	if err != nil {
		return "", err
	}
	logger.Ctx(ctx).Infof("[getRefreshTokenFromCache] get access_token from cache success. val: %s", at)
	return at, nil
}

func (a *AlipanDriverImpl) getRefreshTokenFromCache(ctx context.Context) (string, error) {
	rt, err := a.redisCli.Get(ctx, KeyRefreshToken)
	if err != nil {
		return "", err
	}
	logger.Ctx(ctx).Infof("[getRefreshTokenFromCache] get refresh_toke from cache success. val: %s", rt)
	return rt, nil
}

func (a *AlipanDriverImpl) setTokenToCache(ctx context.Context, accessToken, refreshToken string) error {
	err := a.redisCli.MSet(ctx, map[string]string{
		KeyAccessToken:  accessToken,
		KeyRefreshToken: refreshToken,
	}, map[string]time.Duration{
		KeyAccessToken:  AccessTokenExpire,
		KeyRefreshToken: RefreshTokenExpire,
	})
	if err != nil {
		return err
	}
	logger.Ctx(ctx).Infof("[setTokenToCache] set tokens to cache success.")
	return nil
}

func (a *AlipanDriverImpl) delTokenCache(ctx context.Context) error {
	err := a.redisCli.Del(ctx, KeyAccessToken, KeyRefreshToken)
	if err != nil {
		return err
	}
	logger.Ctx(ctx).Debugf("[delTokenCache] get refresh_toke from cache success.")
	return nil
}
