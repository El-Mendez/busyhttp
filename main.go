package main

import (
	"context"
	_ "embed"
	"errors"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

var ready atomic.Bool
var startupTime time.Time

//go:embed help.html
var helpHtml string

func main() {
	startupWaitTime := validateIntFlag("BUSY_STARTUP_TIME_MS")
	shutdownTime := validateIntFlag("BUSY_SHUTDOWN_TIME_MS")
	if startupWaitTime != 0 {
		log.Printf("STARTING UP: wait for startup time %v\n", startupWaitTime)
		time.Sleep(time.Duration(startupWaitTime) * time.Millisecond)
	} else {
		log.Printf("STARTING UP: no startup time %v\n", startupWaitTime)
	}
	startupTime = time.Now()

	shouldCrash := os.Getenv("BUSY_CRASH")
	if shouldCrash != "" && shouldCrash != "0" {
		panic("BUSY_CRASH panic")
	} else {
		log.Printf("STARTING UP: not set to crash\n")
	}

	readyTime := validateIntFlag("BUSY_READY_TIME_MS")
	if readyTime != 0 {
		go func() {
			log.Printf("STARTING UP: not ready. wait for %v", readyTime)
			time.Sleep(time.Duration(readyTime) * time.Millisecond)
			log.Printf("STARTING UP: ready! timeout finished\n")
			ready.Store(true)
		}()
	} else {
		log.Printf("STARTING UP: ready! no timeout %v\n", startupWaitTime)
		ready.Store(true)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	trustedProxies := strings.Split(os.Getenv("BUSY_TRUSTED_PROXIES"), ",")
	if len(trustedProxies) == 1 && trustedProxies[0] == "" {
		log.Printf("STARTING UP: no trusted proxies\n")
		r.SetTrustedProxies(nil)
	} else {
		err := r.SetTrustedProxies(trustedProxies)
		if err != nil {
			log.Fatalf("Invalid trusted proxies")
		}
		log.Printf("STARTING UP: trusted proxies %v\n", trustedProxies)
	}
	r.TrustedPlatform = os.Getenv("BUSY_TRUSTED_PLATFORM")
	log.Printf("STARTING UP: trusted platform is '%v'\n", r.TrustedPlatform)

	r.Use(gin.Logger())
	var busySecret = os.Getenv("BUSY_SECRET")
	if busySecret != "" {
		log.Printf("STARTING UP: bearer secret is set!\n")
		r.Use(AuthMiddleware(os.Getenv("BUSY_SECRET")))
	} else {
		log.Printf("STARTING UP: no bearer secret is set!\n")
	}
	r.SetHTMLTemplate(template.Must(template.New("index").Parse(helpHtml)))

	specifyRoutes := r.Group("/specify/:instance")
	createRouter(specifyRoutes)

	topRoutes := r.Group("/")
	createRouter(topRoutes)

	var busyAddress = os.Getenv("BUSY_ADDRESS")
	if busyAddress == "" {
		busyAddress = ":8080"
	}

	log.Printf("setting up server on %v\n", busyAddress)
	srv := &http.Server{
		Addr:    busyAddress,
		Handler: r.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	if shutdownTime != 0 {
		log.Println("Starting shutdown time")
		time.Sleep(time.Duration(shutdownTime) * time.Millisecond)
	}

	log.Println("Starting Shutdown ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
