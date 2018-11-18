package main 
import (
    "github.com/gorilla/mux"
    "github.com/gorilla/handlers"
    "github.com/form-with-jwt/services"
    "github.com/subosito/gotenv"
    "os"
    "net/http"
    "fmt"
    "github.com/form-with-jwt/controllers"
)
func main() {
    router := mux.NewRouter()

    router.HandleFunc("/api/user/signup", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/user/upload", controllers.UploadFile).Methods("POST")
    router.Use(services.JwtAuthentication)
    originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
    // loading secret creds from .env
    gotenv.Load()

    port := os.Getenv("PORT")
	if port == "" {
		port = "8080" //localhost
    }

    fmt.Println(port)

	err := http.ListenAndServe(":" + port,   handlers.CORS(originsOk)(router))
	if err != nil {
		fmt.Print(err)
}

}