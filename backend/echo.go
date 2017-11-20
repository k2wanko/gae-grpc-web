package backend

import (
	"time"

	"github.com/k2wanko/gae-grpc-web/echo"
	"github.com/kjk/betterguid"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// EchoService structs
type EchoService struct{}

// Echo structs message
type Echo struct {
	ID      *datastore.Key `datastore:"-"`
	Message string
	Created time.Time
}

// ToMessage returns Echo of protocl buffer
func (e *Echo) ToMessage() *echo.Echo {
	id := ""
	if e.ID != nil {
		id = e.ID.StringID()
	}
	return &echo.Echo{
		Id:      id,
		Message: e.Message,
		Created: e.Created.Unix(),
	}
}

// Echo implements EchoServiceServer
func (*EchoService) Echo(ctx context.Context, req *echo.EchoRequest) (res *echo.EchoResponse, err error) {
	grpc.SendHeader(ctx, metadata.Pairs("Pre-Response-Metadata", "Is-sent-as-headers-unary"))
	grpc.SetTrailer(ctx, metadata.Pairs("Post-Response-Metadata", "Is-sent-as-trailers-unary"))
	msg := req.GetMessage()

	logf(ctx, "Echo Message = %s", msg)

	k := datastore.NewKey(ctx, "Echo", betterguid.New(), 0, nil)
	e := &Echo{
		ID:      k,
		Message: req.GetMessage(),
		Created: time.Now(),
	}

	if _, err := datastore.Put(ctx, k, e); err != nil {
		errorf(ctx, "don't save message: %v", err)
		return nil, err
	}

	res = &echo.EchoResponse{
		Echo: e.ToMessage(),
	}

	return
}

// EchoHistory implements EchoServiceServer
func (*EchoService) EchoHistory(req *echo.EchoHistoryRequest, ss echo.EchoService_EchoHistoryServer) (err error) {
	limit := int(req.GetLimit())
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	ctx := ss.Context()

	it := datastore.NewQuery("Echo").
		Limit(limit).
		Order("-Created").
		Run(ctx)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		e := &Echo{}
		k, err := it.Next(e)
		if err == datastore.Done {
			break
		} else if err != nil {
			errorf(ctx, "don't get history: %v", err)
			return err
		}
		e.ID = k
		ss.Send(&echo.EchoResponse{
			Echo: e.ToMessage(),
		})
	}

	return
}
