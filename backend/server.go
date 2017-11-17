package backend

import (
	"context"
	"html/template"
	"net/http"
	"strings"

	"github.com/k2wanko/gae-grpc-web/gaegrpc"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/k2wanko/gae-grpc-web/echo"
	"google.golang.org/appengine/log"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").ParseFiles("index.html"))

	sv := gaegrpc.NewServer()
	echo.RegisterEchoServiceServer(sv, &EchoService{})

	wgs := grpcweb.WrapServer(sv)
	http.HandleFunc("/", createAppHandler(wgs))
}

func createAppHandler(wgs *grpcweb.WrappedGrpcServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc-web") {
			wgs.ServeHTTP(w, gaegrpc.NewRequest(r))
			gaegrpc.DeleteRequest(r)
		} else {
			serverTop(w, r)
		}
	}
}

func serverTop(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func logf(ctx context.Context, format string, args ...interface{}) {
	log.Infof(ctx, format, args...)
}

func debugf(ctx context.Context, format string, args ...interface{}) {
	log.Debugf(ctx, format, args...)
}

func warnf(ctx context.Context, format string, args ...interface{}) {
	log.Warningf(ctx, format, args...)
}

func errorf(ctx context.Context, format string, args ...interface{}) {
	log.Errorf(ctx, format, args...)
}
