package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gtldhawalgandhi/go-training/3.Intermediate/db"
	"github.com/gtldhawalgandhi/go-training/3.Intermediate/token"
)

// Server will server HTTP requests
type Server struct {
	store   db.Store
	router  *gin.Engine
	tokener token.Tokener
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(store db.Store) (*Server, error) {
	tokener, err := token.NewJWTToken("12345678901234567890123456789012")
	if err != nil {
		return nil, fmt.Errorf("failed to create tokener: %w", err)
	}

	server := &Server{
		store:   store,
		tokener: tokener,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/login", server.loginUser)
	router.GET("/users", server.getUsers)
	router.POST("/users", server.createUser)
	//sunday work
	// router.Group("/auth", server.ValidateToken())
	// {
	// router.Use(server.ValidateToken())
	router.GET("/authUser", server.ValidateToken(), server.getUsers)
	// }
	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
