package connect

import (
	"context"
	"fmt"
	"kitchen/pkg/common/logging"
	"kitchen/pkg/service"
	"kitchen/pkg/service/cors"
	"net"
	"net/http"
	"strings"

	connect_metadata "kitchen/pkg/service/connect/metadata"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"

	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// LoggerName the logger name to use for the server
const LoggerName = "connect.server"

// HandlerFactory is a factory for creating a connect handler
type HandlerFactory[T any] func(T, ...connect.HandlerOption) (string, http.Handler)

// Server is a connect server
type Server struct {
	logger           *zap.Logger
	httpServer       *http.Server
	mux              *http.ServeMux
	preStartHooks    []service.PreStartHook
	preShutdownHooks []service.PreShutdownHook
	shutdownHooks    []service.ShutdownHook
	cfg              service.Config
}

// ServiceRegistrar provides the ability to register services to a server
type ServiceRegistrar interface {
	Register(*Server, ...connect.HandlerOption) string
}

// Service is a service to mount on a server
type Service[T any] struct {
	Factory HandlerFactory[T]
	Handler T
}

// Register registers a service with the supplied server
func (s Service[T]) Register(server *Server, opts ...connect.HandlerOption) string {
	path, route := s.Factory(s.Handler, server.defaultOptions(opts)...)
	server.mux.Handle(path, route)
	return path
}

// NewServer creates a new connect server
func NewServer[T any](cfg service.Config, factory HandlerFactory[T], handler T, opts ...Option) *Server {

	// Parse the configuration options
	var options options
	for _, opt := range opts {
		opt(&options)
	}

	// Create a Server instance
	s := &Server{
		cfg:    cfg,
		logger: logging.NewLogger(LoggerName),
		mux:    http.NewServeMux(),
	}
	// Register all the services
	serviceNames := make([]string, 0, 1+len(options.additionalServices))
	services := make([]ServiceRegistrar, 0, 1+len(options.additionalServices))
	services = append(services, Service[T]{factory, handler})
	services = append(services, options.additionalServices...)

	for _, service := range services {
		path := service.Register(s, options.handlerOptions...)
		serviceName := strings.ReplaceAll(path, "/", "")

		// Register any annotations attached to the specified service
		// annotation.RegisterService(serviceName)
		serviceNames = append(serviceNames, serviceName)
	}

	// Register the reflection service
	reflector := grpcreflect.NewStaticReflector(serviceNames...)
	s.mux.Handle(grpcreflect.NewHandlerV1(reflector))
	s.mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	// Register our hooks
	if hook, ok := any(handler).(interface{ PreStart(context.Context) error }); ok {
		s.RegisterPreStartHook(hook.PreStart)
	}
	if hook, ok := any(handler).(interface{ PreShutdown() error }); ok {
		s.RegisterPreShutdownHook(hook.PreShutdown)
	}
	if hook, ok := any(handler).(interface{ Shutdown() error }); ok {
		s.RegisterShutdownHook(hook.Shutdown)
	}

	// Create the http server
	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.BindAddress, cfg.GrpcPort),
		Handler:           h2c.NewHandler(cors.New(options.corsOptions...).Handler(s.mux), &http2.Server{}),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
	}

	return s
}

// RegisterPreStartHook registers a pre-start hook
func (s *Server) RegisterPreStartHook(hooks ...service.PreStartHook) {
	s.preStartHooks = append(s.preStartHooks, hooks...)
}

// RegisterPreShutdownHook registers a pre-shutdown hook
func (s *Server) RegisterPreShutdownHook(hooks ...service.PreShutdownHook) {
	s.preShutdownHooks = append(s.preShutdownHooks, hooks...)
}

// RegisterShutdownHook registers a shutdown hook
func (s *Server) RegisterShutdownHook(hooks ...service.ShutdownHook) {
	s.shutdownHooks = append(s.shutdownHooks, hooks...)
}

// defaultOptions creates the set of default options
func (s *Server) defaultOptions(opts []connect.HandlerOption) []connect.HandlerOption {
	interceptors := make([]connect.Interceptor, 0, 4)
	if interceptor, err := otelconnect.NewInterceptor(otelconnect.WithTrustRemote()); err == nil {
		interceptors = append(interceptors, interceptor)
	}
	// interceptors = append(interceptors, connect_logging.NewInterceptor(connect_logging.Config{}))
	interceptors = append(interceptors, connect_metadata.NewInterceptor())
	defaults := make([]connect.HandlerOption, 0, len(opts)+2)
	defaults = append(defaults, connect.WithInterceptors(interceptors...), connect.WithRecover(s.recover))
	return append(defaults, opts...)
}

// Start starts this server, this function will block until Stop is called
func (s *Server) Start(ctx context.Context) error {

	s.logger.Info("starting connect service", zap.String("addr", s.cfg.BindAddress), zap.Int("port", s.cfg.GrpcPort))

	// RUn the pre-start hooks
	if err := s.runPreStartHooks(ctx); err != nil {
		return err
	}

	// Start listening for incomming connections
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.BindAddress, s.cfg.GrpcPort))
	if err != nil {
		return err
	}

	s.logger.Info("connect service listening on port", zap.String("addr", s.cfg.BindAddress), zap.Int("port", s.cfg.GrpcPort))

	// Serve our traffic
	return s.httpServer.Serve(lis)
}

// Stop stops this server
func (s *Server) Stop() error {

	s.logger.Info("stopping connect service")

	// Sync the log when complete
	defer s.logger.Sync()

	// Run the registered pre-shutdown hooks
	s.runPreShutdownHooks()

	// Stop the underlying connect server
	err := s.httpServer.Shutdown(context.Background())

	// Run the registered shutdown hooks
	s.runShutdownHooks()

	// Flush the Logger
	return err
}

// runPreStartHooks runs any registered pre-start hooks
func (s *Server) runPreStartHooks(ctx context.Context) error {
	if len(s.preStartHooks) == 0 {
		return nil
	}
	s.logger.Info("running pre-start hooks")
	for i, hook := range s.preStartHooks {
		if err := hook(ctx); err != nil {
			s.logger.Error("pre-start hook failed", zap.Int("index", i))
			return err
		}
	}
	return nil
}

// runPreShutdownHooks runs any registered pre-shutdown hooks
func (s *Server) runPreShutdownHooks() {
	if len(s.preShutdownHooks) == 0 {
		return
	}
	s.logger.Info("running pre-shutdown hooks")
	for i, hook := range s.preShutdownHooks {
		if err := hook(); err != nil {
			s.logger.Error("pre-shutdown hook failed", zap.Int("index", i))
		}
	}
}

// runShutdownHooks runs any registered shutdown hooks
func (s *Server) runShutdownHooks() {
	if len(s.shutdownHooks) == 0 {
		return
	}
	s.logger.Info("running shutdown hooks")
	for i, hook := range s.shutdownHooks {
		if err := hook(); err != nil {
			s.logger.Error("shutdown hook failed", zap.Int("index", i))
		}
	}
}

// recover is a connect recovery handler
func (s *Server) recover(ctx context.Context, spec connect.Spec, header http.Header, a any) error {
	err, _ := a.(error)
	s.logger.Error("*** PANIC ***", zap.Error(err), zap.Stack("stacktrace"), zap.String("grpc_method", spec.Procedure))
	return connect.NewError(connect.CodeInternal, err)
}
