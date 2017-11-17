package gaegrpc

import (
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// HeaderKey is http.Request ID
const HeaderKey = "x-gae-grpc-id"

type requestKey struct{}

var (
	reqs = make(map[string]*http.Request)
	mu   sync.RWMutex
)

func newContextWithRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, requestKey{}, r)
}

// RequestFromContext returns *http.Request
func RequestFromContext(ctx context.Context) *http.Request {
	if r, ok := ctx.Value(requestKey{}).(*http.Request); ok {
		return r
	}
	return nil
}

func requestIDFromContext(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		v := md[HeaderKey]
		if len(v) > 0 {
			return v[0]
		}

	}
	return ""
}

func newAppContext(ctx context.Context) context.Context {
	id := requestIDFromContext(ctx)
	if id != "" {
		mu.RLock()
		r := reqs[id]
		mu.RUnlock()
		ctx = newContextWithRequest(ctx, r)
		ctx = appengine.WithContext(ctx, r)
	}
	return ctx
}

func injectAppContext() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			resp, err = handler(newAppContext(ctx), req)
			return
		}),
		grpc.StreamInterceptor(func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
			err = handler(srv, &wrapServerStream{
				ctx: newAppContext(ss.Context()),
			})
			return
		}),
	}
}

type wrapServerStream struct {
	ctx context.Context
	grpc.ServerStream
}

func (wss *wrapServerStream) Context() context.Context {
	return wss.ctx
}

// requestID returns ID, ID is pointer address.
func requestID(r *http.Request) string {
	return fmt.Sprintf("%x", &r)
}

// NewServer returns grpc.Server for App Engine
func NewServer(opt ...grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(append(injectAppContext(), opt...)...)
}

// NewRequest returns http.Request for GRPC, set the http.Request on memory
func NewRequest(r *http.Request) *http.Request {
	id := requestID(r)
	mu.Lock()
	reqs[id] = r
	mu.Unlock()
	r.Header.Add(HeaderKey, id)
	return r
}

// DeleteRequest deletes the http.Request on memory
func DeleteRequest(r *http.Request) {
	mu.Lock()
	delete(reqs, requestID(r))
	mu.Unlock()
}
