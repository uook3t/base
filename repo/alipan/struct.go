package repo

type GetQrCodeReq struct {
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	Scopes       []ScopeID `json:"scopes"`
}

func (a *AlipanDriverImpl) NewGetQrCodeReq(scopes []ScopeID) *GetQrCodeReq {
	return &GetQrCodeReq{
		ClientID:     a.config.ClientID,
		ClientSecret: a.config.ClientSecret,
		Scopes:       scopes,
	}
}

type GetQrCodeResp struct {
	QrCodeUrl string `json:"qrCodeUrl"`
	SID       string `json:"sid"`
}

type GetQrCodeStatusResp struct {
	Status   QrCodeStatus `json:"status"`
	AuthCode string       `json:"authCode"`
}

type GetAccessTokenReq struct {
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	GrantType    GrantType `json:"grant_type"`
	Code         string    `json:"code"`
	RefreshToken string    `json:"refresh_token"`
}

func (a *AlipanDriverImpl) NewGetAccessTokenReq(code string) *GetAccessTokenReq {
	return &GetAccessTokenReq{
		ClientID:     a.config.ClientID,
		ClientSecret: a.config.ClientSecret,
		GrantType:    GrantTypeAuthCode,
		Code:         code,
	}
}

func (a *AlipanDriverImpl) NewRefreshAccessTokenReq(refreshToken string) *GetAccessTokenReq {
	return &GetAccessTokenReq{
		ClientID:     a.config.ClientID,
		ClientSecret: a.config.ClientSecret,
		GrantType:    GrantTypeRefreshToken,
		RefreshToken: refreshToken,
	}
}

type GetAccessTokenResp struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type GetDownloadURLRequest struct {
	DriveID DriveID `json:"drive_id"`
	FileID  FileID  `json:"file_id"`
}

func NewGetDownloadURLRequest(driveID DriveID, fileID FileID) *GetDownloadURLRequest {
	return &GetDownloadURLRequest{
		DriveID: driveID,
		FileID:  fileID,
	}
}

type GetDownloadURLResponse struct {
	URL        string `json:"url"`
	Expiration string `json:"expiration"`
	Method     string `json:"method"`
}
