package proxy

import (
	"context"
	"github.com/donscoco/gateway/police/load_balance"
	"github.com/donscoco/gateway/proxy/grpc_proxy_ext"
	"google.golang.org/grpc"
	"log"
)

func NewGrpcLoadBalanceHandler(lb load_balance.LoadBalance) grpc.StreamHandler {
	return func() grpc.StreamHandler {
		director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
			nextAddr, err := lb.Get("")
			if err != nil {
				log.Fatal("get next addr fail")
			}

			c, err := grpc.DialContext(ctx, nextAddr, grpc.WithCodec(grpc_proxy_ext.Codec()), grpc.WithInsecure())
			return ctx, c, err
		}
		return grpc_proxy_ext.TransparentHandler(director)
	}()
}
