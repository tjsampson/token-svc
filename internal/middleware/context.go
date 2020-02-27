package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

type key int

const (
	requestIDKey          key = 0
	requestStartTimeKey   key = 1
	userIPKey             key = 2
	jwtIDKey              key = 3
	tokenSvcRequestHeader     = "X-Request-ID"
)

// userIPFromRequest extracts the user IP address from req, if present.
func userIPFromRequest(req *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	return userIP, nil
}

// NewJWTIDContext returns a new Context carrying the JWT ID.
func newJWTIDContext(ctx context.Context, jti string) context.Context {
	return context.WithValue(ctx, jwtIDKey, jti)
}

// NewUserIPContext returns a new Context carrying userIP.
func newUserIPContext(ctx context.Context, userIP net.IP) context.Context {
	return context.WithValue(ctx, userIPKey, userIP)
}

// NewStartTimeContext returns a new Context carrying startTime.
func newStartTimeContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, requestStartTimeKey, time.Now())
}

// NewRequestIDContext returns a new Context carrying startTime.
func newRequestIDContext(req *http.Request) context.Context {
	ctx := req.Context()
	reqID := req.Header.Get(tokenSvcRequestHeader)
	if reqID == "" {
		reqID = uuid.NewV4().String()
	}
	return context.WithValue(ctx, requestIDKey, reqID)
}

func newWrappedReqCtx(req *http.Request) context.Context {
	userIP, _ := userIPFromRequest(req)
	return newUserIPContext(newStartTimeContext(newRequestIDContext(req)), userIP)
}

// JWTIDFromContext returns the JWTD ID from the context
func JWTIDFromContext(ctx context.Context) string {
	return ctx.Value(jwtIDKey).(string)
}

// UserIPFromContext extracts the user IP address from ctx, if present.
func UserIPFromContext(ctx context.Context) net.IP {
	// ctx.Value returns nil if ctx has no value for the key;
	// the net.IP type assertion returns ok=false for nil.
	return ctx.Value(userIPKey).(net.IP)
}

// RequestIDFromContext returns the requestID from the http Context
func RequestIDFromContext(ctx context.Context) string {
	return ctx.Value(requestIDKey).(string)
}

// RequestStartTimeFromContext returns the requestStartTime from the http Context
func RequestStartTimeFromContext(ctx context.Context) time.Time {
	return ctx.Value(requestStartTimeKey).(time.Time)
}
