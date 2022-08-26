package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	. "github.com/kadirgonen/movie-api/app/handlers"
	"github.com/kadirgonen/movie-api/app/pkg/config"
	. "github.com/kadirgonen/movie-api/app/pkg/db"
	. "github.com/kadirgonen/movie-api/app/pkg/graceful"
	logger "github.com/kadirgonen/movie-api/app/pkg/log"
	. "github.com/kadirgonen/movie-api/app/pkg/middleware"
	. "github.com/kadirgonen/movie-api/app/repo"
	. "github.com/kadirgonen/movie-api/app/service"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Loadconfig Failed: %v", err)
	}

	logger.NewLogger(cfg)
	defer logger.Close()

	errload := godotenv.Load("../.env")
	if errload != nil {
		zap.L().Fatal("Error loading .env file")
	}

	// It is possible to integrate different db technologies
	base := DBBase{DbType: &POSTGRES{}}
	db, err := base.DbType.Create(cfg)
	if err != nil {
		zap.L().Fatal("DB cannot init", zap.Error(err))
	}

	gin.SetMode(gin.ReleaseMode)
	g := gin.Default()
	LoggerMiddleware(g)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerConfig.Port),
		Handler:      g,
		ReadTimeout:  time.Duration(cfg.ServerConfig.ReadTimeoutSecs * int64(time.Second)),
		WriteTimeout: time.Duration(cfg.ServerConfig.WriteTimeoutSecs * int64(time.Second)),
	}

	getUp(g, db, cfg)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("main.listen and serve: ", zap.Error(err))
		}
	}()
	zap.L().Debug(os.Getenv("START_SERVER"))

	// HealthCheck On DB
	go healthCheck()

	ShutdownGin(srv, time.Duration(cfg.ServerConfig.TimeoutSecs*int64(time.Second)))

}

func getUp(g *gin.Engine, db *gorm.DB, cfg *config.Config) {
	rootRouter := g.Group(cfg.ServerConfig.RoutePrefix)
	userRooter := rootRouter.Group("/user")
	movieRooter := rootRouter.Group("/movie")

	userRepo := NewUserRepository(db)
	userService := NewUserService(userRepo)
	NewUserHandler(userRooter, userService, cfg)

	movieRepo := NewMovieRepository(db)
	movieService := NewMovieService(movieRepo)
	NewMovieHandler(movieRooter, movieService, cfg)
}

// HealthCheck checks the db is ready with 10 seconds break
func healthCheck() {
	tck := time.NewTicker(10 * time.Second)
	issueOnHealth := make(chan bool)
	go func() {
		for {
			select {
			case <-issueOnHealth:
				return
			case time := <-tck.C:
				zap.L().Debug("Time", zap.Reflect("time:", time))
				resp, err := http.Get("http://127.0.0.1:8080/api/v1/movie-api/status")
				if err != nil {
					zap.L().Fatal("DB doesn't work", zap.Error(err))
					issueOnHealth <- true
				}
				zap.L().Debug("Response", zap.String("resp:", resp.Status))

			}
		}
	}()
	<-issueOnHealth
}
