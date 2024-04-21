package repo

import (
	"context"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/uook3t/base/utils"
)

func (a *AlipanDriverImpl) GetDownloadUrl(ctx context.Context, driveID DriveID, fileID FileID) (*GetDownloadURLResponse, error) {
	req := NewGetDownloadURLRequest(driveID, fileID)
	httpReq := utils.NewRawHttpRequest(http.MethodPost, OpenUrl(PathGetDownloadURL, nil), req)
	httpReq.WithAccessToken(a.accessToken)
	httpResp, err := httpReq.DoSimpleHttp(ctx)
	if err != nil {
		return nil, err
	}
	var resp GetDownloadURLResponse
	err = sonic.Unmarshal(httpResp.Body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (a *AlipanDriverImpl) DownloadLargeFile(ctx context.Context, url string) {

}
