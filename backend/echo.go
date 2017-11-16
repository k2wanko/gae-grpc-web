package backend

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/k2wanko/gae-grpc-web/echo"
)

// EchoService structs
type EchoService struct{}

// Echo implements EchoServiceServer
func (*EchoService) Echo(ctx context.Context, req *echo.EchoRequest) (res *echo.EchoResponse, err error) {
	grpc.SendHeader(ctx, metadata.Pairs("Pre-Response-Metadata", "Is-sent-as-headers-unary"))
	grpc.SetTrailer(ctx, metadata.Pairs("Post-Response-Metadata", "Is-sent-as-trailers-unary"))
	msg := req.GetMessage()
	logf(ctx, "Echo Message = %s", msg)
	res = &echo.EchoResponse{
		Message: msg,
	}
	return
}
