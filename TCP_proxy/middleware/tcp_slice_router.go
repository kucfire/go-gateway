package middleware

import (
	"context"
	"go-gateway/TCP_proxy/server"
	"math"
	"net"
)

// 最多 63 个中间件
var abortIndex int8 = math.MaxInt8 / 2

//
type TCPHandlerFunc func(*TCPSliceRouterContext)

// router 结构体
type TCPSliceRouter struct {
	groups []*TCPSliceGroup
}

// group 结构体
type TCPSliceGroup struct {
	*TCPSliceRouter
	path     string
	handlers []TCPHandlerFunc
}

// router 上下文
type TCPSliceRouterContext struct {
	conn net.Conn
	Ctx  context.Context
	*TCPSliceGroup
	index int8
}

func newTcpSliceRouterContext(conn net.Conn, r *TCPSliceRouter, ctx context.Context) *TCPSliceRouterContext {
	newTcpSliceGroup := &TCPSliceGroup{}
	*newTcpSliceGroup = *r.groups[0] //浅拷贝数组指针,只会使用第一个分组
	c := &TCPSliceRouterContext{
		conn:          conn,
		TCPSliceGroup: newTcpSliceGroup,
		Ctx:           ctx}
	c.Reset()
	return c
}

func (c *TCPSliceRouterContext) Get(key interface{}) interface{} {
	return c.Ctx.Value(key)
}

func (c *TCPSliceRouterContext) Set(key, val interface{}) {
	c.Ctx = context.WithValue(c.Ctx, key, val)
}

type TCPSliceRouterHandler struct {
	coreFunc func(*TCPSliceRouterContext) server.TCPHandler
	router   *TCPSliceRouter
}

func (w *TCPSliceRouterHandler) ServeTCP(ctx context.Context, conn net.Conn) {
	c := newTcpSliceRouterContext(conn, w.router, ctx)
	c.handlers = append(c.handlers, func(c *TCPSliceRouterContext) {
		w.coreFunc(c).ServeTCP(ctx, conn)
	})
	c.Reset()
	c.Next()
}

func NewTcpSliceRouterHandler(coreFunc func(*TCPSliceRouterContext) server.TCPHandler, router *TCPSliceRouter) *TCPSliceRouterHandler {
	return &TCPSliceRouterHandler{
		coreFunc: coreFunc,
		router:   router,
	}
}

// 构造 router
func NewTcpSliceRouter() *TCPSliceRouter {
	return &TCPSliceRouter{}
}

// 创建 Group
func (g *TCPSliceRouter) Group(path string) *TCPSliceGroup {
	if path != "/" {
		panic("only accept path=/")
	}
	return &TCPSliceGroup{
		TCPSliceRouter: g,
		path:           path,
	}
}

// 构造回调方法
func (g *TCPSliceGroup) Use(middlewares ...TCPHandlerFunc) *TCPSliceGroup {
	g.handlers = append(g.handlers, middlewares...)
	existsFlag := false
	for _, oldGroup := range g.TCPSliceRouter.groups {
		if oldGroup == g {
			existsFlag = true
		}
	}
	if !existsFlag {
		g.TCPSliceRouter.groups = append(g.TCPSliceRouter.groups, g)
	}
	return g
}

// 从最先加入中间件开始回调
func (c *TCPSliceRouterContext) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// 跳出中间件方法
func (c *TCPSliceRouterContext) Abort() {
	c.index = abortIndex
}

// 是否跳过了回调
func (c *TCPSliceRouterContext) IsAborted() bool {
	return c.index >= abortIndex
}

// 重置回调
func (c *TCPSliceRouterContext) Reset() {
	c.index = -1
}
