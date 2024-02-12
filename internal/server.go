package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kaato137/quickrest/internal/pkg/filewatch"
	"github.com/kaato137/quickrest/internal/pkg/rwhandler"
)

type Server struct {
	cfg      *Config
	cfgMutex sync.RWMutex

	mux         *rwhandler.RWHandler
	renderer    *Renderer
	reqRecorder *RequestRecorder

	logger *slog.Logger
}

func NewServerFromConfig(cfg *Config) (*Server, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	s := &Server{cfg: cfg, logger: logger}

	if err := s.setupMux(); err != nil {
		return nil, fmt.Errorf("setup mux: %w", err)
	}

	s.renderer = NewRenderer()
	s.reqRecorder = NewRequestRecorder(cfg.RecordDir)

	return s, nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(rw, r)
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
	return filewatch.WatchFilePath(s.cfg.Path).
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
			s.logger.Error("Error in file watcher", err)
			return true
		}).
		Run(context.Background())
}

func (s *Server) handleResponse(route RouteConfig) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if route.Latency > 0 || route.Jitter > 0 {
			if err := s.waitLatency(r, route); err != nil {
				s.logger.Error("Failed during waiting latency", err)
				return
			}
		}

		rw.Header().Add("Content-Type", route.ContentType)
		for k, v := range route.Headers {
			rw.Header().Set(k, v)
		}

		rw.WriteHeader(route.StatusCode)

		if err := s.renderBody(rw, r, route); err != nil {
			s.logger.Error("Failed to render body", err)
			return
		}

		if route.Record {
			if err := s.reqRecorder.Record(formatRouteFilename(route), r); err != nil {
				s.logger.Error("Failed to record request", err)
				return
			}
		}
	}
}

func (s *Server) renderBody(rw http.ResponseWriter, r *http.Request, route RouteConfig) error {
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

	newCfg, err := LoadConfigFromFile(s.cfg.Path)
	if err != nil {
		return fmt.Errorf("load config from file: %w", err)
	}

	s.cfg = newCfg

	return nil
}

func (s *Server) waitLatency(r *http.Request, route RouteConfig) error {
	select {
	case <-time.After(calcWaitDuration(route)):
		return nil
	case <-r.Context().Done():
		return r.Context().Err()
	}
}

func formatRouteFilename(route RouteConfig) string {
	date := time.Now().Format("2006-01-02")
	rt := strings.ReplaceAll(route.Path, "/", " ")

	return fmt.Sprintf("%s-%s.log", rt, date)
}

func formatResponseBody(rc RouteConfig, r *http.Request) []byte {
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

func prepareRenderContext(rc RouteConfig, r *http.Request) RenderContext {
	ctx := make(RenderContext)
	for _, wc := range rc.Wildcards {
		ctx[wc] = r.PathValue(wc)
	}

	return ctx
}

func calcWaitDuration(route RouteConfig) time.Duration {
	waitDuration := route.Latency

	if route.Jitter > 0 {
		waitDuration += time.Duration(rand.Int63n(int64(route.Jitter*2))) - route.Jitter
	}

	return waitDuration
}
