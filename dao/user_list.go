package dao

import (
	"context"
	"time"

	"github.com/akatsukisun2020/go_components/logger"
	"github.com/akatsukisun2020/go_components/redis"
)

type userListClient struct {
	redisCli *redis.RedisZset
}

func NewUserListClient() *userListClient {
	opts := redis.RedisOptions{} // TODO: 调试配置信息
	return &userListClient{
		redisCli: redis.NewRedisZset(opts),
	}
}

type UserListItem struct {
	UserID    string
	EventTime int64
}

func makeUserListKey() string {
	return "UserList20230811"
}

func (cli *userListClient) ZAdd(ctx context.Context, userid string) error {
	key := makeUserListKey()
	zsetItem := &redis.ZsetItem{Score: time.Now().UnixMilli(), Value: userid}
	if err := cli.redisCli.ZAdd(ctx, key, []*redis.ZsetItem{zsetItem}); err != nil {
		logger.ErrorContextf(ctx, "userListClient.ZAdd error, userid:%s, err:%v", userid, err)
		return err
	}

	return nil
}

func (cli *userListClient) GetAll(ctx context.Context) ([]*UserListItem, error) {
	key := makeUserListKey()

	zsetItems, err := cli.redisCli.ZRange(ctx, key, 0, -1)
	if err != nil {
		logger.ErrorContextf(ctx, "userListClient.GetAll error, err:%v", err)
		return []*UserListItem{}, err
	}

	var res []*UserListItem
	for _, item := range zsetItems {
		if userid, ok := item.Value.(string); ok {
			res = append(res, &UserListItem{
				UserID:    userid,
				EventTime: item.Score,
			})
		}
	}

	return res, nil
}
