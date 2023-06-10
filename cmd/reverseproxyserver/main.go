package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"formal.com/go-reverse-proxy-homework/internal/config"
	"formal.com/go-reverse-proxy-homework/internal/logging"
	"formal.com/go-reverse-proxy-homework/internal/proxy"
	"formal.com/go-reverse-proxy-homework/internal/util"
	"github.com/gorilla/mux"
)

func main() {
	ctx := context.Background()
	ctx = logging.InitLogger(ctx, "./logs/application.log")

	configuration := config.ParseConfigFile("./internal/config/test_data.json")
	configuration.PrintConfigurations(ctx)

	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Failed to start listener: ", err)
	}

	router := mux.NewRouter()
	for _, cfg := range configuration.Configs {
		proxy := proxy.NewProxy(ctx, cfg)

		router.HandleFunc("/"+cfg.SourcePath, proxy)
	}

	servers := util.NewServerGroup()

	servers.Go(func() error {
		server := http.Server{
			Handler:     router,
			BaseContext: func(listener net.Listener) context.Context { return ctx },
		}
		err := server.Serve(listener)
		return fmt.Errorf("server error: %w", err)
	})

	// setup mirroring backends to test proxy
	for _, cfg := range configuration.Configs {
		backendUrl, err := url.Parse(cfg.TargetURL)
		if err != nil {
			log.Fatal("Failed to parse target URL: ", err)
		}

		backendMirroringListener, err := net.Listen("tcp", backendUrl.Host)
		if err != nil {
			log.Fatal("Failed to start listener: ", err)
		}

		servers.Go(func() error {
			server := http.Server{
				Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					for key, values := range req.Header {
						for _, value := range values {
							rw.Header().Add(key, value)
						}
					}

					_, err := io.Copy(rw, req.Body)
					if err != nil {
						log.Fatal("Failed to copy body: ", err)
					}
					err = req.Body.Close()
					if err != nil {
						log.Fatal("Failed to close body: ", err)
					}
				}),
				BaseContext: func(listener net.Listener) context.Context { return ctx },
			}
			err := server.Serve(backendMirroringListener)
			return fmt.Errorf("server error: %w", err)
		})
	}

	err = servers.Wait()
	if err != nil {
		log.Fatal("Server error: ", err)
		os.Exit(1)
	}
}
