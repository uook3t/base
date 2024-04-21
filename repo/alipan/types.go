package repo

type DriveID string
type FileID string
type ScopeID string

const (
	ScopeUserBase   ScopeID = "user:base"
	ScopeFileRead   ScopeID = "file:all:read"
	ScopeFileWrite  ScopeID = "file:all:write"
	ScopeSharedRead ScopeID = "album:shared:read"
)

type GrantType string

const (
	GrantTypeAuthCode     GrantType = "authorization_code"
	GrantTypeRefreshToken GrantType = "refresh_token"
)

type QrCodeStatus string

const (
	QrCodeStatusWaiting      QrCodeStatus = "WaitLogin"
	QrCodeStatusScanSuccess  QrCodeStatus = "ScanSuccess"
	QrCodeStatusLoginSuccess QrCodeStatus = "LoginSuccess"
	QrCodeStatusExpired      QrCodeStatus = "QRCodeExpired"
)

func IsStopQrCodeStatus(c QrCodeStatus) bool {
	switch c {
	case QrCodeStatusLoginSuccess, QrCodeStatusExpired:
		return true
	default:
		return false
	}
}
