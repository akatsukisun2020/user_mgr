package main

import (
	"context"

	pb "github.com/akatsukisun2020/proto_proj/user_mgr"
	"github.com/akatsukisun2020/user_mgr/service"
)

type UserMgr struct {
	pb.UnimplementedUserMgrHttpServer
}

func (s *UserMgr) UserRegister(ctx context.Context, req *pb.UserRegisterReq) (*pb.UserRegisterRsp, error) {
	return nil, nil
}

func (s *UserMgr) UserLogin(ctx context.Context, req *pb.UserLoginReq) (*pb.UserLoginRsp, error) {
	return service.UserLogin(ctx, req)
}

func (s *UserMgr) CheckLogin(ctx context.Context, req *pb.CheckLoginReq) (*pb.CheckLoginRsp, error) {
	return service.CheckLogin(ctx, req)
}

func (s *UserMgr) RefreshLogin(ctx context.Context, req *pb.RefreshLoginReq) (*pb.RefreshLoginRsp, error) {
	return service.RefreshLogin(ctx, req)
}
