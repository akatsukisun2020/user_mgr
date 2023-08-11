package login

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/akatsukisun2020/go_components/logger"
	pb "github.com/akatsukisun2020/proto_proj/user_mgr"
)

// EncodeToUID uid编码
func EncodeToUID(ctx context.Context, loginType pb.ELoginType, loginID string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d#%s", loginType, loginID)))
}

// ParseFromUID uid解码
func ParseFromUID(ctx context.Context, uid string) (pb.ELoginType, string, error) {
	bytes, err := base64.StdEncoding.DecodeString(uid)
	if err != nil {
		return pb.ELoginType_LoginWx, "", err
	}

	arr := strings.Split(string(bytes), "#")
	if len(arr) != 2 {
		return pb.ELoginType_LoginWx, "", fmt.Errorf("ParseFromUID: uid's format error")
	}
	loginType, err := strconv.Atoi(arr[0])
	loginID := arr[1]
	if err != nil || loginID == "" {
		return pb.ELoginType_LoginWx, "", fmt.Errorf("ParseFromUID: uid's format error")
	}

	return pb.ELoginType(loginType), loginID, nil
}

// backendLogin 后端真实的登录接口
type backendLogin interface {
	// 根据登录拼争，获取对应的后台登录信息
	Login(ctx context.Context, credential *pb.LoginCredential) (*pb.LoginInfo, error)
}

var gBackendLoginMgr map[pb.ELoginType]backendLogin

func init() {
	if gBackendLoginMgr == nil {
		gBackendLoginMgr = make(map[pb.ELoginType]backendLogin)
	}

	gBackendLoginMgr[pb.ELoginType_LoginWx] = new(wxBackendLogin)
}

// BackendLogin 后台登录接口
func BackendLogin(ctx context.Context, credential *pb.LoginCredential) (*pb.LoginInfo, error) {
	if _, ok := gBackendLoginMgr[credential.LoginType]; !ok {
		logger.ErrorContextf(ctx, "BackendLogin error, login type not existed, credential:%v", credential)
		return nil, fmt.Errorf("login type not existed")
	}

	return gBackendLoginMgr[credential.LoginType].Login(ctx, credential)
}
