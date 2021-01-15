package grpc_proxy_middleware

import (
	"gatewayDemo/dao"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// 匹配接入方式 给予请求信息
func GRPCHeaderTransferModeMiddleware(serviceDetail *dao.ServiceDetail) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("missing metadata from context")
		}

		for _, item := range strings.Split(serviceDetail.GRPCRule.HeaderTransfor, ",") {
			items := strings.Split(item, " ")
			if len(items) != 3 {
				continue
			}
			if items[0] == "add" || items[0] == "edit" {
				md.Set(items[1], items[2])
			}

			if items[0] == "del" {
				// 由于md是map结构，所以可以直接使用内置的delete函数进行删除map内的数据
				delete(md, items[1])
			}
		}

		if err := ss.SetHeader(md); err != nil {
			return errors.WithMessage(err, "sendHeader error")
		}

		if err := handler(srv, ss); err != nil {
			// log.Printf("RPC failed with error %v\n", err)
			return err
		}

		return nil
	}
}
