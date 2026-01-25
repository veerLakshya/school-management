package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	mw "school-management/internal/api/middlewares"
	"school-management/internal/api/router"
	"school-management/internal/repository/sqlconnect"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error-------", err)
		return
	}

	_, err = sqlconnect.ConnectDB()

	if err != nil {
		log.Println("Error-------", err)
		return
	}

	port := os.Getenv("API_PORT")

	// cert := "cert.pem"
	// key := "key.pem"

	mux := router.NewRouter()

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// rl := mw.NewRateLimiter(5, time.Minute)

	// hppOptions := mw.HPPOptions{
	// 	CheckQuery:                  true,
	// 	CheckBody:                   true,
	// 	CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
	// 	Whitelist:                   []string{"sortBy", "name", "age", "class"},
	// }

	// secureMux := mw.Cors(rl.Middleware(mw.ResponseTime(mw.SecurityHeaders(mw.Compression(mw.Hpp(hppOptions)(mux))))))
	// function to properly chain middlewares
	// secureMux := utils.ApplyMiddlewares(mux, mw.Hpp(hppOptions), mw.Compression, mw.SecurityHeaders, mw.ResponseTime, rl.Middleware, mw.Cors)
	secureMux := mw.SecurityHeaders(mux)

	//custom server
	server := &http.Server{
		Addr:      port,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port: ", port)
	err = server.ListenAndServe()
	// err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error starting the server: ", err)
	}
}
