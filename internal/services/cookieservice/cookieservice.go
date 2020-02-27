package cookieservice

import (
	"context"
	"net/http"

	"github.com/tjsampson/token-svc/internal/config"
	"github.com/tjsampson/token-svc/internal/log"

	"github.com/gorilla/securecookie"
	"go.uber.org/zap"
)

// Oven is the interface for the cookieoven
type Provider interface {
	BakeCookie(ctx context.Context, cookieChan chan *http.Cookie, cookieData map[string]string)
	DecodedCookie(ctx context.Context, r *http.Request) map[string]string
}

type provider struct {
	logger    log.Factory
	cfg       *config.Config
	secCookie *securecookie.SecureCookie
}

// New returns a new implementation of the Cookie.Oven interface
func New(cfg *config.Config, logger log.Factory) Provider {
	return &provider{
		logger:    logger.With(zap.String("package", "cookieoven")),
		cfg:       cfg,
		secCookie: securecookie.New([]byte(cfg.Cookie.BlockKey), []byte(cfg.Cookie.BlockKey)),
	}
}

// DecodeCookie attempted to decode the incoming cookie
// be sure and check for valid nil (cookie doesn't exist)
func (p *provider) DecodedCookie(ctx context.Context, r *http.Request) map[string]string {
	if cookie, err := r.Cookie(p.cfg.Cookie.Name); err == nil {
		value := make(map[string]string)
		if err = p.secCookie.Decode(p.cfg.Cookie.Name, cookie.Value, &value); err == nil {
			return value
		}
	}
	return nil
}

func (p *provider) BakeCookie(ctx context.Context, cookieChan chan *http.Cookie, cookieData map[string]string) {
	p.secCookie.MaxAge(86400 * int(p.cfg.Cookie.LifeSpanDays))
	p.secCookie.SetSerializer(securecookie.JSONEncoder{})
	encodedCookieData, err := p.secCookie.Encode(p.cfg.Cookie.Name, cookieData)

	if err != nil {
		p.logger.For(ctx).Error("failed BakeCookie Encode", zap.Error(err))
		cookieChan <- &http.Cookie{}
		close(cookieChan)
		return
	}

	cookie := &http.Cookie{
		Name:     p.cfg.Cookie.Name,
		Value:    encodedCookieData,
		Path:     "/",
		HttpOnly: true,
		Domain:   p.cfg.Cookie.Domain,
		Secure:   true,
	}

	cookieChan <- cookie
	close(cookieChan)
}
