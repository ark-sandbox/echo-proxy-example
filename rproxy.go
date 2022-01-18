package main

import (
	"net/http"
	"net/http/httputil"

	"github.com/labstack/echo"
)

// Proxy returns a Proxy middleware.
func RProxy(tgtUrl string, args []string,
	requestRewriter func(string, echo.Context, []string) func(*http.Request)) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {

			req := c.Request()
			res := c.Response()

			// Fix header
			// Basically it's not good practice to unconditionally pass incoming x-real-ip header to upstream.
			// However, for backward compatibility, legacy behavior is preserved unless you configure Echo#IPExtractor.
			if req.Header.Get(echo.HeaderXRealIP) == "" {
				req.Header.Set(echo.HeaderXRealIP, c.RealIP())
			}
			if req.Header.Get(echo.HeaderXForwardedProto) == "" {
				req.Header.Set(echo.HeaderXForwardedProto, c.Scheme())
			}
			if c.IsWebSocket() && req.Header.Get(echo.HeaderXForwardedFor) == "" { // For HTTP, it is automatically set by Go HTTP reverse proxy.
				req.Header.Set(echo.HeaderXForwardedFor, c.RealIP())
			}

			proxyHTTP(requestRewriter(tgtUrl, c, args)).ServeHTTP(res, req)

			if e, ok := c.Get("_error").(error); ok {
				err = e
			}

			return
		}
	}
}

func proxyHTTP(requestRewriter func(*http.Request)) http.Handler {
	director := requestRewriter
	return &httputil.ReverseProxy{Director: director}
}
