package login

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// EncodeToAccessToken 编码AccessToken
func EncodeToAccessToken() string {
	accessToken := fmt.Sprintf("%s_%d", uuid.New().String(), time.Now().UnixMilli())
	return base64.StdEncoding.EncodeToString([]byte(accessToken)) // base64编码
}

// ParseFromAccessToken 解码AccessToken
func ParseFromAccessToken(accesstokenBase64 string) (string, int64, error) {
	bytes, err := base64.StdEncoding.DecodeString(accesstokenBase64)
	if err != nil {
		return "", 0, fmt.Errorf("accesstoken's format error, accesstokenBase64:%s, err:%s", accesstokenBase64, err.Error())
	}

	accesstoken := string(bytes)
	arr := strings.Split(string(accesstoken), "_")
	if len(arr) != 2 {
		return "", 0, fmt.Errorf("accesstoken's format error, accesstoken:%s", accesstoken)
	}

	t, err := strconv.ParseInt(arr[1], 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("accesstoken's format error, accesstoken:%s, err:%s", accesstoken, err.Error())
	}

	return arr[0], t, nil
}
