package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mwitkow/go-conntrack"
	"github.com/mwitkow/go-flagz"
	"github.com/mwitkow/go-flagz/protobuf"
	"github.com/mwitkow/go-httpwares/logging/logrus"
	"github.com/mwitkow/go-httpwares/tags"
	"github.com/mwitkow/go-httpwares/tracing/debug"
	validator "github.com/mwitkow/go-proto-validators"
	pb_config "github.com/mwitkow/kedge/_protogen/winch/config"
	"github.com/mwitkow/kedge/lib/map"
	"github.com/mwitkow/kedge/lib/sharedflags"
	"github.com/mwitkow/kedge/winch"
	"github.com/pressly/chi"
	log "github.com/sirupsen/logrus"
	_ "golang.org/x/net/trace" // so /debug/request gets registered.
)

var (
	flagHttpPort = sharedflags.Set.Int("server_http_port", 8070, "TCP port to listen on for HTTP1.1/REST calls.")

	flagHttpMaxWriteTimeout = sharedflags.Set.Duration("server_http_max_write_timeout", 10*time.Second, "HTTP server config, max write duration.")
	flagHttpMaxReadTimeout  = sharedflags.Set.Duration("server_http_max_read_timeout", 10*time.Second, "HTTP server config, max read duration.")
	flagMapperConfig        = protoflagz.DynProto3(sharedflags.Set,
		"server_mapper_config",
		&pb_config.MapperConfig{},
		"Contents of the Winch Mapper configuration. Content or read from file if _path suffix.").
		WithFileFlag("../../misc/winch_mapper.json").WithValidator(validateMapper)
	flagAuthConfig = protoflagz.DynProto3(sharedflags.Set,
		"server_auth_config",
		&pb_config.AuthConfig{},
		"Contents of the Winch Auth configuration. Content or read from file if _path suffix.").
		WithFileFlag("../../misc/winch_auth.json").WithValidator(validateMapper)
)

func validateMapper(msg proto.Message) error {
	if val, ok := msg.(validator.Validator); ok {
		if err := val.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := sharedflags.Set.Parse(os.Args); err != nil {
		log.WithError(err).Fatal("failed parsing flags")
	}
	if err := flagz.ReadFileFlags(sharedflags.Set); err != nil {
		log.WithError(err).Fatal("failed reading flagz from files")
	}

	log.SetOutput(os.Stdout)
	logEntry := log.NewEntry(log.StandardLogger())

	routes, err := winch.NewStaticRoutes(
		flagMapperConfig.Get().(*pb_config.MapperConfig),
		flagAuthConfig.Get().(*pb_config.AuthConfig),
	)
	if err != nil {
		log.WithError(err).Fatal("failed reading flagz from files")
	}

	var httpPlainListener net.Listener
	httpPlainListener = buildListenerOrFail("http_plain", *flagHttpPort)
	log.Infof("listening for HTTP Plain on: %v", httpPlainListener.Addr().String())

	http.Handle("/debug/flagz", http.HandlerFunc(flagz.NewStatusEndpoint(sharedflags.Set).ListFlags))

	pacHandle, err := winch.NewPacFromFlags(httpPlainListener.Addr().String())
	if err != nil {
		log.WithError(err).Fatalf("failed to init PAC handler")
	}

	http.Handle("/wpad.dat", pacHandle)
	http.Handle("/", winch.New(kedge_map.RouteMapper(routes.Get()),
		// TODO(bplotka): Only for debug purposes.
		&tls.Config{InsecureSkipVerify: true},
	))
	winchServer := &http.Server{
		WriteTimeout: *flagHttpMaxWriteTimeout,
		ReadTimeout:  *flagHttpMaxReadTimeout,
		ErrorLog:     http_logrus.AsHttpLogger(logEntry),
		Handler: chi.Chain(
			http_ctxtags.Middleware("winch"),
			http_debug.Middleware(),
			http_logrus.Middleware(logEntry),
		).Handler(http.DefaultServeMux),
	}

	// Serve.
	if err := winchServer.Serve(httpPlainListener); err != nil {
		log.Fatalf("http_plain server error: %v", err)
	}
}

func buildListenerOrFail(name string, port int) net.Listener {
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed listening for '%v' on %v: %v", name, port, err)
	}
	return conntrack.NewListener(listener,
		conntrack.TrackWithName(name),
		conntrack.TrackWithTcpKeepAlive(20*time.Second),
		conntrack.TrackWithTracing(),
	)
}