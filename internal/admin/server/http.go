package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"ktserver/internal/admin/conf"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Config, r *Router, logger log.Logger) *http.Server {

	var opts []http.ServerOption

	if c.Server.Admin.Network != "" {
		opts = append(opts, http.Network(c.Server.Admin.Network))
	}
	if c.Server.Admin.Addr != "" {
		opts = append(opts, http.Address(c.Server.Admin.Addr))
	}
	if c.Server.Admin.Timeout != 0 {
		opts = append(opts, http.Timeout(c.Server.Admin.Timeout))
	}

	srv := http.NewServer(opts...)
	r.installRouter() // install router
	srv.HandlePrefix("/", r.g)
	return srv
}
