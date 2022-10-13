package main

//require package
import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/isanaID/mongo-redis-golang-gin/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client
	redisclient *redis.Client
)

func init() {
	//load env
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load env", err)
	}

	//create context

	ctx = context.TODO()

	//connect to mongo
	mongoconn := options.Client().ApplyURI(config.DBUri)
	mongoclient, err = mongo.Connect(ctx, mongoconn)

	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Connected to MongoDB!")

	//connect to redis
	redisclient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})

	if _, err := redisclient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	err = redisclient.Set(ctx, "test", "Welcome to golang with redis and mongodb", 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Redis!")

	//create the gin engine instance
	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load env", err)
	}

	defer mongoclient.Disconnect(ctx)

	value, err := redisclient.Get(ctx, "test").Result()

	if err == redis.Nil {
		fmt.Println("key does not exist")
	} else if err != nil {
		panic(err)
	}

	router := server.Group("/api")
	router.GET("healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": value,
		})
	})

	log.Fatal(server.Run(":" + config.Port))
}
