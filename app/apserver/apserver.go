package main

import (
	"context"
	"github.com/eolinker/eosc/log"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type consoleServer struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	httpServer *http.Server
	grpcServer *grpc.Server
}

func newConsoleServer(httpServer *http.Server, grpcServer *grpc.Server) *consoleServer {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &consoleServer{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		httpServer: httpServer,
		grpcServer: grpcServer,
	}
}

func (c *consoleServer) close() {
	if c.cancelFunc == nil {
		return
	}
	c.cancelFunc()
	c.cancelFunc = nil

	err := c.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Errorf("关闭httpServer失败. err: %s", err)
	}
	c.grpcServer.GracefulStop()
}

func (c *consoleServer) Wait() error {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGINT)

	for {
		sig := <-sigc
		log.Infof("Caught signal pid:%d ppid:%d signal %s: .\n", os.Getpid(), os.Getppid(), sig.String())
		//log.Debug(os.Interrupt.String(), sig.String(), sig == os.Interrupt)
		switch sig {
		case os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGINT:
			{
				c.close()
				return nil
			}

		default:
			continue
		}
	}
}
