package grpc_proxy_middleware

import (
	"encoding/json"
	"fmt"
	"gatewayDemo/dao"
	"gatewayDemo/public"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// 匹配接入方式 给予请求信息
func GRPCJwtFLowCountModeMiddleware() grpc.StreamServerInterceptor {
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

		// 统计项
		// 1. 全站统计
		// 2. 服务统计
		// 3. 租户统计
		// 租户统计
		appCounter, err := public.FlowCounterHandler.GetFlowCounter(public.FlowAppPrefix + appInfo.AppID)
		if err != nil {
			return err
		}
		appCounter.Increase()
		if appInfo.QPD > 0 && appCounter.TotalCount > appInfo.QPD {
			return errors.New(fmt.Sprintf("该租户已限流 limit : %v current : %v",
				appInfo.QPD,
				appCounter.TotalCount))
		}

		// dayappCount, _ := appCounter.GetDayData(time.Now())
		// fmt.Printf("{%s serviceCount —— qps:%v total:%v}\n", appInfo.AppID, appCounter.QPS, appCounter.TotalCount)

		if err := handler(srv, ss); err != nil {
			// log.Printf("RPC failed with error %v\n", err)
			return err
		}

		return nil
	}
}
