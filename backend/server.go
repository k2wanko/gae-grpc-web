package backend

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/k2wanko/gae-grpc-web/echo"
	"github.com/k2wanko/gae-grpc-web/gaegrpc"
	"google.golang.org/appengine/log"
)

var tpl *template.Template

func init() {
	l, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(l)
	}
	time.Local = l

	tpl = template.Must(template.New("").ParseFiles("index.html"))

	sv := gaegrpc.NewServer()
	echo.RegisterEchoServiceServer(sv, &EchoService{})

	wh := gaegrpc.NewWrapHandler(grpcweb.WrapServer(sv))
	http.HandleFunc("/", createAppHandler(wh))
}

func createAppHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc-web") {
			h.ServeHTTP(w, r)
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
