package main

import (
	"context"
	"os"
	"time"

	"github.com/Agmer17/srd_lab_creative/internal/bootstrap"
	"github.com/Agmer17/srd_lab_creative/pkg"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("COULDN'T FIND OR READ THE ENV : " + err.Error())
	}

	googleClientId := os.Getenv("GOOGLE_OAUTH_CLIENT")
	googleClientSecret := os.Getenv("GOOGLE_OAUTH_SECRET")
	databaseUrl := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	redisUrl := os.Getenv("REDIS_URL")

	pkg.JwtInit(jwtSecret)

	mainAppCtx := context.Background()

	// database pool and connection test!
	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		panic("COULDN'T SETUP THE DATABASE : " + err.Error())
	}
	config.MaxConns = 15
	config.MinConns = 3
	config.MaxConnIdleTime = 20 * time.Minute
	config.MaxConnLifetime = 10 * time.Minute

	pool, err := pgxpool.NewWithConfig(mainAppCtx, config)
	if err != nil {
		panic("COULDN'T SETUP THE DATABASE : " + err.Error())
	}

	err = pool.Ping(mainAppCtx)
	if err != nil {
		panic("DB CONNECTION FAILED : " + err.Error())
	}
	// =================================

	// redis setup
	// setup redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "",
		DB:       0,
	})
	_, err = rdb.Ping(mainAppCtx).Result()
	if err != nil {
		panic(err)
	}
	// =========
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	app := bootstrap.NewApp(mainAppCtx, r, googleClientId, googleClientSecret, pool, rdb)
	app.Run()

}
