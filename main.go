package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var progressMap = make(map[string]int)
var progressChannels = make(map[string]chan int)
var progressMutex sync.Mutex


func main() {
	ConfigRuntime()
	StartGin()
}

func envVariable(key string) string {
  err := godotenv.Load(".env")

  if err != nil {
    log.Fatalf("Error loading .env file")
  }

  return os.Getenv(key)
}


// ConfigRuntime sets the number of operating system threads.
func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}
// StartGin starts gin web server with setting router.
func StartGin() {
	gin.SetMode(gin.ReleaseMode)
	envVariable("ENV")

	router := gin.Default()
	config := cors.DefaultConfig()
	config.ExposeHeaders = []string{"Server-Timing"}
	if envVariable("ENV") == "production" {
		config.AllowOrigins = []string{envVariable("CLIENT_ORIGIN")}
	} else {
		config.AllowAllOrigins = true
	}
  	
  	router.Use(cors.New(config))

	// Serve static files for the client
	router.GET("/health", HealthCheck)

	// SSE route to receive progress updates
	router.GET("/events", EventProgressGET)
	// Set a lower memory limit for multipart forms (default is 32 MiB)
  	// router.MaxMultipartMemory = 5 << 20  // 5 MiB
	router.POST("/api/convert", convertImagePOST)

	port := envVariable("PORT") 
	if port == "" {
		port = "8080"
	}

	host := "localhost:"
	if envVariable("ENV") == "production" {
		host = "0.0.0.0:"
	}

	if err := router.Run(host + port); err != nil {
        log.Panicf("error: %s", err)
	}
}