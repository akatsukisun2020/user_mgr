package login

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/akatsukisun2020/go_components/logger"
	pb "github.com/akatsukisun2020/proto_proj/user_mgr"
	svrConf "github.com/akatsukisun2020/user_mgr/config"
)

// 微信后台登录: https://docs.qq.com/doc/DV1JSVm5kSHNDS3lH
type wxBackendLogin struct {
}

// 参考文档： https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-login/code2Session.html
type jscode2sessionReq struct {
	appid     string
	secret    string
	jsCode    string
	grantType string
}

type jscode2sessionRsp struct {
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	Openid     string `json:"openid"`
	Errmsg     string `json:"errmsg"`
	Errcode    int32  `json:"errcode"`
}

func (wbl *wxBackendLogin) Login(ctx context.Context, credential *pb.LoginCredential) (*pb.LoginInfo, error) {
	if credential.GetLoginType() != pb.ELoginType_LoginWx {
		logger.ErrorContextf(ctx, "login type error, not ELoginType_LoginWx")
		return nil, fmt.Errorf("login type error, not ELoginType_LoginWx")
	}

	req := &jscode2sessionReq{
		appid:     svrConf.GetUserConfig().WXLoginAppID,
		secret:    svrConf.GetUserConfig().WXLoginSecret,
		jsCode:    credential.GetCode(),
		grantType: "authorization_code",
	}
	rsp := &jscode2sessionRsp{}
	if err := jscode2session(ctx, req, rsp); err != nil {
		logger.ErrorContextf(ctx, "jscode2session error, jscode:%s, err:%v", req.jsCode, err)
		return nil, err
	}

	return &pb.LoginInfo{
		LoginType:   pb.ELoginType_LoginWx,
		OpenId:      rsp.Openid,
		SessionKey:  rsp.SessionKey,
		AccessToken: EncodeToAccessToken(),
	}, nil
}

func jscode2session(ctx context.Context, req *jscode2sessionReq, rsp *jscode2sessionRsp) error {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=%s",
		req.appid, req.secret, req.jsCode, req.grantType)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.ErrorContextf(ctx, "jscode2session.NewRequest error, url:%s, err:%v", url, err)
		return err
	}
	client := http.Client{
		Timeout:   2 * time.Second,       // 超时时间控制
		Transport: http.DefaultTransport, // 默认
	}
	response, err := client.Do(request)
	if err != nil {
		logger.ErrorContextf(ctx, "jscode2session.Do error, url:%s, err:%v", url, err)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		logger.ErrorContextf(ctx, "jscode2session.Do error, url:%s, response:%v", url, response)
		return fmt.Errorf("jscode2session error, StatusCode is 200")
	}
	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.ErrorContextf(ctx, "jscode2session.ReadAll error, url:%s, err:%v", url, err)
		return err
	}
	err = json.Unmarshal(respBytes, rsp)
	if err != nil {
		logger.ErrorContextf(ctx, "jscode2session.Unmarshal error, url:%s, err:%v", url, err)
		return err
	}

	if rsp.Errcode != 0 {
		logger.ErrorContextf(ctx, "jscode2session's response error, url:%s, response:%v", url, response)
		return fmt.Errorf("jscode2session error, code:%d, msg:%s", rsp.Errcode, rsp.Errmsg)
	}
	return nil
}
