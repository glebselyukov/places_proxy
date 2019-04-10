package api

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"

	"github.com/prospik/places_proxy/pkg/conv"
	logging "github.com/prospik/places_proxy/pkg/logger"
)

type route struct {
	methods [][]byte
	Handler
	f handleFunc
}

type handleFunc func(ctx *fasthttp.RequestCtx, log logging.Logger)

// Handlers http controller
type Handler interface {
	// Places handler
	Places(ctx *fasthttp.RequestCtx, log logging.Logger)
}

// http methods
var (
	GET  = []byte("GET")
	POST = []byte("POST")
)

// Router helper tool for request executing
type Router struct {
	routes map[string]*route
	log    logging.Logger
}

// NewRouter constructor for Router
func NewRouter(log logging.Logger) *Router {
	return &Router{
		routes: make(map[string]*route),
		log:    log,
	}
}

// RegisterPlacesRoutes registers paths belonging to certain operations.
func (r *Router) RegisterPlacesRoutes() {
	places := NewPlacesHandler()
	r.Register("/api/places", places, places.Places, POST)
}

// Register add new route to routes map
func (r *Router) Register(path string, handler Handler, handle handleFunc, methods ...[]byte) {
	if r.routes == nil {
		r.routes = make(map[string]*route)
	}
	r.routes[path] = &route{
		Handler: handler,
		methods: methods,
		f:       handle,
	}
}

// ServeHTTP entry point for fasthttp server
func (r *Router) ServeHTTP(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(xContentTypeOptionsHeader, xContentTypeOptionsNosniff)
	requestID := ctx.Request.Header.Peek(requestHeader)
	if requestID == nil {
		requestID = []byte(fmt.Sprintf("%d", ctx.ID()))
	}
	ctx.Response.Header.SetBytesV(requestHeader, requestID)
	headerTags := make([]*logging.Tag, 0, ctx.Request.Header.Len())
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		headerTags = append(headerTags, logging.Any(conv.B2S(key), conv.B2S(value)))
	})

	requestBody := ctx.Request.Body()
	requestURI := ctx.URI().FullURI()
	log := r.log.Copy(
		logging.Any("request_id", conv.B2S(requestID))).
		Copy(headerTags...).
		Copy(logging.Any("body", conv.B2S(requestBody)))

	log.Info(fmt.Sprintf("%s", conv.B2S(requestURI)))

	defer func() {
		r.corsHandler(ctx)
	}()

	if ctx.IsOptions() || ctx.IsHead() {
		return
	}

	path := strings.ToLower(conv.B2S(ctx.Path()))
	if len(path) > 0 && path[len(path)-1:] == "/" {
		path = path[:len(path)-1]
	}
	route, ok := r.routes[path]
	if !ok {
		r.errorHandler(ctx, notFoundError)
		return
	}

	var isAllowedMethod bool
	for _, method := range route.methods {
		if bytes.Equal(ctx.Method(), method) {
			isAllowedMethod = true
			break
		}
	}
	if !isAllowedMethod {
		r.errorHandler(ctx, methodNotAllowedError)
		return
	}

	route.f(ctx, log)

	acceptContentType := ctx.Request.Header.Peek(acceptHeader)
	if acceptContentType != nil && ctx.Response.Header.ContentType() == nil {
		switch strings.ToLower(conv.B2S(acceptContentType)) {
		case strings.ToLower(jsonContentType):
			ctx.Response.Header.SetContentType(jsonContentType)
		default:
			ctx.Response.Header.SetContentType(textContentType)
		}
	}
}
