package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jailtonjunior94/go-products/configs"
	_ "github.com/jailtonjunior94/go-products/docs"
	"github.com/jailtonjunior94/go-products/internal/entity"
	"github.com/jailtonjunior94/go-products/internal/infra/database"
	"github.com/jailtonjunior94/go-products/internal/infra/webserver/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	swagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title          Go Products API
// @version        1.0
// @description    Product API with authentication
// @termsOfService http://swagger.io/terms

// @contact.name  Jailton Junior
// @contact.url   http://jailton.junior.net
// @contact.email jailton.junior94@outlook.com

// @license.name Jailton Junior License
// @license.url  http://jailton.junior.net

// @BasePath                   /
// @securityDefinitions.apiKey ApiKeyAuth
// @in                         header
// @name                       Authorization
func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB, configs.TokenAuthKey, configs.JWTExpiresIn)

	r := chi.NewRouter()
	r.Use(LogRequest)

	r.Post("/token", userHandler.GetJWT)
	r.Post("/users", userHandler.CreateUser)

	r.Route("/products", func(r chi.Router) {
		r.Use(Authorization)
		r.Use(jwtauth.Authenticator)

		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.GetProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	r.Get("/docs/*", swagger.Handler(swagger.URL("http://localhost:8000/docs/doc.json")))

	http.ListenAndServe(":8000", r)
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[REQUEST] %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretKey := "-----BEGIN CERTIFICATE-----\n" +
			"" +
			"\n-----END CERTIFICATE-----"

		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(secretKey))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenReq := tokenFromHeader(r)
		token, err := jwt.Parse(tokenReq, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("Token inválido")
			}
			return key, nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println("token is valid")
		}

		next.ServeHTTP(w, r)
	})
}

func tokenFromHeader(r *http.Request) string {
	// Get token from authorization header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}
