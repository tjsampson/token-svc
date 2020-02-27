package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	internalerrors "github.com/tjsampson/token-svc/internal/errors"
	"github.com/tjsampson/token-svc/internal/log"
	"github.com/tjsampson/token-svc/internal/serviceprovider"
	"github.com/tjsampson/token-svc/pkg/metrics"

	"github.com/opentracing-contrib/go-gorilla/gorilla"
	nethttp "github.com/opentracing-contrib/go-stdlib/nethttp"
	"go.uber.org/zap"
)

// TimeoutHandler will timeout the request if the TokenSvcTimeout is exceeded (currently 30 seconds)
// This will result in a 503 HTTP Error
func TimeoutHandler(timeoutSecs uint16) Adapter {
	return func(h http.Handler) http.Handler {
		return http.TimeoutHandler(h, time.Duration(timeoutSecs)*time.Second, fmt.Sprintf("Sorry! Service Timeout Exceeded: %v", timeoutSecs))
	}
}

// TracingHandler will adding jaeger tracing to all call
func TracingHandler(appCtxProvider *serviceprovider.Context) Adapter {
	return func(h http.Handler) http.Handler {
		return gorilla.Middleware(
			appCtxProvider.TraceProvider.Tracer,
			h,
			nethttp.OperationNameFunc(func(r *http.Request) string {
				return "HTTP " + r.Method + " " + r.RequestURI
			}),
		)
	}
}

// LogMetricsHandler adapts the incoming request with Logging/Metrics (Observability)
func LogMetricsHandler(logger log.Factory, metricProvider *metrics.Provider) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Establish new Context with Request-ID, StartTime, UserIP, etc.
			// (continue to add "contextual data" here)
			ctx := newWrappedReqCtx(r)

			// Log the incoming Request
			logger.For(ctx).Info("incoming request",
				zap.Time("request_start", RequestStartTimeFromContext(ctx)),
				zap.String("request_id", RequestIDFromContext(ctx)),
				zap.String("user_agent", r.UserAgent()),
				zap.Int64("content_length_bytes", r.ContentLength),
				zap.String("host", r.Host),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("uri", r.RequestURI),
				zap.String("method", r.Method),
				zap.String("user_ip", UserIPFromContext(ctx).String()),
				zap.String("protocol", r.Proto),
			)
			// Tally the incoming request metrics
			metricProvider.StatHTTPRequestCount.WithLabelValues(r.RequestURI, r.Method, r.Proto).Inc()
			metricProvider.StatRequestSaturationGuage.WithLabelValues(r.RequestURI, r.Method, r.Proto).Inc()

			// Run this on the way out (i.e. outgoing response)
			defer func() {
				// Calculate the request duration (i.e. latency)
				ms := float64(time.Since(RequestStartTimeFromContext(ctx))) / float64(time.Millisecond)

				// Log the outgoing response
				logger.For(ctx).Info("outgoing response",
					zap.String("request_id", RequestIDFromContext(ctx)),
					zap.Float64("duration_ms", ms),
				)

				// Tally the outgoing response metrics
				metricProvider.StatRequestDurationHistogram.WithLabelValues(r.RequestURI, r.Method, r.Proto).Observe(ms)
				metricProvider.StatRequestDurationGuage.WithLabelValues(r.RequestURI, r.Method, r.Proto).Set(ms)
				metricProvider.StatRequestSaturationGuage.WithLabelValues(r.RequestURI, r.Method, r.Proto).Dec()

			}()
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AuthHandler handles api auth
// first we check for open routes (which are configurable)
// if route is not open then we check the JWT and the HTTPS Cookie
// we extract the JWT, validate it's contents, then extract/decode the cookie data
// we also compare the JWT ID inside the JWT and what was baked into the secure cookie
// if all is good, we create a new context with the JTI value
func AuthHandler(appCtx *serviceprovider.Context) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			var isOK bool

			// PERF: This should be a Map
			for _, openEndPoint := range appCtx.Config.API.OpenEndPoints {
				if openEndPoint == r.RequestURI {
					isOK = true
					break
				}
			}
			if isOK {
				h.ServeHTTP(w, r)
			} else {

				invalidAuth := func(err error) {
					appCtx.Logger.For(ctx).Error("AuthHandler - Invalid Auth", zap.Error(err))
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(internalerrors.RestError{Message: "invalid auth", Code: http.StatusUnauthorized})
				}

				validAuth := func(tokenID string) {
					ctx := newJWTIDContext(ctx, tokenID)
					appCtx.Logger.For(ctx).Info("AuthHandler - Authenticated")
					h.ServeHTTP(w, r.WithContext(ctx))
				}

				// TODO: Abstract this out and use go routines
				if jwtToken := extractAuthBearerToken(r); len(jwtToken) > 0 {
					if tokenClaims, validToken := appCtx.JwtClient.IsValidAccessToken(ctx, jwtToken); validToken {
						if cookies := appCtx.CookieOven.DecodedCookie(ctx, r); cookies != nil {
							if cookies[appCtx.Config.Cookie.KeyJWTAccessID] == tokenClaims.Id {
								// check if user creds are valid
								user, err := appCtx.UserRepo.ReadByEmail(ctx, cookies[appCtx.Config.Cookie.KeyEmail])
								if err != nil {
									invalidAuth(err)
									return
								}
								cacheJTI, err := appCtx.RedisClient.Get(ctx, fmt.Sprintf("%v-%v", appCtx.Config.Token.AccessCacheKeyID, user.ID))
								if cacheJTI == tokenClaims.Id {
									validAuth(tokenClaims.Id)
									return
								}
							}
						}
					}
				}
				invalidAuth(fmt.Errorf("Invalid JWT Token"))
				return
			}
		})
	}
}

// RouteHandlerSig is the route handler signature
type RouteHandlerSig func(appCtxProvider *serviceprovider.Context, res http.ResponseWriter, req *http.Request) (int, interface{}, error)

// Handler is the wrapper that provides context to the app handler
type Handler struct {
	AppCtx       *serviceprovider.Context
	RouteHandler RouteHandlerSig
}
