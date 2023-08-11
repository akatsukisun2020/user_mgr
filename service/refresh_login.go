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

// TODO:[重要]
// (1) 统一配置化
// (2) 统一缓存化
// (3) 登录态校验接口filter化

// RefreshLogin 刷新登录态
func RefreshLogin(ctx context.Context, req *pb.RefreshLoginReq) (*pb.RefreshLoginRsp, error) {
	rsp := new(pb.RefreshLoginRsp)
	st := time.Now()
	defer func() {
		logger.DebugContextf(ctx, "RefreshLogin timecost:%d, req:%v, rsp:%v",
			time.Since(st).Milliseconds(), req, rsp)
	}()

	if req.GetUserId() == "" || req.GetAccessToken() == "" {
		logger.ErrorContextf(ctx, "RefreshLogin param error, req:%v", req)
		rsp.RetCode, rsp.RetMsg = codes.ERROR_PARMA, "参数错误"
		return rsp, nil
	}

	// 校验登录态
	// 以后应该在nginx上嵌入这个校验功能比较好点
	checkReq := &pb.CheckLoginReq{
		UserId:      req.GetUserId(),
		AccessToken: req.GetAccessToken(),
	}
	checkRsp, err := CheckLogin(ctx, checkReq)
	if err != nil || checkRsp.GetRetCode() != 0 || checkRsp.GetResult() != pb.ELoginResult_LoginSuccess {
		logger.ErrorContextf(ctx, "RefreshLogin CheckLogin error, checkReq:%v, checkRsp:%v", checkReq, checkRsp)
		rsp.RetCode, rsp.RetMsg = codes.ERROR_TOKENCHECK, "TOKEN校验失败"
		return rsp, nil
	}

	// 生成新的token&更新存储token登录态信息
	client := dao.NewUserInfoClient()
	userInfo, err := client.Get(ctx, req.GetUserId())
	if err != nil || userInfo == nil {
		logger.ErrorContextf(ctx, "RefreshLogin Get error, userid:%s, err:%v", req.GetUserId(), err)
		rsp.RetCode, rsp.RetMsg = codes.ERROR_QUERYREDIS, "查询存储失败"
		return rsp, nil
	}
	userInfo.LoginInfo.AccessToken = login.EncodeToAccessToken()
	if err = client.Set(ctx, req.GetUserId(), userInfo); err != nil {
		logger.ErrorContextf(ctx, "RefreshLogin Set error, userid:%s, err:%v", req.GetUserId(), err)
		rsp.RetCode, rsp.RetMsg = codes.ERROR_SETREDIS, "写存储失败"
		return rsp, err
	}

	rsp.Info = userInfo.GetLoginInfo()
	return rsp, nil
}
