package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

var channelAmqp *amqp.Channel

type Request struct {
	URL string `json:"url"`
}

func ParserHandler(c *gin.Context) {
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
		data, _ := json.Marshal(request)
	err := channelAmqp.Publish(
		"",
		"rss_urls",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(data),
		})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while publishing to RabbitMQ"})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"message": "success"})
}

func init() {
	amqpConnection, err := amqp.Dial("amqp://guest:password@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}

	channelAmqp, _ = amqpConnection.Channel()
}

func main() {
	router := gin.Default()
	router.POST("/parse", ParserHandler)
	router.Run(":5000")
}



/*

docker run -d --name rabbitmq -e RABBITMQ_DEFAUL_USER=user -e RABBITMQ_DEFAULT_PASS=password -p 8080:15672 -p 5672:5672 rabbitmq:3.8.21-management-alpine
*/