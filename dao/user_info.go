package dao

import (
	"context"
	"fmt"
	"reflect"

	"github.com/akatsukisun2020/go_components/logger"
	"github.com/akatsukisun2020/go_components/redis"
	pb "github.com/akatsukisun2020/proto_proj/user_mgr"
	svrConf "github.com/akatsukisun2020/user_mgr/config"
	"google.golang.org/protobuf/proto"
)

type userInfoClient struct {
	redisCli *redis.RedisKV
}

func NewUserInfoClient() *userInfoClient {
	opts := redis.RedisOptions{
		redis.WithAddrOption(svrConf.GetUserConfig().RedisAddr),
		redis.WithPasswordOption(svrConf.GetUserConfig().RedisPasswd),
	}
	return &userInfoClient{
		redisCli: redis.NewRedisKV(opts),
	}
}

func makeuserInfoKey(userid string) string {
	return "UserInfo20230811_" + userid
}

func (cli *userInfoClient) Get(ctx context.Context, userid string) (*pb.UserInfo, error) {
	key := makeuserInfoKey(userid)
	reply, err := cli.redisCli.Get(ctx, key)
	if err != nil {
		logger.ErrorContextf(ctx, "userInfoClient.Get Get error, userid:%s, err:%v", userid, err)
		return nil, err
	}

	if reply == nil { // 目标数据不存在
		return nil, nil
	}

	bytesData, ok := reply.([]byte)
	if !ok {
		logger.ErrorContextf(ctx, "userInfoClient'reflectType error, userid:%s, type:%v", userid, reflect.TypeOf(reply))
		return nil, fmt.Errorf("userInfoClient'reflectType error, userid:%s, type:%v", userid, reflect.TypeOf(reply))
	}

	userInfo := new(pb.UserInfo)
	if err = proto.Unmarshal(bytesData, userInfo); err != nil {
		logger.ErrorContextf(ctx, "userInfoClient.Get Unmarshal error, userid:%s, err:%v", userid, err)
		return nil, err
	}

	return userInfo, nil
}

func (cli *userInfoClient) Set(ctx context.Context, userid string, userinfo *pb.UserInfo) error {
	byteDatas, err := proto.Marshal(userinfo)
	if err != nil {
		logger.ErrorContextf(ctx, "userInfoClient.Set Marshal error, userid:%s, err:%v", userid, err)
		return err
	}

	key := makeuserInfoKey(userid)
	if err = cli.redisCli.Set(ctx, key, byteDatas); err != nil {
		logger.ErrorContextf(ctx, "userInfoClient.Set Set error, userid:%s, err:%v", userid, err)
		return err
	}

	return nil
}

func (cli *userInfoClient) Del(ctx context.Context, userid string) error {
	key := makeuserInfoKey(userid)
	return cli.Del(ctx, key)
}
