package grpc_proxy_middleware

import (
	"encoding/json"
	"fmt"
	"gatewayDemo/dao"
	"gatewayDemo/public"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// 匹配接入方式 给予请求信息
func GRPCFJwtLowLimitModeMiddleware() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("missing metadata from context")
		}

		app := md.Get("appInfo")
		if len(app) == 0 {
			if err := handler(srv, ss); err != nil {
				// log.Printf("RPC failed with error %v\n", err)
				return err
			}

			return nil
		}
		appInfo := &dao.AppInfo{}
		if err := json.Unmarshal([]byte(app[0]), appInfo); err != nil {
			return err
		}

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

		if appInfo.QPS > 0 {
			clientLimit, err := public.FlowLimiterHandler.GetFlowLimiter(
				public.FlowAppPrefix+appInfo.AppID+"_"+ClientIP,
				float64(appInfo.QPS))
			if err != nil {
				return err
			}

			if !clientLimit.Allow() {
				return errors.New(fmt.Sprintf("app client ip:%v flow limit %v", ClientIP, appInfo.QPS))
			}
		}

		if err := handler(srv, ss); err != nil {
			// log.Printf("RPC failed with error %v\n", err)
			return err
		}

		return nil
	}
}
