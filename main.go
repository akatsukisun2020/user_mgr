package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"regexp"

	commonConf "github.com/akatsukisun2020/go_components/config"
	"github.com/akatsukisun2020/go_components/logger"
	serverinit "github.com/akatsukisun2020/go_components/server_init"
	pb "github.com/akatsukisun2020/proto_proj/user_mgr"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// 提供http接口服务
func main() {
	ctx := context.Background()

	// 【系统初始化】自定义的服务初始流程
	serverinit.ServerInit()

	s := grpc.NewServer(grpc.UnaryInterceptor(logger.UnaryLoggerInterceptor))
	pb.RegisterUserMgrHttpServer(s, &UserMgr{})

	// 1. 启动grpc服务
	systemConf := commonConf.GetSystemConfig()
	trpcTaget := fmt.Sprintf(":%d", systemConf.ServerConfig.TrpcPort)
	lis, err := net.Listen("tcp", trpcTaget)
	if err != nil {
		logger.FatalContextf(ctx, "failed to listen:%v", err)
	}
	go func() {
		if err := s.Serve(lis); err != nil {
			logger.FatalContextf(ctx, "failed to serve:%v", err)
		}
	}()

	// 2. 启动http服务
	httpTaget := fmt.Sprintf(":%d", systemConf.ServerConfig.HttpPort)
	conn, err := grpc.Dial(trpcTaget, grpc.WithInsecure())
	if err != nil {
		logger.FatalContextf(ctx, "failed to dial server:%v", err)
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterUserMgrHttpHandler(ctx, gwmux, conn)
	if err != nil {
		logger.FatalContextf(ctx, "Failed to register gateway:%v", err)
	}
	gwServer := &http.Server{
		Addr:    httpTaget,
		Handler: cors(gwmux), // 支持跨域的方式。
	}

	logger.InfoContextf(ctx, "Serving grpc-gateway on http, Addr: %s", httpTaget)
	logger.FatalContextf(ctx, "ListenAndServe ERROR :%v", gwServer.ListenAndServe())
}

func allowedOrigin(origin string) bool {
	if viper.GetString("cors") == "*" {
		return true
	}
	if matched, _ := regexp.MatchString(viper.GetString("cors"), origin); matched {
		return true
	}
	return false
}

// http自定义方式支持跨域，参考：https://fale.io/blog/2021/07/28/cors-headers-with-grpc-gateway ==> 具体原理还不知道！！
func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allowedOrigin(r.Header.Get("Origin")) {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
		}
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
