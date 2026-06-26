# github.com/go-chi/chi/v5 (third-party Go module)

Package chi is a small, idiomatic and composable router for building HTTP

Import path: github.com/go-chi/chi/v5   Added to workspace: chiweb

package chi // import "github.com/go-chi/chi/v5"

Package chi is a small, idiomatic and composable router for building HTTP
services.

chi supports the four most recent major versions of Go.

Example:

    package main

    import (
    	"net/http"

    	"github.com/go-chi/chi/v5"
    	"github.com/go-chi/chi/v5/middleware"
    )

    func main() {
    	r := chi.NewRouter()
    	r.Use(middleware.Logger)
    	r.Use(middleware.Recoverer)

    	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    		w.Write([]byte("root."))
    	})

    	http.ListenAndServe(":3333", r)
    }

See github.com/go-chi/chi/_examples/ for more in-depth examples.

URL patterns allow for easy matching of path components in HTTP requests.
The matching components can then be accessed using chi.URLParam(). All patterns
must begin with a slash.

A simple named placeholder {name} matches any sequence of characters up to
the next / or the end of the URL. Trailing slashes on paths must be handled
explicitly.

A placeholder with a name followed by a colon allows a regular expression match,
for example {number:\\d+}. The regular expression syntax is Go's normal regexp
RE2 syntax, except that / will never be matched. An anonymous regexp pattern
is allowed, using an empty string before the colon in the placeholder, such as
{:\\d+}

The special placeholder of asterisk matches the rest of the requested URL.
Any trailing characters in the pattern are ignored. This is the only placeholder
which will match / characters.

Examples:

    "/user/{name}" matches "/user/jsmith" but not "/user/jsmith/info" or "/user/jsmith/"
    "/user/{name}/info" matches "/user/jsmith/info"
    "/page/*" matches "/page/intro/latest"
    "/page/{other}/latest" also matches "/page/intro/latest"
    "/date/{yyyy:\\d\\d\\d\\d}/{mm:\\d\\d}/{dd:\\d\\d}" matches "/date/2017/04/01"

VARIABLES

var (
	// RouteCtxKey is the context.Context key to store the request context.
	RouteCtxKey = &contextKey{"RouteContext"}
)

FUNCTIONS

func RegisterMethod(method string)
    RegisterMethod adds support for custom HTTP method handlers, available via
    Router#Method and Router#MethodFunc

func URLParam(r *http.Request, key string) string
    URLParam returns the url parameter from a http.Request object.

func URLParamFromCtx(ctx context.Context, key string) string
    URLParamFromCtx returns the url parameter from a http.Request Context.

func Walk(r Routes, walkFn WalkFunc) error
    Walk walks any router tree that implements Routes interface.


TYPES

type ChainHandler struct {
	Endpoint http.Handler

	Middlewares Middlewares
	// Has unexported fields.
}
    ChainHandler is a http.Handler with support for handler composition and
    execution.

func (c *ChainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)

type Context struct {
	Routes Routes

	// Routing path/method override used during the route search.
	// See Mux#routeHTTP method.
	RoutePath   string
	RouteMethod string

	// URLParams are the stack of routeParams captured during the
	// routing lifecycle across a stack of sub-routers.
	URLParams RouteParams

	// Routing pattern stack throughout the lifecycle of the request,
	// across all connected routers. It is a record of all matching
	// patterns across a stack of sub-routers.
	RoutePatterns []string

	// Has unexported fields.
}
    Context is the default routing context set on the root node of a request
    context to track route patterns, URL parameters and an optional routing
    path.

func NewRouteContext() *Context
    NewRouteContext returns a new routing Context object.

func RouteContext(ctx context.Context) *Context
    RouteContext returns chi's routing Context object from a http.Request
    Context.

func (x *Context) Reset()
    Reset a routing context to its initial state.

func (x *Context) RoutePattern() string
    RoutePattern builds the routing pattern string for the particular request,
    at the particular point during routing. This means, the value will change
    throughout the execution of a request in a router. That is why it's advised
    to only use this value after calling the next handler.

    For example,

        func Instrument(next http.Handler) http.Handler {
        	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        		next.ServeHTTP(w, r)
        		routePattern := chi.RouteContext(r.Context()).RoutePattern()
        		measure(w, r, routePattern)
        	})
        }

func (x *Context) URLParam(key string) string
    URLParam returns the corresponding URL parameter value from the request
    routing context.

type Middlewares []func(http.Handler) http.Handler
    Middlewares type is a slice of standard middleware handlers with methods to
    compose middleware chains and http.Handler's.

func Chain(middlewares ...func(http.Handler) http.Handler) Middlewares
    Chain returns a Middlewares type from a slice of middleware handlers.

func (mws Middlewares) Handler(h http.Handler) http.Handler
    Handler builds and returns a http.Handler from the chain of middlewares,
    with `h http.Handler` as the final handler.

func (mws Middlewares) HandlerFunc(h http.HandlerFunc) http.Handler
    HandlerFunc builds and returns a http.Handler from the chain of middlewares,
    with `h http.Handler` as the final handler.

type Mux struct {
	// Has unexported fields.
}
    Mux is a simple HTTP route multiplexer that parses a request path, records
    any URL params, and executes an end handler. It implements the http.Handler
    interface and is friendly with the standard library.

    Mux is designed to be fast, minimal and offer a powerful API for building
    modular and composable HTTP services with a large set of handlers. It's
    particularly useful for writing large REST API services that break a handler
    into many smaller parts composed of middlewares and end handlers.

func NewMux() *Mux
    NewMux returns a newly initialized Mux object that implements the Router
    interface.

func NewRouter() *Mux
    NewRouter returns a new Mux object that implements the Router interface.

func (mx *Mux) Connect(pattern string, handlerFn http.HandlerFunc)
    Connect adds the route `pattern` that matches a CONNECT http method to
    execute the `handlerFn` http.HandlerFunc.

func (mx *Mux) Delete(pattern string, handlerFn http.HandlerFunc)
    Delete adds the route `pattern` that matches a DELETE http method to execute
    the `handlerFn` http.HandlerFunc.

func (mx *Mux) Find(rctx *Context, method, path string) string
    Find searches the routing tree for the pattern that matches the method/path.

    Note: the *Context state is updated during execution, so manage the state
    carefully or make a NewRouteContext().

func (mx *Mux) Get(pattern string, handlerFn http.HandlerFunc)
    Get adds the route `pattern` that matches a GET http method to execute the
    `handlerFn` http.HandlerFunc.

func (mx *Mux) Group(fn func(r Router)) Router
    Group creates a new inline-Mux with a copy of middleware stack. It's useful
    for a group of handlers along the same routing path that use an additional
    set of middlewares. See _examples/.

func (mx *Mux) Handle(pattern string, handler http.Handler)
    Handle adds the route `pattern` that matches any http method to execute the
    `handler` http.Handler.

func (mx *Mux) HandleFunc(pattern string, handlerFn http.HandlerFunc)
    HandleFunc adds the route `pattern` that matches any http method to execute
    the `handlerFn` http.HandlerFunc.

func (mx *Mux) Head(pattern string, handlerFn http.HandlerFunc)
    Head adds the route `pattern` that matches a HEAD http method to execute the
    `handlerFn` http.HandlerFunc.

func (mx *Mux) Match(rctx *Context, method, path string) bool
    Match searches the routing tree for a handler that matches the method/path.
    It's similar to routing a http request, but without executing the handler
    thereafter.

    Note: the *Context state is updated during execution, so manage the state
    carefully or make a NewRouteContext().

func (mx *Mux) Method(method, pattern string, handler http.Handler)
    Method adds the route `pattern` that matches `method` http method to execute
    the `handler` http.Handler.

func (mx *Mux) MethodFunc(method, pattern string, handlerFn http.HandlerFunc)
    MethodFunc adds the route `pattern` that matches `method` http method to
    execute the `handlerFn` http.HandlerFunc.

func (mx *Mux) MethodNotAllowed(handlerFn http.HandlerFunc)
    MethodNotAllowed sets a custom http.HandlerFunc for routing paths where the
    method is unresolved. The default handler returns a 405 with an empty body.

func (mx *Mux) MethodNotAllowedHandler(methodsAllowed ...methodTyp) http.HandlerFunc
    MethodNotAllowedHandler returns the default Mux 405 responder whenever a
    method cannot be resolved for a route.

func (mx *Mux) Middlewares() Middlewares
    Middlewares returns a slice of middleware handler functions.

func (mx *Mux) Mount(pattern string, handler http.Handler)
    Mount attaches another http.Handler or chi Router as a subrouter along a
    routing path. It's very useful to split up a large API as many independent
    routers and compose them as a single service using Mount. See _examples/.

    Note that Mount() simply sets a wildcard along the `pattern` that
    will continue routing at the `handler`, which in most cases is another
    chi.Router. As a result, if you define two Mount() routes on the exact same
    pattern the mount will panic.

func (mx *Mux) NotFound(handlerFn http.HandlerFunc)
    NotFound sets a custom http.HandlerFunc for routing paths that could not be
    found. The default 404 handler is `http.NotFound`.

func (mx *Mux) NotFoundHandler() http.HandlerFunc
    NotFoundHandler returns the default Mux 404 responder whenever a route
    cannot be found.

func (mx *Mux) Options(pattern string, handlerFn http.HandlerFunc)
    Options adds the route `pattern` that matches an OPTIONS http method to
    execute the `handlerFn` http.HandlerFunc.

func (mx *Mux) Patch(pattern string, handlerFn http.HandlerFunc)
    Patch adds the route `pattern` that matches a PATCH http method to execute
    the `handlerFn` http.HandlerFunc.

func (mx *Mux) Post(pattern string, handlerFn http.HandlerFunc)
    Post adds the route `pattern` that matches a POST http method to execute the
    `handlerFn` http.HandlerFunc.

func (mx *Mux) Put(pattern string, handlerFn http.HandlerFunc)
    Put adds the route `pattern` that matches a PUT http method to execute the
    `handlerFn` http.HandlerFunc.

func (mx *Mux) Route(pattern string, fn func(r Router)) Router
    Route creates a new Mux and mounts it along the `pattern` as a subrouter.
    Effectively, this is a short-hand call to Mount. See _examples/.

func (mx *Mux) Routes() []Route
    Routes returns a slice of routing information from the tree, useful for
    traversing available routes of a router.

func (mx *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request)
    ServeHTTP is the single method of the http.Handler interface that makes
    Mux interoperable with the standard library. It uses a sync.Pool to get and
    reuse routing contexts for each request.

func (mx *Mux) Trace(pattern string, handlerFn http.HandlerFunc)
    Trace adds the route `pattern` that matches a TRACE http method to execute
    the `handlerFn` http.HandlerFunc.

func (mx *Mux) Use(middlewares ...func(http.Handler) http.Handler)
    Use appends a middleware handler to the Mux middleware stack.

    The middleware stack for any Mux will execute before searching for a
    matching route to a specific handler, which provides opportunity to respond
    early, change the course of the request execution, or set request-scoped
    values for the next http.Handler.

func (mx *Mux) With(middlewares ...func(http.Handler) http.Handler) Router
    With adds inline middlewares for an endpoint handler.

type Route struct {
	SubRoutes Routes
	Handlers  map[string]http.Handler
	Pattern   string
}
    Route describes the details of a routing handler. Handlers map key is an
    HTTP method

type RouteParams struct {
	Keys, Values []string
}
    RouteParams is a structure to track URL routing parameters efficiently.

func (s *RouteParams) Add(key, value string)
    Add will append a URL parameter to the end of the route param

type Router interface {
	http.Handler
	Routes

	// Use appends one or more middlewares onto the Router stack.
	Use(middlewares ...func(http.Handler) http.Handler)

	// With adds inline middlewares for an endpoint handler.
	With(middlewares ...func(http.Handler) http.Handler) Router

	// Group adds a new inline-Router along the current routing
	// path, with a fresh middleware stack for the inline-Router.
	Group(fn func(r Router)) Router

	// Route mounts a sub-Router along a `pattern` string.
	Route(pattern string, fn func(r Router)) Router

	// Mount attaches another http.Handler along ./pattern/*
	Mount(pattern string, h http.Handler)

	// Handle and HandleFunc adds routes for `pattern` that matches
	// all HTTP methods.
	Handle(pattern string, h http.Handler)
	HandleFunc(pattern string, h http.HandlerFunc)

	// Method and MethodFunc adds routes for `pattern` that matches
	// the `method` HTTP method.
	Method(method, pattern string, h http.Handler)
	MethodFunc(method, pattern string, h http.HandlerFunc)

	// HTTP-method routing along `pattern`
	Connect(pattern string, h http.HandlerFunc)
	Delete(pattern string, h http.HandlerFunc)
	Get(pattern string, h http.HandlerFunc)
	Head(pattern string, h http.HandlerFunc)
	Options(pattern string, h http.HandlerFunc)
	Patch(pattern string, h http.HandlerFunc)
	Post(pattern string, h http.HandlerFunc)
	Put(pattern string, h http.HandlerFunc)
	Trace(pattern string, h http.HandlerFunc)

	// NotFound defines a handler to respond whenever a route could
	// not be found.
	NotFound(h http.HandlerFunc)

	// MethodNotAllowed defines a handler to respond whenever a method is
	// not allowed.
	MethodNotAllowed(h http.HandlerFunc)
}
    Router consisting of the core routing methods used by chi's Mux, using only
    the standard net/http.

type Routes interface {
	// Routes returns the routing tree in an easily traversable structure.
	Routes() []Route

	// Middlewares returns the list of middlewares in use by the router.
	Middlewares() Middlewares

	// Match searches the routing tree for a handler that matches
	// the method/path - similar to routing a http request, but without
	// executing the handler thereafter.
	Match(rctx *Context, method, path string) bool

	// Find searches the routing tree for the pattern that matches
	// the method/path.
	Find(rctx *Context, method, path string) string
}
    Routes interface adds two methods for router traversal, which is also used
    by the `docgen` subpackage to generation documentation for Routers.

type WalkFunc func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error
    WalkFunc is the type of the function called for each method and route
    visited by Walk.
