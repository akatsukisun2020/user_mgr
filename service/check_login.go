package service

import (
	"context"
	"time"

	"github.com/akatsukisun2020/go_components/logger"
	pb "github.com/akatsukisun2020/proto_proj/user_mgr"
	"github.com/akatsukisun2020/user_mgr/codes"
	"github.com/akatsukisun2020/user_mgr/dao"
	"github.com/akatsukisun2020/user_mgr/login"
)

// TODO: 登录态等公用的信息，看看能否抽到 cookie中. 做一个统一的filter..

// CheckLogin 用户登录校验
func CheckLogin(ctx context.Context, req *pb.CheckLoginReq) (*pb.CheckLoginRsp, error) {
	rsp := new(pb.CheckLoginRsp)
	st := time.Now()
	defer func() {
		logger.DebugContextf(ctx, "CheckLogin timecost:%d, req:%v, rsp:%v",
			time.Since(st).Milliseconds(), req, rsp)
	}()

	if req.GetUserId() == "" || req.GetAccessToken() == "" {
		logger.ErrorContextf(ctx, "CheckLogin param error, req:%v", req)
		rsp.RetCode, rsp.RetMsg = codes.ERROR_PARMA, "参数错误"
		return rsp, nil
	}

	// 查询自己的token信息，并校验
	userInfo, err := dao.NewUserInfoClient().Get(ctx, req.GetUserId())
	if err != nil {
		logger.ErrorContextf(ctx, "CheckLogin Get error, userid:%s, err:%v", req.GetUserId(), err)
		rsp.RetCode, rsp.RetMsg = codes.ERROR_QUERYREDIS, "查询存储失败"
		rsp.Result = pb.ELoginResult_LoginFailForElse
		return rsp, nil
	}
	if userInfo == nil {
		logger.ErrorContextf(ctx, "CheckLogin Get error, userid:%s, err:%v", req.GetUserId(), err)
		rsp.RetCode, rsp.RetMsg = codes.ERROR_NOUSER, "用户不存在"
		rsp.Result = pb.ELoginResult_LoginFailForNotExist
		return rsp, nil
	}

	if req.GetAccessToken() != userInfo.GetLoginInfo().GetAccessToken() {
		logger.ErrorContextf(ctx, "CheckLogin accesstoken error, userid:%s, req.accesstoken:%s, redis.accesstoken:%s",
			req.GetUserId(), req.GetAccessToken(), userInfo.GetLoginInfo().GetAccessToken())
		rsp.RetCode, rsp.RetMsg = codes.ERROR_TOKENCHECK, "TOKEN校验失败"
		rsp.Result = pb.ELoginResult_LoginFailForElse
		return rsp, nil
	}

	_, eventTime, err := login.ParseFromAccessToken(req.GetAccessToken())
	if err != nil { // 格式不对
		logger.ErrorContextf(ctx, "CheckLogin accesstoken error, userid:%s, req.accesstoken:%s, redis.accesstoken:%s",
			req.GetUserId(), req.GetAccessToken(), userInfo.GetLoginInfo().GetAccessToken())
		rsp.RetCode, rsp.RetMsg = codes.ERROR_TOKENCHECK, "TOKEN校验失败"
		rsp.Result = pb.ELoginResult_LoginFailForElse
		return rsp, nil
	}
	if eventTime+5*60*1000 < time.Now().UnixMilli() { // TODO: 5min过期做成配置形式 ==> 重要
		logger.ErrorContextf(ctx, "CheckLogin accesstoken error, accesstoken is expired, userid:%s, accesstoken:%s",
			req.GetUserId(), req.GetAccessToken())
		rsp.RetCode, rsp.RetMsg = codes.ERROR_TOKENEXPIRE, "TOKEN过期"
		rsp.Result = pb.ELoginResult_LoginFailForExpired
		return rsp, nil
	}

	return rsp, nil
}
