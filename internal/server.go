package internal

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kaato137/quickrest/internal/conf"
	"github.com/kaato137/quickrest/internal/pkg/filewatch"
	"github.com/kaato137/quickrest/internal/pkg/rwhandler"
)

type Server struct {
	cfg      *conf.Config
	cfgMutex sync.RWMutex

	reqID uint64

	mux         *rwhandler.RWHandler
	renderer    *Renderer
	reqRecorder *RequestRecorder

	closers []func()

	logger Logger
}

func NewServerFromConfig(cfg *conf.Config) (*Server, error) {
	logger := NewLogger()

	s := &Server{cfg: cfg, logger: logger}

	if err := s.setupMux(); err != nil {
		return nil, fmt.Errorf("setup mux: %w", err)
	}

	s.renderer = NewRenderer()
	s.reqRecorder = NewRequestRecorder(cfg.RecordDir)

	s.logger.Info("Listen on", "addr", cfg.Address)

	return s, nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(rw, r)
}

func (s *Server) Close() {
	for _, closeFn := range s.closers {
		closeFn()
	}
}

func (s *Server) setupMux() error {
	router := s.setupRouter()
	s.mux = rwhandler.New(router)

	if err := s.setupConfigReload(); err != nil {
		return fmt.Errorf("setup config reload: %w", err)
	}

	return nil
}

func (s *Server) setupRouter() *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range s.cfg.Routes {
		mux.HandleFunc(route.Path, s.handleResponse(route))
	}

	return mux
}

func (s *Server) setupConfigReload() error {
	s.logger.Info("Setup config reload", "interval", s.cfg.ReloadInterval)
	closer, err := filewatch.WatchFilePath(s.cfg.Path).
		WithInterval(s.cfg.ReloadInterval).
		OnChange(func() error {
			s.logger.Info("Config changed. Reloading...")

			if err := s.reloadConfigFile(); err != nil {
				return err
			}
			s.mux.SetHandler(s.setupRouter())

			s.logger.Info("Config reloaded successfully")

			return nil
		}).
		OnError(func(err error) bool {
			s.logger.Error("Error on config reload. Keeping old configuration", "err", err)
			return true
		}).
		Run(context.Background())

	s.appendCloser(closer)

	return err
}

func (s *Server) handleResponse(route conf.RouteConfig) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		now := time.Now()
		reqID := s.ReqID()
		s.logger.Info("Request started", "id", reqID, "method", r.Method, "url", r.URL.String())
		defer func(now time.Time) {
			took := time.Since(now)
			s.logger.Info("Request ended", "id", reqID, "took", took, "code", route.StatusCode)
		}(now)

		if route.Latency > 0 || route.Jitter > 0 {
			if err := s.waitLatency(r, route); err != nil {
				s.logger.Error("Failed during waiting latency", "err", err)
				return
			}
		}

		rw.Header().Add("Content-Type", route.ContentType)
		for k, v := range route.Headers {
			rw.Header().Set(k, v)
		}

		rw.WriteHeader(route.StatusCode)

		if err := s.renderBody(rw, r, route); err != nil {
			s.logger.Error("Failed to render body", "err", err)
			return
		}

		if route.Record {
			if err := s.reqRecorder.Record(formatRouteFilename(route), r); err != nil {
				s.logger.Error("Failed to record request", "err", err)
				return
			}
		}
	}
}

func (s *Server) renderBody(rw http.ResponseWriter, r *http.Request, route conf.RouteConfig) error {
	var (
		body []byte
		err  error
	)
	if route.BodyJS != "" {
		renderCtx := prepareRenderContext(route, r)
		body, err = s.renderer.Render(route.BodyJS, renderCtx)
		if err != nil {
			return fmt.Errorf("render js template: %w", err)
		}
	} else {
		body = formatResponseBody(route, r)
	}

	if _, err := rw.Write(body); err != nil {
		return fmt.Errorf("write body: %w", err)
	}

	return nil
}

func (s *Server) reloadConfigFile() error {
	s.cfgMutex.Lock()
	defer s.cfgMutex.Unlock()

	newCfg, err := conf.LoadConfigFromFile(s.cfg.Path)
	if err != nil {
		return fmt.Errorf("load config from file: %w", err)
	}

	s.cfg = newCfg

	return nil
}

func (s *Server) waitLatency(r *http.Request, route conf.RouteConfig) error {
	select {
	case <-time.After(calcWaitDuration(route)):
		return nil
	case <-r.Context().Done():
		return r.Context().Err()
	}
}

func (s *Server) ReqID() uint64 {
	return atomic.AddUint64(&s.reqID, 1)
}

func (s *Server) appendCloser(fn func()) {
	s.closers = append(s.closers, fn)
}

func formatRouteFilename(route conf.RouteConfig) string {
	date := time.Now().Format("2006-01-02")
	rt := strings.ReplaceAll(route.Path, "/", " ")

	return fmt.Sprintf("%s-%s.log", rt, date)
}

func formatResponseBody(rc conf.RouteConfig, r *http.Request) []byte {
	resolvedBody := rc.Body
	for _, c := range rc.Wildcards {
		new := r.PathValue(c)

		if new == "" {
			continue
		}

		old := fmt.Sprintf("{%s}", c)

		resolvedBody = strings.ReplaceAll(resolvedBody, old, new)
	}

	return []byte(resolvedBody)
}

func prepareRenderContext(rc conf.RouteConfig, r *http.Request) RenderContext {
	ctx := make(RenderContext)
	for _, wc := range rc.Wildcards {
		ctx[wc] = r.PathValue(wc)
	}

	return ctx
}

func calcWaitDuration(route conf.RouteConfig) time.Duration {
	waitDuration := route.Latency

	if route.Jitter > 0 {
		waitDuration += time.Duration(rand.Int63n(int64(route.Jitter*2))) - route.Jitter
	}

	return waitDuration
}
