package httputil

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"sync"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/netutil/urlutil"
	"github.com/AdguardTeam/golibs/service"
)

// loggerKeyServer is the key used by [Server] to identify itself.
const loggerKeyServer = "server"

// ServerConfig is the configuration structure for *Server.
type ServerConfig struct {
	// BaseLogger is used to create the initial logger for the server.  It must
	// not be nil.
	BaseLogger *slog.Logger

	// Server is the underlying http server which contains most of
	// configuration.  It must not be nil.
	Server *http.Server

	// InitialAddress is the initial address for the server.  It may have a zero
	// port, in which case the real port will be set in [Server.Start].  It must
	// be set.
	InitialAddress netip.AddrPort
}

// Server contains an *http.Server as well as entities and data associated with
// it.
type Server struct {
	baseLogger *slog.Logger
	http       *http.Server
	listener   net.Listener
	logger     *slog.Logger
	// mu protects http, logger, listener, and url.
	mu          *sync.Mutex
	url         *url.URL
	initialAddr netip.AddrPort
}

// NewServer returns properly initialized *Server that is ready to serve HTTP
// queries.  The TCP listener is not started.  c must not be nil and must be
// valid.
func NewServer(c *ServerConfig) (s *Server) {
	u := &url.URL{
		Scheme: urlutil.SchemeHTTP,
		Host:   c.InitialAddress.String(),
	}

	if c.Server.TLSConfig != nil {
		u.Scheme = urlutil.SchemeHTTPS
	}

	logger := c.BaseLogger.With(loggerKeyServer, u)

	return &Server{
		mu:          &sync.Mutex{},
		http:        c.Server,
		logger:      logger,
		url:         u,
		baseLogger:  c.BaseLogger,
		initialAddr: c.InitialAddress,
	}
}

// LocalAddr returns the local address of the server if it has started
// listening.  Otherwise, it returns nil.
func (s *Server) LocalAddr() (addr net.Addr) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if l := s.listener; l != nil {
		return l.Addr()
	}

	return nil
}

// type check
var _ service.Interface = (*Server)(nil)

// Start implements [service.Interface] for *Server.  It blocks if the server
// starts successfully.
func (s *Server) Start(ctx context.Context) (err error) {
	defer slogutil.RecoverAndLog(ctx, s.logger)

	tcpListener, err := net.ListenTCP("tcp", net.TCPAddrFromAddrPort(s.initialAddr))
	if err != nil {
		s.logger.ErrorContext(ctx, "listening tcp", slogutil.KeyError, err)

		return fmt.Errorf("listening tcp: %w", err)
	}

	var listener net.Listener
	if s.http.TLSConfig == nil {
		listener = tcpListener
	} else {
		listener = tls.NewListener(tcpListener, s.http.TLSConfig)
	}

	func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		s.listener = listener

		// Reassign the address in case the port was zero.
		s.url.Host = listener.Addr().String()
		s.logger = s.baseLogger.With(loggerKeyServer, s.url)
		s.http.ErrorLog = slog.NewLogLogger(s.logger.Handler(), slog.LevelDebug)
	}()

	s.logger.InfoContext(ctx, "starting")

	err = s.http.Serve(listener)
	if err == nil || errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	s.logger.ErrorContext(ctx, "serving", slogutil.KeyError, err)

	return fmt.Errorf("serving: %w", err)
}

// Shutdown implements [service.Interface] for *Server.
func (s *Server) Shutdown(ctx context.Context) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var errs []error
	err = s.http.Shutdown(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("shutting down server %s: %w", s.url, err))
	}

	// Close the listener separately, as it might not have been closed if the
	// context has been canceled.
	//
	// NOTE:  The listener could remain uninitialized if [net.ListenTCP] failed
	// during server start.
	if l := s.listener; l != nil {
		err = l.Close()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			errs = append(errs, fmt.Errorf("closing listener for server %s: %w", s.url, err))
		}
	}

	return errors.Join(errs...)
}
