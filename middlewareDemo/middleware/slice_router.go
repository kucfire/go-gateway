package middleware

import (
	"context"
	"math"
	"net/http"
	"strings"
)

const abortIndex int8 = math.MaxInt8 / 2 // 最多63

type HandlerFunc func(*SliceRouterContext)

type SliceRouter struct {
	groups []*SliceGroup
}

type SliceGroup struct {
	*SliceRouter
	path     string
	handlers []HandlerFunc
}

type SliceRouterContext struct {
	RW  http.ResponseWriter
	Req *http.Request
	Ctx context.Context
	*SliceGroup
	index int8
}

//
func newSliceRouterContext(rw http.ResponseWriter, req *http.Request, r *SliceRouter) *SliceRouterContext {
	newSliceGroup := &SliceGroup{}

	// 最长url前缀匹配
	matchUrlLen := 0
	for _, group := range r.groups {
		// fmt.Println("req.RequestURI")
		// fmt.Println(req.RequestURI)
		if strings.HasPrefix(req.RequestURI, group.path) {
			pathLen := len(group.path)
			if pathLen > matchUrlLen {
				matchUrlLen = pathLen
				*newSliceGroup = *group // 浅拷贝数组指针
			}
		}
	}

	c := &SliceRouterContext{
		RW:         rw,
		Req:        req,
		SliceGroup: newSliceGroup,
		Ctx:        req.Context(),
	}
	c.Reset()
	return c
}

//
func (c *SliceRouterContext) Get(key interface{}) interface{} {
	return c.Ctx.Value(key)
}

//
func (c *SliceRouterContext) Set(key, val interface{}) {
	c.Ctx = context.WithValue(c.Ctx, key, val)
}

// 从最先加入的中间件开始回调
func (c *SliceRouterContext) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		// fmt.Println("c.index")
		// fmt.Println(c.index)
		c.handlers[c.index](c)
		c.index++
	}
}

// 跳出中间件的方法
func (c *SliceRouterContext) Abort() {
	// 直接跳到最大值63，c.next就会跳出循环
	c.index = abortIndex
}

// 是否跳过回调
func (c *SliceRouterContext) IsAborted() bool {
	return c.index >= abortIndex
}

func (c *SliceRouterContext) Reset() {
	c.index = -1
}

// 构建router
func NewSliceRouter() *SliceRouter {
	return &SliceRouter{}
}

func (g *SliceRouter) Group(path string) *SliceGroup {
	return &SliceGroup{
		SliceRouter: g,
		path:        path,
	}
}

// 构造回调方法
func (g *SliceGroup) Use(middlewares ...HandlerFunc) *SliceGroup {
	g.handlers = append(g.handlers, middlewares...)
	existsFlag := false
	for _, oldGroup := range g.SliceRouter.groups {
		if oldGroup == g {
			existsFlag = true
		}
	}
	if !existsFlag {
		g.SliceRouter.groups = append(g.SliceRouter.groups, g)
	}
	return g
}

type SliceRouterHandler struct {
	coreFunc func(*SliceRouterContext) http.Handler
	router   *SliceRouter
}

//
func (w *SliceRouterHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c := newSliceRouterContext(rw, req, w.router)
	if w.coreFunc != nil {
		c.handlers = append(c.handlers, func(c *SliceRouterContext) {
			w.coreFunc(c).ServeHTTP(rw, req)
		})
	}
	c.Reset() // index 变为-1
	c.Next()
}

//
func NewSliceRouterHandler(coreFunc func(*SliceRouterContext) http.Handler, router *SliceRouter) *SliceRouterHandler {
	return &SliceRouterHandler{
		coreFunc: coreFunc,
		router:   router,
	}
}
