package grpc_proxy_middleware

import (
	"gatewayDemo/dao"
	"gatewayDemo/public"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// jwt auth token
func GRPCJwtOauthTokenModeMiddleware(serviceDetail *dao.ServiceDetail) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("missing metadata from context")
		}

		// decode jwt token
		// app_id 与app_list 取得 appInfo
		// appInfo 放到 gin.context
		auth := md.Get("Authorization")
		authToken := ""
		if len(auth) > 0 {
			authToken = auth[0]
		}
		token := strings.ReplaceAll(authToken, "Bearer ", "")
		matched := false
		if token != "" {
			claim, err := public.JwtDecode(token)
			if err != nil {
				return errors.WithMessage(err, "JwtDecode")
			}
			// fmt.Println("claims.Issuer : ", claim.Issuer)
			appList := dao.AppManagerHandler.GetAppList()
			for _, appInfo := range appList {
				if appInfo.AppID == claim.Issuer {
					md.Set("appInfo", public.ObjToJson(appInfo))
					matched = true
					break
				}
			}
		}

		if serviceDetail.AccessControl.OpenAuth == 1 && !matched {
			return errors.New("not mactched vail app")
		}

		if err := handler(srv, ss); err != nil {
			// log.Printf("RPC failed with error %v\n", err)
			return err
		}

		return nil
	}
}
