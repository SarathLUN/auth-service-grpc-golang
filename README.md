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

Now, let’s connect to the Redis database. Also, let’s create an instance of the Gin Engine after the Redis database has
also been connected.

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
router.GET("/healthchecker", func (ctx *gin.Context) {
ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
})

log.Fatal(server.Run(":" + config.Port))
}
```

Here is a breakdown of what I did above:

- First I loaded the environment variables with the `config.LoadConfig(".")` function.
- Next, I created a router group and passed `/api` as an argument since all requests to the server will have `/api`
  after the hostname.
- Lastly, I created a GET route `/health-checker` to return the message we stored in the Redis database. This is only
  needed for testing the API.

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

Now, let’s test the API by sending a GET request
to [http://localhost:8000/api/health-checker](http://localhost:8000/api/health-checker).

Before you make the API request, run this command to start both the Docker containers and the Go application.

```shell
docker-compose up -d && go run main.go

```

Now open an API testing tool (Postman, Thunder Client, HTTP Client, REST Client, etc ) and send a GET request to
[http://localhost:8000/api/health-checker](http://localhost:8000/api/health-checker)

![auth-service-grpc-golang-health-checker](images/auth-service-grpc-golang-health-checker.png)

show redis key:

![](images/redis-cli-01.png)

## 2. Golang & MongoDB: JWT Authentication and Authorization

In this article, you’ll learn how to implement RS256 JWT Authentication and Authorization with Golang, Gin Gonic,
MongoDB-Go-driver, and Docker-compose.

### How to Generate Public and Private Keys

To generate the public and private keys, navigate to [this website](https://www.base64encode.org/) and click on the “**
Generate New Keys**” button.

**Private Key**

```text
-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgHxZA63YHZayYkq/S4FgcLZs2VpgdX/S2EI2HtlxXLkg/lqB2WcX
tljFC5zD6DbaZ+DaYwGTAooIfCMH/8ZwLRgdvRouGvEYbAx7Cbe5393F62LYR0QU
nsX8qUufyeuAeL90QEERADrqfS5IIKQxhF5WwMqflOuWmvObM+UdZyN/AgMBAAEC
gYA05BfwhhBv3Qh6CswvB7ppnfsFuinAaQqbfzG4A8QQ6rvF+51hgHCLaaT61Umc
q2Gr/rhGHyXhWKKUs/5Thbtn+vEWjcySjQs1lzpy5EMv8v1qe9a2JDIKAfdLi8u8
IFrLfYbWHHceyDrL0JPime5CnNt1JpXDl+Uv7xB8Xp6jIQJBAPX88ACrXvqD1EcU
89yV4Z9sn/DHUuN5D57vZd+P1cyuT2HMjNIL6Hi6oDgHN3LzsKpUwcpkUSBWC0yp
XlddeGUCQQCBaKZ+9PL4RnR/eyfd0igHEzC6TkvTN29JGzqdckB/KSPFutMb5hN8
8wsl9Xg5ugk5GDm+bLeTr58hUUzHVSQTAkEA3guCQbhvtcbGGZvqPRX2YkhVX60t
PzSKSzdiSBm4LliyQ1+3ReuF32TnN2wcO2GJ+VckzyMasAu7i2nfRRGUmQJAOiFH
QRpJFkuXwzgQNW1Uu7DcH6RSHUP2gnFvyA/tdvGE9yov8IVOVb9bT6bPZabyw6nl
LmvpnrayFaZOhe/x7wJAG2F6ayLqxoxwJjeocHD8PIJy7ZhUOUzRbDlh/Z3kI+Tx
7viBy6Wl4EeUG36tT/1yejIKEz3eYGHZEkTnmEFBOw==
-----END RSA PRIVATE KEY-----
```

**Public Key**

```text
-----BEGIN PUBLIC KEY-----
MIGeMA0GCSqGSIb3DQEBAQUAA4GMADCBiAKBgHxZA63YHZayYkq/S4FgcLZs2Vpg
dX/S2EI2HtlxXLkg/lqB2WcXtljFC5zD6DbaZ+DaYwGTAooIfCMH/8ZwLRgdvRou
GvEYbAx7Cbe5393F62LYR0QUnsX8qUufyeuAeL90QEERADrqfS5IIKQxhF5WwMqf
lOuWmvObM+UdZyN/AgMBAAE=
-----END PUBLIC KEY-----
```

Next, copy the private key and visit [this website](https://www.base64encode.org/) to encode it in base64. Later we’ll
use the base64 Golang package to decode it back to ASCII string.

On the [Base64 Decode and Encode website](https://www.base64encode.org/), paste the copied private key into the input
field and click on the “Encode” button.

```text
LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlDV3dJQkFBS0JnSHhaQTYzWUhaYXlZa3EvUzRGZ2NMWnMyVnBnZFgvUzJFSTJIdGx4WExrZy9scUIyV2NYCnRsakZDNXpENkRiYVorRGFZd0dUQW9vSWZDTUgvOFp3TFJnZHZSb3VHdkVZYkF4N0NiZTUzOTNGNjJMWVIwUVUKbnNYOHFVdWZ5ZXVBZUw5MFFFRVJBRHJxZlM1SUlLUXhoRjVXd01xZmxPdVdtdk9iTStVZFp5Ti9BZ01CQUFFQwpnWUEwNUJmd2hoQnYzUWg2Q3N3dkI3cHBuZnNGdWluQWFRcWJmekc0QThRUTZydkYrNTFoZ0hDTGFhVDYxVW1jCnEyR3IvcmhHSHlYaFdLS1VzLzVUaGJ0bit2RVdqY3lTalFzMWx6cHk1RU12OHYxcWU5YTJKRElLQWZkTGk4dTgKSUZyTGZZYldISGNleURyTDBKUGltZTVDbk50MUpwWERsK1V2N3hCOFhwNmpJUUpCQVBYODhBQ3JYdnFEMUVjVQo4OXlWNFo5c24vREhVdU41RDU3dlpkK1AxY3l1VDJITWpOSUw2SGk2b0RnSE4zTHpzS3BVd2Nwa1VTQldDMHlwClhsZGRlR1VDUVFDQmFLWis5UEw0Um5SL2V5ZmQwaWdIRXpDNlRrdlROMjlKR3pxZGNrQi9LU1BGdXRNYjVoTjgKOHdzbDlYZzV1Z2s1R0RtK2JMZVRyNThoVVV6SFZTUVRBa0VBM2d1Q1FiaHZ0Y2JHR1p2cVBSWDJZa2hWWDYwdApQelNLU3pkaVNCbTRMbGl5UTErM1JldUYzMlRuTjJ3Y08yR0orVmNrenlNYXNBdTdpMm5mUlJHVW1RSkFPaUZIClFScEpGa3VYd3pnUU5XMVV1N0RjSDZSU0hVUDJnbkZ2eUEvdGR2R0U5eW92OElWT1ZiOWJUNmJQWmFieXc2bmwKTG12cG5yYXlGYVpPaGUveDd3SkFHMkY2YXlMcXhveHdKamVvY0hEOFBJSnk3WmhVT1V6UmJEbGgvWjNrSStUeAo3dmlCeTZXbDRFZVVHMzZ0VC8xeWVqSUtFejNlWUdIWkVrVG5tRUZCT3c9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQ==
```

Copy the encoded private key and add it to the `app.env` file as `ACCESS_TOKEN_PRIVATE_KEY`.

Go back to the website where you generated the keys and copy the corresponding public key. Visit the **Base64 Decode and
Encode website** and encode it to base64 before adding it to the `app.env` file as `ACCESS_TOKEN_PUBLIC_KEY`.

```text
LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlHZU1BMEdDU3FHU0liM0RRRUJBUVVBQTRHTUFEQ0JpQUtCZ0h4WkE2M1lIWmF5WWtxL1M0RmdjTFpzMlZwZwpkWC9TMkVJMkh0bHhYTGtnL2xxQjJXY1h0bGpGQzV6RDZEYmFaK0RhWXdHVEFvb0lmQ01ILzhad0xSZ2R2Um91Ckd2RVliQXg3Q2JlNTM5M0Y2MkxZUjBRVW5zWDhxVXVmeWV1QWVMOTBRRUVSQURycWZTNUlJS1F4aEY1V3dNcWYKbE91V212T2JNK1VkWnlOL0FnTUJBQUU9Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ==
```

Repeat the process for the refresh token.

**Private Key**

```text
-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQCWHTGTCXQcWAJlU4YozYt4hGhcRDqtdbCUp/bfsQ7pqqcgPubd
fpZxQU4sLXr1CzPRc0qyuB1Gx7Bb/nBAu3YWdRo1jIHqSq6fILeIo6FtaIzJWdjD
TXnzGRGm3CjzvQoc/nxgl928Usi+9CIC/BNQJnlRfLYdyFjJDLu2uXctWwIDAQAB
AoGAFbY9QF4irCLSneM5yVcJcDH+gjHgJsyZ8LZmj64lUQi8yX4vmCPCELl4XA1r
9ltbcFTznevhWCZzzKp9g+9ZQeYp6VVC038tmeP9B2klu6DD6cRvv2HvaCkTBOcG
cTmvVfZ6j01cqhGGa63bYseksXvTUlAlNAyflzcC2iQvi8ECQQD05Q2q32Yy8O4j
jDgHfPb64HYcnV3AHDITRbHOjtbArvcvEu33rydoIbpXkvjqxhR2tabNJfmXzOLw
maGwHMV7AkEAnOvYXi/nj0WI+b3RRE/C0WVIPOTe9SIQ17DXv8/1Y9T8cLzuua1z
BIMQ61YsGH/y8UnqRltIhNF5M34QatSBoQJARxwzycgwzt3y6FOQ+Wup89M1uL1h
OaV2TO5R4R7M8QpdMEgY3XD4yl3w0/kTWmjGDGpW0v8R6KhPxXl238Ka/QJAMV4Z
YfhcNzy/DSLklUMqnUnp4g52/lPLTLA8EAJSxyeYi7dEu3GP/9lkF0Wy9bo7KiIA
7EWHmRF1ufdtHs7gQQJAGT13JLj/GVUhCpyZ86vEScwppSOTIBsbFIiQO8CW9C45
fsX+wZ2tqzEc70id4S0ysg7YYEOELwQiyavZAOB15Q==
-----END RSA PRIVATE KEY-----
```

**Public Key**

```text
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCWHTGTCXQcWAJlU4YozYt4hGhc
RDqtdbCUp/bfsQ7pqqcgPubdfpZxQU4sLXr1CzPRc0qyuB1Gx7Bb/nBAu3YWdRo1
jIHqSq6fILeIo6FtaIzJWdjDTXnzGRGm3CjzvQoc/nxgl928Usi+9CIC/BNQJnlR
fLYdyFjJDLu2uXctWwIDAQAB
-----END PUBLIC KEY-----
```

**base64 of private key**

```text
LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlDV3dJQkFBS0JnUUNXSFRHVENYUWNXQUpsVTRZb3pZdDRoR2hjUkRxdGRiQ1VwL2Jmc1E3cHFxY2dQdWJkCmZwWnhRVTRzTFhyMUN6UFJjMHF5dUIxR3g3QmIvbkJBdTNZV2RSbzFqSUhxU3E2ZklMZUlvNkZ0YUl6SldkakQKVFhuekdSR20zQ2p6dlFvYy9ueGdsOTI4VXNpKzlDSUMvQk5RSm5sUmZMWWR5RmpKREx1MnVYY3RXd0lEQVFBQgpBb0dBRmJZOVFGNGlyQ0xTbmVNNXlWY0pjREgrZ2pIZ0pzeVo4TFptajY0bFVRaTh5WDR2bUNQQ0VMbDRYQTFyCjlsdGJjRlR6bmV2aFdDWnp6S3A5Zys5WlFlWXA2VlZDMDM4dG1lUDlCMmtsdTZERDZjUnZ2Mkh2YUNrVEJPY0cKY1RtdlZmWjZqMDFjcWhHR2E2M2JZc2Vrc1h2VFVsQWxOQXlmbHpjQzJpUXZpOEVDUVFEMDVRMnEzMll5OE80agpqRGdIZlBiNjRIWWNuVjNBSERJVFJiSE9qdGJBcnZjdkV1MzNyeWRvSWJwWGt2anF4aFIydGFiTkpmbVh6T0x3Cm1hR3dITVY3QWtFQW5PdllYaS9uajBXSStiM1JSRS9DMFdWSVBPVGU5U0lRMTdEWHY4LzFZOVQ4Y0x6dXVhMXoKQklNUTYxWXNHSC95OFVucVJsdEloTkY1TTM0UWF0U0JvUUpBUnh3enljZ3d6dDN5NkZPUStXdXA4OU0xdUwxaApPYVYyVE81UjRSN004UXBkTUVnWTNYRDR5bDN3MC9rVFdtakdER3BXMHY4UjZLaFB4WGwyMzhLYS9RSkFNVjRaCllmaGNOenkvRFNMa2xVTXFuVW5wNGc1Mi9sUExUTEE4RUFKU3h5ZVlpN2RFdTNHUC85bGtGMFd5OWJvN0tpSUEKN0VXSG1SRjF1ZmR0SHM3Z1FRSkFHVDEzSkxqL0dWVWhDcHlaODZ2RVNjd3BwU09USUJzYkZJaVFPOENXOUM0NQpmc1grd1oydHF6RWM3MGlkNFMweXNnN1lZRU9FTHdRaXlhdlpBT0IxNVE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQ==
```

**base64 of public key**

```text
LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlHZk1BMEdDU3FHU0liM0RRRUJBUVVBQTRHTkFEQ0JpUUtCZ1FDV0hUR1RDWFFjV0FKbFU0WW96WXQ0aEdoYwpSRHF0ZGJDVXAvYmZzUTdwcXFjZ1B1YmRmcFp4UVU0c0xYcjFDelBSYzBxeXVCMUd4N0JiL25CQXUzWVdkUm8xCmpJSHFTcTZmSUxlSW82RnRhSXpKV2RqRFRYbnpHUkdtM0NqenZRb2MvbnhnbDkyOFVzaSs5Q0lDL0JOUUpubFIKZkxZZHlGakpETHUydVhjdFd3SURBUUFCCi0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo=
```

### Update Environment Variables

You can copy and paste the environment variables below into your `app.env` file if you find it difficult to encode the
private and public keys.

I also added the access and refresh tokens expiration time and max-age to the `app.env` file.

**app.env**

```dotenv
PORT=8000

MONGO_INITDB_ROOT_USERNAME=root
MONGO_INITDB_ROOT_PASSWORD=password123

MONGODB_LOCAL_URI=mongodb://root:password123@localhost:6000

REDIS_URL=localhost:6379

ACCESS_TOKEN_PRIVATE_KEY=LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlDV3dJQkFBS0JnSHhaQTYzWUhaYXlZa3EvUzRGZ2NMWnMyVnBnZFgvUzJFSTJIdGx4WExrZy9scUIyV2NYCnRsakZDNXpENkRiYVorRGFZd0dUQW9vSWZDTUgvOFp3TFJnZHZSb3VHdkVZYkF4N0NiZTUzOTNGNjJMWVIwUVUKbnNYOHFVdWZ5ZXVBZUw5MFFFRVJBRHJxZlM1SUlLUXhoRjVXd01xZmxPdVdtdk9iTStVZFp5Ti9BZ01CQUFFQwpnWUEwNUJmd2hoQnYzUWg2Q3N3dkI3cHBuZnNGdWluQWFRcWJmekc0QThRUTZydkYrNTFoZ0hDTGFhVDYxVW1jCnEyR3IvcmhHSHlYaFdLS1VzLzVUaGJ0bit2RVdqY3lTalFzMWx6cHk1RU12OHYxcWU5YTJKRElLQWZkTGk4dTgKSUZyTGZZYldISGNleURyTDBKUGltZTVDbk50MUpwWERsK1V2N3hCOFhwNmpJUUpCQVBYODhBQ3JYdnFEMUVjVQo4OXlWNFo5c24vREhVdU41RDU3dlpkK1AxY3l1VDJITWpOSUw2SGk2b0RnSE4zTHpzS3BVd2Nwa1VTQldDMHlwClhsZGRlR1VDUVFDQmFLWis5UEw0Um5SL2V5ZmQwaWdIRXpDNlRrdlROMjlKR3pxZGNrQi9LU1BGdXRNYjVoTjgKOHdzbDlYZzV1Z2s1R0RtK2JMZVRyNThoVVV6SFZTUVRBa0VBM2d1Q1FiaHZ0Y2JHR1p2cVBSWDJZa2hWWDYwdApQelNLU3pkaVNCbTRMbGl5UTErM1JldUYzMlRuTjJ3Y08yR0orVmNrenlNYXNBdTdpMm5mUlJHVW1RSkFPaUZIClFScEpGa3VYd3pnUU5XMVV1N0RjSDZSU0hVUDJnbkZ2eUEvdGR2R0U5eW92OElWT1ZiOWJUNmJQWmFieXc2bmwKTG12cG5yYXlGYVpPaGUveDd3SkFHMkY2YXlMcXhveHdKamVvY0hEOFBJSnk3WmhVT1V6UmJEbGgvWjNrSStUeAo3dmlCeTZXbDRFZVVHMzZ0VC8xeWVqSUtFejNlWUdIWkVrVG5tRUZCT3c9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQ==


ACCESS_TOKEN_PUBLIC_KEY=LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlHZU1BMEdDU3FHU0liM0RRRUJBUVVBQTRHTUFEQ0JpQUtCZ0h4WkE2M1lIWmF5WWtxL1M0RmdjTFpzMlZwZwpkWC9TMkVJMkh0bHhYTGtnL2xxQjJXY1h0bGpGQzV6RDZEYmFaK0RhWXdHVEFvb0lmQ01ILzhad0xSZ2R2Um91Ckd2RVliQXg3Q2JlNTM5M0Y2MkxZUjBRVW5zWDhxVXVmeWV1QWVMOTBRRUVSQURycWZTNUlJS1F4aEY1V3dNcWYKbE91V212T2JNK1VkWnlOL0FnTUJBQUU9Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ
ACCESS_TOKEN_EXPIRED_IN=15m
ACCESS_TOKEN_MAXAGE=15


REFRESH_TOKEN_PRIVATE_KEY=LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlDV3dJQkFBS0JnUUNXSFRHVENYUWNXQUpsVTRZb3pZdDRoR2hjUkRxdGRiQ1VwL2Jmc1E3cHFxY2dQdWJkCmZwWnhRVTRzTFhyMUN6UFJjMHF5dUIxR3g3QmIvbkJBdTNZV2RSbzFqSUhxU3E2ZklMZUlvNkZ0YUl6SldkakQKVFhuekdSR20zQ2p6dlFvYy9ueGdsOTI4VXNpKzlDSUMvQk5RSm5sUmZMWWR5RmpKREx1MnVYY3RXd0lEQVFBQgpBb0dBRmJZOVFGNGlyQ0xTbmVNNXlWY0pjREgrZ2pIZ0pzeVo4TFptajY0bFVRaTh5WDR2bUNQQ0VMbDRYQTFyCjlsdGJjRlR6bmV2aFdDWnp6S3A5Zys5WlFlWXA2VlZDMDM4dG1lUDlCMmtsdTZERDZjUnZ2Mkh2YUNrVEJPY0cKY1RtdlZmWjZqMDFjcWhHR2E2M2JZc2Vrc1h2VFVsQWxOQXlmbHpjQzJpUXZpOEVDUVFEMDVRMnEzMll5OE80agpqRGdIZlBiNjRIWWNuVjNBSERJVFJiSE9qdGJBcnZjdkV1MzNyeWRvSWJwWGt2anF4aFIydGFiTkpmbVh6T0x3Cm1hR3dITVY3QWtFQW5PdllYaS9uajBXSStiM1JSRS9DMFdWSVBPVGU5U0lRMTdEWHY4LzFZOVQ4Y0x6dXVhMXoKQklNUTYxWXNHSC95OFVucVJsdEloTkY1TTM0UWF0U0JvUUpBUnh3enljZ3d6dDN5NkZPUStXdXA4OU0xdUwxaApPYVYyVE81UjRSN004UXBkTUVnWTNYRDR5bDN3MC9rVFdtakdER3BXMHY4UjZLaFB4WGwyMzhLYS9RSkFNVjRaCllmaGNOenkvRFNMa2xVTXFuVW5wNGc1Mi9sUExUTEE4RUFKU3h5ZVlpN2RFdTNHUC85bGtGMFd5OWJvN0tpSUEKN0VXSG1SRjF1ZmR0SHM3Z1FRSkFHVDEzSkxqL0dWVWhDcHlaODZ2RVNjd3BwU09USUJzYkZJaVFPOENXOUM0NQpmc1grd1oydHF6RWM3MGlkNFMweXNnN1lZRU9FTHdRaXlhdlpBT0IxNVE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQ==


REFRESH_TOKEN_PUBLIC_KEY=LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlHZk1BMEdDU3FHU0liM0RRRUJBUVVBQTRHTkFEQ0JpUUtCZ1FDV0hUR1RDWFFjV0FKbFU0WW96WXQ0aEdoYwpSRHF0ZGJDVXAvYmZzUTdwcXFjZ1B1YmRmcFp4UVU0c0xYcjFDelBSYzBxeXVCMUd4N0JiL25CQXUzWVdkUm8xCmpJSHFTcTZmSUxlSW82RnRhSXpKV2RqRFRYbnpHUkdtM0NqenZRb2MvbnhnbDkyOFVzaSs5Q0lDL0JOUUpubFIKZkxZZHlGakpETHUydVhjdFd3SURBUUFCCi0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo=

REFRESH_TOKEN_EXPIRED_IN=60m
REFRESH_TOKEN_MAXAGE=60
```

Next, update the `config/default.go` file with the environment variables provided in the `app.env` file for Viper to
load and make them available in the project.

**config/default.go**

```go
package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBUri                  string        `mapstructure:"MONGODB_LOCAL_URI"`
	RedisUri               string        `mapstructure:"REDIS_URL"`
	Port                   string        `mapstructure:"PORT"`
	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`
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

### Creating the User models with structs

In the JSON web Token authentication flow, we need to start by signing up the user. To achieve that let’s define a
struct that specifies the fields required to register a user.

To define the user struct, you need to specify the field name followed by the field type and an optional tag which will
be used by `Gin Gonic` and `MongoDB` for validation, marshaling, and unmarshaling.

Now create a `SignUpInput` struct to specify the fields required to register a new user.

The `bson` tag will be used by MongoDB since MongoDB stores data as BSON documents. Also, the `json` tag will be used by
Gin Gonic for validating user inputs.

Lastly, the `binding` tag specifies additional validation rules that will be used by Gin Gonic to validate user inputs.

**models/user.model.go**

```go
package models

import "time"

type SignUpInput struct {
	Name            string    `json:"name" bson:"name" bindinng:"required"`
	Email           string    `json:"email" bson:"email" binding:"required"`
	Password        string    `json:"password" bson:"password" binding:"required,min=8"`
	ConfirmPassword string    `json:"confirm_password" bson:"confirm_password,omitempty" binding:"required"`
	Role            string    `json:"role" bson:"role"`
	Verified        bool      `json:"verified" bson:"verified"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}
```

With the user login logic, we’ll only require the user to provide an email and password. To achieve that create
a `SignInInput` struct with email and password fields and make them required.

**models/user.model.go**

```go
package models

import "time"

// ? SignUpInput struct

type SignInInput struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}
```

Next, create a `DBResponse` struct to define the fields that will be returned by MongoDB.

**models/user.model.go**

```go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ? SignUpInput struct

// ? SignInInput struct

type DBResponse struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Name            string             `json:"name" bson:"name"`
	Email           string             `json:"email" bson:"email"`
	Password        string             `json:"password" bson:"password"`
	PasswordConfirm string             `json:"passwordConfirm,omitempty" bson:"passwordConfirm,omitempty"`
	Role            string             `json:"role" bson:"role"`
	Verified        bool               `json:"verified" bson:"verified"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}
```

Lastly, create a struct to specify the fields that should be included in the JSON response and a function to filter out
the sensitive fields.

```go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ? SignUpInput struct

// ? SignInInput struct

// ? DBResponse struct

type UserResponse struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Role      string             `json:"role,omitempty" bson:"role,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func FilteredResponse(user *DBResponse) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
```

Below is the complete code for the user model

```go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ? SignUpInput struct
type SignUpInput struct {
	Name            string    `json:"name" bson:"name" binding:"required"`
	Email           string    `json:"email" bson:"email" binding:"required"`
	Password        string    `json:"password" bson:"password" binding:"required,min=8"`
	PasswordConfirm string    `json:"passwordConfirm" bson:"passwordConfirm,omitempty" binding:"required"`
	Role            string    `json:"role" bson:"role"`
	Verified        bool      `json:"verified" bson:"verified"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

// ? SignInInput struct
type SignInInput struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

// ? DBResponse struct
type DBResponse struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Name            string             `json:"name" bson:"name"`
	Email           string             `json:"email" bson:"email"`
	Password        string             `json:"password" bson:"password"`
	PasswordConfirm string             `json:"passwordConfirm,omitempty" bson:"passwordConfirm,omitempty"`
	Role            string             `json:"role" bson:"role"`
	Verified        bool               `json:"verified" bson:"verified"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

// ? UserResponse struct
type UserResponse struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Role      string             `json:"role,omitempty" bson:"role,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func FilteredResponse(user *DBResponse) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
```

### Create an Auth and User Interfaces

Typically, in Golang RESTful API, there are controllers and services. Controllers don't have direct access to database,
instead they call services to either mutate or query database.

#### Authentication interface

Now let define an `AuthService` interface that has method `SignUpUser` and `SignInUser`

**services/auth.service.go**

```go
package services

import "github.com/SarathLUN/auth-service-grpc-golang/models"

type AuthService interface {
	SignUpUser(*models.SignUpInput) (*models.DBResponse, error)
	SignInUser(*models.SignInInput) (*models.DBResponse, error)
}
```

#### User Interface

Also, define a `UserService` interface that has method `FindUserById` and `FindUserByEmail`.

**services/user.service.go**

```go
package services

import "github.com/SarathLUN/auth-service-grpc-golang/models"

type UserService interface {
  FindUserById(string) (*models.DBResponse, error)
  FindUserByEmail(string) (*models.DBResponse, error)
}
```

### Create utility function to hash and verify password

Another important port of authentication is to encrypt the plain password provided by the user before saving it to 
the database.

To secure the password, we use a hashing algorithm to transform the plain password into a hashed string.

In the event of a compromised database, the hacker can not easily decrypt the hashed password to get the original 
plain password.

When two strings are hashed the outputs are the same so we use **salt** to ensure that no two users having the same 
password end up with the same hashed password.

There are different ways of hashing a string but in this tutorial, I will use the `Golang Bcrypt package` to hash 
the user's password.

To hash a string with Bcrypt, we specify a **Cost Factor** which is the amount of time needed to calculate a single 
hash.

The higher the **Cost Factor** the longer the hashing time and the more difficult it is to brute force. Modern 
computers now have powerful CPUs so using a Cost Factor of 12 should be fine.

Now, let's define two functions to hash and verify the user's password.

**utils/password.go**

```go
package utils

import (
  "fmt"
  "golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil {
    return "", fmt.Errorf("could not hash password %w", err)
  }
  return string(hashedPassword),nil
}

func VerifyPassword(hashedPassword string, candidatePassword string) error {
  return bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(candidatePassword))
}
```

### Create services that interact with the database

Next, let's define a `SignUpUser` service that will be called by a controller to create a new user in database.

#### Auth interface implementation

I included the MongoDB Collection struct and the context interface of the context package in the `AuthServiceImpl` 
struct to help us perform basic CRUD operation against the MongoDB database.

Then I use the `NewAuthService` constructor function to instantiate the MongoDB collection struct and Context 
interface. Also, the `NewAuthService` constructor implements the `AuthService` interface we defined above.

**services/auth.service.impl.go**

```go

package services

import (
  "context"
  "errors"
  "github.com/SarathLUN/auth-service-grpc-golang/models"
  "github.com/SarathLUN/auth-service-grpc-golang/utils"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "strings"
  "time"
)

type AuthServiceImpl struct {
  collection *mongo.Collection
  ctx context.Context
}

func NewAuthService(collection *mongo.Collection, ctx context.Context) AuthService {
  return &AuthServiceImpl{
    collection: collection,
    ctx: ctx,
  }
}

func (uc *AuthServiceImpl) SignUpUser(user *models.SignUpInput) (*models.DBResponse, error) {
  user.CreatedAt = time.Now()
  user.UpdatedAt = user.CreatedAt
  user.Email = strings.ToLower(user.Email)
  user.ConfirmPassword = ""
  user.Verified = true
  user.Role = "user"
  hashedPassword, _ := utils.HashPassword(user.Password)
  user.Password = hashedPassword
  res, err := uc.collection.InsertOne(uc.ctx, &user)
  if err != nil {
    if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000{
      return nil, errors.New("user with this email already is existed")
    }
    return nil, err
  }

  // create unique index for the email field
  opt := options.Index()
  opt.SetUnique(true)
  index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: opt}

  if _, err := uc.collection.Indexes().CreateOne(uc.ctx, index); err != nil{
    return nil, errors.New("could not create index for email")
  }

  var newUser *models.DBResponse
  query := bson.M{"_id":res.InsertedID}

  err = uc.collection.FindOne(uc.ctx,query).Decode(&newUser)
  if err != nil {
    return nil, err
  }
  return newUser, nil
}

func (uc AuthServiceImpl) SignInUser(*models.SignInInput) (*models.DBResponse, error) {
  return nil,nil
}
```

Here is the summary of what I did in the `SignUpUser` function:

- First, I added to the database with the `InsertOne()` function.
- Next, I added a unique index on the email field to ensure that no two users can have the same email address.
- Lastly, I added the `FindOne()` function to find and return the user that was added to database.

#### User interface implementation 

In order to implement the `UserService` interface, I define the `FindUserByID()` and `FindUserByEmail()` function 
receivers.

In the `FindUserByID()` function, I convert the `id` string to a MongoDB ObjectID for the query to work.

Also, I convert the `email` string to lowercase in the `FindUserByEmail()` function before the calling the `FindOne()
` function, since MongoDB fields are not case-sensitive.

**services/user.service.impl.go**

```go
package services

import (
	"context"
	"github.com/SarathLUN/auth-service-grpc-golang/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type UserServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserServiceImpl(collection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{
		collection: collection,
		ctx:        ctx,
	}
}

func (us *UserServiceImpl) FindUserById(id string) (*models.DBResponse, error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	var user *models.DBResponse
	query := bson.M{"_id": oid}
	err := us.collection.FindOne(us.ctx, query).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}
	return user, nil
}

func (us *UserServiceImpl) FindUserByEmail(email string) (*models.DBResponse, error) {
	var user *models.DBResponse
	query := bson.M{"email": strings.ToLower(email)}
	err := us.collection.FindOne(us.ctx, query).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}
	return user, nil
}

```

### Create a utility function to sign and verify JWT tokens

User authentication can be done with difference strategies but each strategy has a drawback or security flaws.

The concept of authentication can get complex so library like Passport provide difference strategies to simplify the 
process.

In this tutorial I use JSON Web Token to sign in and persist the user to enable them to request protected routes 
without to login on every request.

#### Create Json Web Token

JSON Web Token are used to implement stateless authentication but that alone is not enough since the only way that 
the token can be invalidated is when it has expired.

You can read more about the flaws of JSON Web Token in this [article](https://redis.com/blog/json-web-tokens-jwt-are-dangerous-for-user-sessions/).

Now let's define a `CreateToken()` function to generate either **access** or **refresh** tokens.

**utils/token.go**

```go

```
