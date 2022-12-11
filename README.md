# auth-service-grpc-golang

## 1. API with Golang + MongoDB + Redis + Gin Gonic: Project Setup

### How to Setup Golang with MongoDB and Redis

#### Creating MongoDB and Redis Database with Docker-compose

The most straightforward way to get MongoDB and Redis database instances running on your machine is to use Docker and
Docker-compose. I assume you already have Docker and Docker-compose installed on your computer.

To initialize a new Golang project, create a folder and open it with your preferred text editor. I recommend Visual
Studio Code (VS Code) because it’s now the default text editor for new and professional developers.

Create a `docker-compose.yml` file in your root directory and paste the configurations below into it.

**docker-compose.yml**

```yaml
version: '3'
services:
  mongodb:
    image: mongo
    container_name: mongodb
    restart: always
    env_file:
      - ./app.env

    ports:
      - '6000:27017'
    volumes:
      - mongodb:/data/db

  redis:
    image: redis:alpine
    container_name: redis
    ports:
      - '6379:6379'
    volumes:
      - redisDB:/data
volumes:
  mongodb:
  redisDB:
```

To provide the credentials needed by the MongoDB Docker image, we need to create a `app.env` file in the root directory.

You can add the `app.env` file to your `.gitignore` file to omit it from the commits.

**app.env**

```dotenv
PORT=8000
MONGO_INITDB_ROOT_USERNAME=root
MONGO_INITDB_ROOT_PASSWORD=password123
MONGODB_LOCAL_URI=mongodb://root:password123@localhost:6000
REDIS_URL=localhost:6379

```

Once we have the configurations in place we need to run the docker containers.

```shell
docker-compose up -d
```

#### Setup Environment Variables

The most essential part of our Golang application is to set up environment variables. This will allow us to store
sensitive information and we can easily exclude them from being committed to a repository.

We are going to use the popular library `viper` package to load and validate the environment variables.

Before we can start installing third-party packages we need to initialize the Golang application. Open the integrated
terminal/command-line and run this command to initialize a new Golang project.

```shell
go mod init github.com/example/golang-test
```

Where `example` is your profile name on GitHub and `golang-test` is the name of your Golang project. Once the Golang
project has been initialized, run this command to install the [Golang viper package](https://github.com/spf13/viper).

```shell
go get github.com/spf13/viper
```

Next, create a `config` folder in the root directory and within this config folder create a `default.go` file. We need
to define a struct that will contain all the variables in the `app.env` file and also provide the appropriate value
types.

**config/default.go**

```go
type Config struct {
DBUri    string `mapstructure:"MONGODB_LOCAL_URI"`
RedisUri string `mapstructure:"REDIS_URL"`
Port     string `mapstructure:"PORT"`
}

```

The Golang viper package uses github.com/mitchellh/mapstructure under the hood for unmarshaling values which
uses `mapstructure` tags by default. You need to provide the correct variable name to the `mapstructure` tag.

Next, let’s define a function to load and read the content of the `app.env` file. The `LoadConfig` function will take
the path to the `app.env` file as a parameter and return the `config` instance and a possible error if it exists.

**config/default.go**

```go
package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBUri    string `mapstructure:"MONGODB_LOCAL_URI"`
	RedisUri string `mapstructure:"REDIS_URL"`
	Port     string `mapstructure:"PORT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
```

#### How to Connect Golang App to Redis and MongoDB

In this project and the upcoming series, I’ll use Gin Gonic as my web framework since it has a performance that is 40
times faster compared to other web frameworks.

Also, we’ll be using the [MongoDB Golang driver](https://www.mongodb.com/docs/drivers/go/current/) to connect, access,
mutate and perform aggregations on the MongoDB database running in the Docker container.

Run the command below to install the [MongoDB Go driver](https://www.mongodb.com/docs/drivers/go/current/)
, [Go Redis](https://github.com/go-redis/redis), and [Gin Gonic](https://github.com/gin-gonic/gin) packages.

```shell
go get github.com/go-redis/redis/v8 go.mongodb.org/mongo-driver/mongo github.com/gin-gonic/gin

```

Now create a `main.go` file in the root directory and add these variables to be used later.

Don’t worry about the imports because IntelliJ IDEA / VS Code will automatically import them once you save the file
after using the packages.

```go
package main

// ? Require the packages
import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/example/golang-test/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ? Create required variables that we'll re-assign later
var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client
	redisclient *redis.Client
)

```

Next, let’s create an init function that will run before the main function. We first need to evoke
the `config.LoadConfig(".")` function by passing the path to the `app.env` file as an argument.

**main.go**

```go
// init function that will run before `main` function
func init() {
    // load the .env variable
    config, err := config.LoadConfig(".")
    if err != nil {
        log.Fatalln("Could not load environment variables", err)
    }

```

Once the environment variables have been loaded and read by Viper, we can now connect to the MongoDB database instance.

**main.go**

```go
package main

// ? Require the packages

// ? Create required variables that we'll re-assign later

// ? Init function that will run before the "main" function
func init() {

	// ? Load the .env variables

	// ? Create a context
	ctx = context.TODO()

	// ? Connect to MongoDB
	mongoconn := options.Client().ApplyURI(config.DBUri)
	mongoclient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")
}
```

Now, let’s connect to the Redis database. Also, let’s create an instance of the Gin Engine after the Redis database has also been connected.

**main.go**

```go
package main

// ? Require the packages

// ? Create required variables that we'll re-assign later

// ? Init function that will run before the "main" function
     func init() {

	// ? Load the .env variables

	// ? Create a context
	ctx = context.TODO()

	// ? Connect to MongoDB

	// ? Connect to Redis
	redisclient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})

	if _, err := redisclient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	err = redisclient.Set(ctx, "test", "Welcome to Golang with Redis and MongoDB", 
     0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis client connected successfully...")

	// ? Create the Gin Engine instance
	server = gin.Default()
   }
```

Here comes the good part. Let’s define the `main` function.

**main.go**

```go
func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoclient.Disconnect(ctx)

	value, err := redisclient.Get(ctx, "test").Result()

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		panic(err)
	}

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})

	log.Fatal(server.Run(":" + config.Port))
}
```
Here is a breakdown of what I did above:
- First I loaded the environment variables with the `config.LoadConfig(".")` function.
- Next, I created a router group and passed `/api` as an argument since all requests to the server will have `/api` after the hostname.
- Lastly, I created a GET route `/health-checker` to return the message we stored in the Redis database. This is only needed for testing the API.

Below is the complete code for the `main.go` file.

**main.go**

```go
package main

// ? Require the packages
import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/example/golang-test/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ? Create required variables that we'll re-assign later
var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client
	redisclient *redis.Client
)

// ? Init function that will run before the "main" function
func init() {

	// ? Load the .env variables
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	// ? Create a context
	ctx = context.TODO()

	// ? Connect to MongoDB
	mongoconn := options.Client().ApplyURI(config.DBUri)
	mongoclient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	// ? Connect to Redis
	redisclient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})

	if _, err := redisclient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	err = redisclient.Set(ctx, "test", "Welcome to Golang with Redis and MongoDB", 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis client connected successfully...")

	// ? Create the Gin Engine instance
	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoclient.Disconnect(ctx)

	value, err := redisclient.Get(ctx, "test").Result()

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		panic(err)
	}

	router := server.Group("/api")
	router.GET("/health-checker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})

	log.Fatal(server.Run(":" + config.Port))
}
```
#### Test the Golang API

Now, let’s test the API by sending a GET request to [http://localhost:8000/api/health-checker](http://localhost:8000/api/health-checker).

Before you make the API request, run this command to start both the Docker containers and the Go application.

```shell
docker-compose up -d && go run main.go

```

Now open an API testing tool (Postman, Thunder Client, HTTP Client, REST Client, etc ) and send a GET request to 
[http://localhost:8000/api/health-checker](http://localhost:8000/api/health-checker)

![auth-service-grpc-golang-health-checker](images/auth-service-grpc-golang-health-checker.png)

show redis key:

![](images/redis-cli-01.png)
