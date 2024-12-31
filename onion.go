package onion

import (
    "fmt"
    "net/http"
    "strings"
)

type HandlerFunc func(*Context)

type Context struct {
    Response http.ResponseWriter
    Request  *http.Request
}

func (c *Context) String(statusCode int, msg string) {
    c.Response.WriteHeader(statusCode)
    c.Response.Write([]byte(msg))
}

type App struct {
    mux         *http.ServeMux
    middlewares []HandlerFunc
    notFound    HandlerFunc
}

type Route struct {
    Method  string
    Pattern string
    Handler HandlerFunc
}

func New() *App {
    return &App{
        mux:         http.NewServeMux(),
        middlewares: []HandlerFunc{},
        notFound: func(c *Context) {
            http.NotFound(c.Response, c.Request)
        },
    }
}

func (a *App) Use(mw HandlerFunc) {
    a.middlewares = append(a.middlewares, mw)
}

func (a *App) NotFoundHandler(fn HandlerFunc) {
    a.notFound = fn
}

func (a *App) MapRoutes(routeGroups ...[]Route) {
    for _, groupRoutes := range routeGroups {
        for _, r := range groupRoutes {
            a.handle(r.Method, r.Pattern, r.Handler)
        }
    }
}

func (a *App) handle(method, pattern string, handler HandlerFunc) {
    a.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
        if r.Method != method {
            a.notFound(&Context{Response: w, Request: r})
            return
        }

        c := &Context{Response: w, Request: r}
        for _, mw := range a.middlewares {
            mw(c)
        }

        handler(c)
    })
}


func (a *App) Group(prefix string) *Group {
    g := &Group{
        prefix:    prefix,
        parentApp: a,
    }
    return g
}

func (a *App) Run(addr string) error {
    fmt.Printf("Onion server listening on %s
", addr)
    return http.ListenAndServe(addr, a.mux)
}

type Group struct {
    prefix    string
    parentApp *App
}

func (g *Group) Use(mw HandlerFunc) {
    g.parentApp.Use(mw)
}

func (g *Group) MapRoutes(routes []Route) {
    for _, r := range routes {
        fullPattern := g.prefix + r.Pattern
        g.parentApp.handle(r.Method, fullPattern, r.Handler)
    }
}
