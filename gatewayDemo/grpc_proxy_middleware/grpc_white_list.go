package grpc_proxy_middleware

import (
	"fmt"
	"gatewayDemo/dao"
	"gatewayDemo/public"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// 匹配接入方式 给予请求信息
func GRPCWhiteListModeMiddleware(serviceDetail *dao.ServiceDetail) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// 获取来源的ClientIP
		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("peer not find with context")
		}
		split := strings.Split(peerCtx.Addr.String(), ":")
		ClientIP := ""
		if len(split) == 2 {
			ClientIP = split[0]
		}

		ipList := []string{}
		if serviceDetail.AccessControl.WhiteList != "" {
			ipList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && len(ipList) > 0 {
			if !public.InStringSlice(ipList, ClientIP) {
				return errors.New(fmt.Sprintf("%s not in white ip list", ClientIP))
			}
		}

		if err := handler(srv, ss); err != nil {
			// log.Printf("RPC failed with error %v\n", err)
			return err
		}

		return nil
	}
}
