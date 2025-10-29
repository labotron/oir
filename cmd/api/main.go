package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/suttapak/odtimagereplacer"
)

func main() {
	// Parse command-line flags
	port := flag.String("port", "8080", "Server port")
	host := flag.String("host", "0.0.0.0", "Server host")
	mode := flag.String("mode", "release", "Gin mode: debug, release, or test")
	flag.Parse()

	// Set Gin mode
	switch *mode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup router
	router := odtimagereplacer.SetupRouter()

	// Server address
	addr := fmt.Sprintf("%s:%s", *host, *port)

	// Print startup information
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║         ODT Image Replacer API Server                   ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Printf("  Mode:    %s\n", *mode)
	fmt.Printf("  Address: http://%s\n", addr)
	fmt.Println("\n  Endpoints:")
	fmt.Println("    POST /api/replace          - Replace images (JSON response)")
	fmt.Println("    POST /api/replace/download - Replace images (file download)")
	fmt.Println("    GET  /health               - Health check")
	fmt.Println("    GET  /info                 - Service information")
	fmt.Println("\n  Example request:")
	fmt.Println(`    curl -X POST http://localhost:8080/api/replace \`)
	fmt.Println(`      -H "Content-Type: application/json" \`)
	fmt.Println(`      -d @example.json`)
	fmt.Println("\n══════════════════════════════════════════════════════════")

	// Start server
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
