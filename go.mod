module github.com/akatsukisun2020/user_mgr

go 1.16

replace github.com/akatsukisun2020/proto_proj => /home/painter/github/proto_proj

replace github.com/akatsukisun2020/go_components => /home/painter/github/go_components

require (
	github.com/akatsukisun2020/go_components v0.0.0-00010101000000-000000000000
	github.com/akatsukisun2020/proto_proj v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/spf13/viper v1.15.0
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.30.0
)
