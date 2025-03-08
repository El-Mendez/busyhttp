package main

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime"
	"net/http"
	"os"
	"strconv"
	"time"
)

func ping(c *gin.Context) {
	c.Status(http.StatusOK)
}

func timeData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"uptime_ms": time.Since(startupTime).Milliseconds(),
		"uptime":    time.Since(startupTime).String(),
		"date":      time.Now(),
		"timezone":  time.Local,
	})
}

func echo(c *gin.Context) {
	var jsonBody map[string]interface{}
	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		if err.Error() == "EOF" {
			jsonBody = nil // Default to empty JSON
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"headers":    c.Request.Header,
		"method":     c.Request.Method,
		"ip":         c.Request.RemoteAddr,
		"detectedIP": c.ClientIP(),
		"body":       jsonBody,
		"path":       c.Request.URL.Path,
	})
}

func wait(c *gin.Context) {
	sleepTime := c.Param("ms")
	if sleepTime == "" {
		c.String(http.StatusBadRequest, "did not send 'ms' parameter")
		return
	}

	intTime, err := strconv.Atoi(sleepTime)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("%v is not a valid ms parameter. Expected a number", sleepTime))
	}

	time.Sleep(time.Duration(intTime) * time.Millisecond)
	c.String(http.StatusOK, "ok")
}

func info(c *gin.Context) {
	path, _ := os.Getwd()
	hostname, _ := os.Hostname()

	c.JSON(http.StatusOK, gin.H{
		"hostname": hostname,
		"path":     path,
		"gid":      os.Getgid(),
		"uid":      os.Getuid(),
		"pid":      os.Getpid(),
		"ppid":     os.Getppid(),
		"env":      os.Environ(),
	})
}

func crash(c *gin.Context) {
	os.Exit(1)
}

func isReady(c *gin.Context) {
	if ready.Load() {
		c.String(http.StatusOK, "ok")
	} else {
		c.String(http.StatusBadRequest, "not ready")
	}
}

func exit(c *gin.Context) {
	c.String(http.StatusOK, "crashed")
	os.Exit(0)
}

func statusCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.String(http.StatusBadRequest, "did not send 'code' parameter")
		return
	}

	intCode, err := strconv.Atoi(code)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("%v is not a valid code parameter. Expected a number", code))
	}

	c.String(intCode, fmt.Sprintf("status %v", intCode))
}

type HelpItem struct {
	Name        string
	Description string
}

func help(c *gin.Context) {
	variables := []HelpItem{
		{"BUSY_STARTUP_TIME_MS", "The startup time for the http server"},
		{"BUSY_CRASH", "Immediately crashes after startup time. Anything other than 0 activates it."},
		{"BUSY_SHUTDOWN_TIME_MS", "The time after SIGINT needed to start shutting down the http server"},
		{"BUSY_TRUSTED_PROXIES", "Trusted proxies CIDR"},
		{"BUSY_READY_TIME_MS", "Time after the http server is up needed for the ready endpoint to return 200."},
		{"BUSY_TRUSTED_PLATFORM", "The header for IP mapping with trusted proxies"},
		{"BUSY_SECRET", "If auth is needed, you can set this variable and it must be sent in the Bearer format"},
		{"BUSY_ADDRESS", "The server address, default :8080"},
	}
	endpoints := []HelpItem{
		{"/ping", "returns 200 Ok"},
		{"/help", "Shows this message"},
		{"/ready", "returns 200 (depends on BUSY_READY_TIME_MS)"},
		{"/info", "Returns info about this server"},
		{"/time", "Returns info about time data"},
		{"/echo", "Returns info about the request received"},
		{"/crash", "Immediately exits with 1."},
		{"/wait/:ms", "Waits for :ms and then returns 200."},
		{"/exit", "Immediately exits with 0"},
		{"/status/:code", "Returns the :code status code"},
		{"/file?filename=:filename", "Gets or creates a file"},
	}
	c.HTML(http.StatusOK, "index", gin.H{
		"title":     "Debug Info",
		"variables": variables,
		"endpoints": endpoints,
	})
}

func readFile(c *gin.Context) {
	filename := c.Query("filename")
	content, err := os.ReadFile(filename)

	mimeType := mime.TypeByExtension(SplitAtLast(filename, "."))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err, "path": filename})
		return
	}

	c.Data(http.StatusOK, mimeType, content)
}

func uploadFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	location := c.Query("location")

	if file == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No location provided"})
		return
	}

	if err := c.SaveUploadedFile(file, location); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("File %s uploaded successfully", file.Filename),
	})
}
