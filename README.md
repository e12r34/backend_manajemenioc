# Introduction
This program is your favorite golang fiber-backend that update ioc from many source so your ioc is up to date. This apps also can convert to many format (csv,json,thorlite, and many custom sensor) so that it compatible to many sensor.

# Technology
This code written in golang that use:
* [Go Fiber](https://gofiber.io/)

Your favorite lightweight web framework
* [swagger](https://github.com/swaggo/fiber-swagger)

Help you test your API and for nice documentation
* [mongodb](https://www.mongodb.com/)

Your favorite general purpose nosql database
* [godotenv](https://github.com/joho/godotenv)

So you can read config from .env
* [gocsv](https://github.com/gocarina/gocsv)

so you can handle csv easily


# Installation
1. Download and Install [MongoDB](https://www.mongodb.com/).
1. run command:

        go get -u github.com/gofiber/fiber/v2 go.mongodb.org/mongo-driver/mongo github.com/joho/godotenv github.com/go-playground/validator/v10 github.com/gocarina/gocsv

1. fill the `.env` and with your typical configuration
1. run your program with

        go run main.go
    
    or

        go build
    
    run executable `fiberioc.exe` or `fiberioc.sh` based on your OS