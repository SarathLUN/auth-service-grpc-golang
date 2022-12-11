package main

import (
	"context"
	"fmt"
	"github.com/SarathLUN/auth-service-grpc-golang/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
)

// ? create variable that we'll re-assign later
var (
	server      *gin.Engine
	ctx         context.Context
	mongoClient *mongo.Client
	redisClient *redis.Client
)

// init function that will run before `main` function
func init() {
	// load the .env variable
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalln("Could not load environment variables", err)
	}

	// create context
	ctx = context.TODO()

	// connect to MongoDB
	mongoConn := options.Client().ApplyURI(config.DBUri)
	mongoClient, err := mongo.Connect(ctx, mongoConn)
	if err != nil {
		log.Fatalln("Could not connect with MongoDB: ", err)
	}
	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalln("can't ping MongoDB: ", err)
	}
	fmt.Println("MongoDB successful connected...")

	// Connect to Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Fatalln("can't ping Redis: ", err)
	}
	err = redisClient.Set(ctx, "test", "Welcome to Golang with Redis and MongoDB",
		0).Err()
	if err != nil {
		log.Fatalln("can't connect with Redis: ", err)
	}
	fmt.Println("Redis client connect successful...")

	// Create the Gin Engine instant
	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalln("can't load config: ", err)
	}
	defer mongoClient.Disconnect(ctx)

	value, err := redisClient.Get(ctx, "test").Result()
	if err == redis.Nil {
		fmt.Println("key: test doesn't exist", err)
	} else if err != nil {
		log.Fatalln(err)
	}

	router := server.Group("/api")
	router.GET("/health-checker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": value,
		})
	})
	log.Fatalln(server.Run(":" + config.Port))
}
