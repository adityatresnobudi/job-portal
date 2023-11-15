package router

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adityatresnobudi/job-portal/db"
	"github.com/adityatresnobudi/job-portal/handler"
	"github.com/adityatresnobudi/job-portal/logger"
	"github.com/adityatresnobudi/job-portal/middleware"
	"github.com/adityatresnobudi/job-portal/repository"
	"github.com/adityatresnobudi/job-portal/usecase"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func NewRouter(h *handler.Handler) *gin.Engine {
	router := gin.Default()
	router.ContextWithFallback = true

	router.Use(requestid.New())
	router.Use(middleware.Logger(logger.NewLogger()))
	router.Use(middleware.GlobalErrorMiddleware())

	job := router.Group("/jobs", middleware.WithTimeout())
	job.GET("", h.GetJobs)
	job.POST("", middleware.Auth(), h.CreateNewJobs)
	job.PUT("/:id/close", middleware.Auth(), h.CloseJobs)
	job.PUT("/:id/update", middleware.Auth(), h.ChangeJobs)

	user := router.Group("/auth", middleware.WithTimeout())
	user.POST("/register", h.CreateUser)
	user.POST("/login", h.LoginUser)

	userJob := router.Group("/users", middleware.WithTimeout())
	userJob.POST("/apply", middleware.Auth(), h.ApplyJob)

	return router
}

func Serve() {
	db, err := db.Connect()
	if err != nil {
		log.Println(err)
	}

	jr := repository.NewJobRepository(db)
	ju := usecase.NewJobUsecase(jr)

	ur := repository.NewUserRepository(db)
	uu := usecase.NewUserUsecase(ur)

	ujr := repository.NewUserJobRepository(db)
	uju := usecase.NewUserJobUsecase(ujr)

	h := handler.NewHandler(ju, uu, uju)
	router := NewRouter(h)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
	log.Println("timeout of 5 seconds.")

	log.Println("Server exiting")
}
