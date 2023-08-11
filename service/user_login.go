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

// UserLogin 用户登录
func UserLogin(ctx context.Context, req *pb.UserLoginReq) (*pb.UserLoginRsp, error) {
	rsp := new(pb.UserLoginRsp)
	st := time.Now()
	defer func() {
		logger.DebugContextf(ctx, "UserLogin timecost:%d, req:%v, rsp:%v",
			time.Since(st).Milliseconds(), req, rsp)
	}()

	if req.GetCredential().GetCode() == "" {
		logger.ErrorContextf(ctx, "UserLogin param error, req:%v", req)
		rsp.RetCode, rsp.RetMsg = codes.ERROR_PARMA, "参数错误"
		return rsp, nil
	}

	// 登录
	loginInfo, err := login.BackendLogin(ctx, req.GetCredential())
	if err != nil {
		logger.ErrorContextf(ctx, "BackendLogin error, req:%v, err:%v", req, err)
		return rsp, err
	}

	// 更新登录信息
	err = modifyLoginInfo(ctx, loginInfo)
	if err != nil {
		logger.ErrorContextf(ctx, "modifyLoginInfo error, req:%v, err:%v", req, err)
		return rsp, err
	}

	return rsp, nil
}

func getUidFromLoginInfo(ctx context.Context, loginInfo *pb.LoginInfo) string {
	var loginID string
	switch loginInfo.GetLoginType() {
	case pb.ELoginType_LoginWx:
		loginID = loginInfo.GetOpenId()
	}

	return login.EncodeToUID(ctx, loginInfo.GetLoginType(), loginID)
}

// 修改登录信息
func modifyLoginInfo(ctx context.Context, loginInfo *pb.LoginInfo) error {
	userid := getUidFromLoginInfo(ctx, loginInfo)

	// redis查找
	client := dao.NewUserInfoClient()
	userInfo, err := client.Get(ctx, userid)
	if err != nil {
		logger.ErrorContextf(ctx, "modifyLoginInfo Get error, userid:%s, err:%v", userid, err)
		return err
	}

	if userInfo == nil { // 还没登录过
		userInfo = &pb.UserInfo{
			UserId: userid,
		}
		if err = dao.NewUserListClient().ZAdd(ctx, userid); err != nil {
			logger.ErrorContextf(ctx, "modifyLoginInfo ZAdd error, userid:%s, err:%v", userid, err)
			// 非关键
		}
	}
	userInfo.LoginInfo = loginInfo

	if err = client.Set(ctx, userid, userInfo); err != nil {
		logger.ErrorContextf(ctx, "modifyLoginInfo Set error, userid:%s, err:%v", userid, err)
		return err
	}

	return nil
}
