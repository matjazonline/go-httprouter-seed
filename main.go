package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"context"
	"github.com/gorilla/handlers"
	"google.golang.org/api/option"
	api "riseup.cloud/gae/api"
)

var firebaseApp *firebase.App
var firebaseAuth *auth.Client
var ADMIN_AUTHORIZATION = "ADMIN"

func main() {
}

func init() {
	_=initSerivces(context.Background())
	router := httprouter.New()
	setupAPIRoutes(router)
}

func initSerivces(ctx context.Context) error {
	// Store the startup time of the server.
	//startupTime = time.Now()

	// Initialize a Google Cloud Storage client.
	opt := option.WithCredentialsFile("riseup-cloud-firebase-adminsdk-yemdr-5bcddf820d.json")
	var err error
	firebaseApp, err = firebase.NewApp(ctx, nil, opt)
	if err != nil {
		fmt.Errorf("error initializing app: %v", err)
		return err
	}

	firebaseAuth, err = firebaseApp.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
		return err
	}

	return nil
}

func setupAPIRoutes(router *httprouter.Router) {

	var jsonAdminReqAdapters = api.AdminJsonAdaptrs(firebaseAuth, ADMIN_AUTHORIZATION, false)

	router.GET("/test", api.TestGET(api.BasicGETLifecycle()))
	router.POST("/test", api.TestPOST(jsonAdminReqAdapters))

	router.HandleMethodNotAllowed = false
	router.RedirectTrailingSlash = false
	router.HandleOPTIONS = false
	headersOk := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Referer", "User-Agent", "Authorization", "ns"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:4200"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, handlers.CORS(headersOk, originsOk, methodsOk)(router)); err != nil {
		log.Fatal(err)
	}
}

