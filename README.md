# auth-service-grpc-golang

## 1. API with Golang + MongoDB + Redis + Gin Gonic: Project Setup

### How to Setup Golang with MongoDB and Redis

#### Creating MongoDB and Redis Database with Docker-compose

The most straightforward way to get MongoDB and Redis database instances running on your machine is to use Docker 
and Docker-compose. I assume you already have Docker and Docker-compose installed on your computer.

To initialize a new Golang project, create a folder and open it with your preferred text editor. I recommend Visual Studio Code (VS Code) because it’s now the default text editor for new and professional developers.

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

The most essential part of our Golang application is to set up environment variables. This will allow us to store sensitive information and we can easily exclude them from being committed to a repository.

We are going to use the popular library `viper` package to load and validate the environment variables.

Before we can start installing third-party packages we need to initialize the Golang application.
Open the integrated terminal/command-line and run this command to initialize a new Golang project.

```shell
go mod init github.com/example/golang-test
```

Where `example` is your profile name on GitHub and `golang-test` is the name of your Golang project.
Once the Golang project has been initialized, run this command to install the [Golang viper package](https://github.com/spf13/viper).

```shell
go get github.com/spf13/viper
```

Next, create a `config` folder in the root directory and within this config folder create a `default.go` file.
We need to define a struct that will contain all the variables in the `app.env` file and also provide the appropriate value types.

**config/default.go**

```go
type Config struct {
	DBUri    string `mapstructure:"MONGODB_LOCAL_URI"`
	RedisUri string `mapstructure:"REDIS_URL"`
	Port     string `mapstructure:"PORT"`
}

```

The Golang viper package uses github.com/mitchellh/mapstructure under the hood for unmarshaling values which uses `mapstructure` tags by default.
You need to provide the correct variable name to the `mapstructure` tag.

Next, let’s define a function to load and read the content of the `app.env` file. The `LoadConfig` function will take the path to the `app.env` file as a parameter and return the `config` instance and a possible error if it exists.

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

