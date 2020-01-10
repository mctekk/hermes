package main

import (
	"log"
	"net/http"
	"github.com/streadway/amqp"
	"github.com/joho/godotenv"
	"os"
	"encoding/json"
  )
  
type Message struct {
	UserId string
	MessageTypeId string
	AddedDate string 
	ObjectAction string
	ObjectType string
	Text string
}

func failOnError(err error, msg string) {
	if err != nil {
	  log.Fatalf("%s: %s", msg, err)
	}
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func main(){
	err := godotenv.Load()
	rabbitserver := os.Getenv("RABBITMQ_SERVER")
	conn, err := amqp.Dial(rabbitserver)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	queue := os.Args[1]
	q, err := ch.QueueDeclare(
	  queue, // name
	  false,   // durable
	  false,   // delete when unused
	  false,   // exclusive
	  false,   // no-wait
	  nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	  )
	  failOnError(err, "Failed to register a consumer")
	  
	  forever := make(chan bool)
	  
	  go func() {
		for d := range msgs {
		  log.Printf("Received a message: %s", d.Body)
		  processMessage(d.Body)
		}
	  }()
	  
	  log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	  <-forever
}

func processMessage(body []byte){
	bodyJson := string(body)
	log.Printf(bodyJson)
	var message Message
	json.Unmarshal([]byte(bodyJson), &message)
	// log.Printf("Species: %s, Description: %s", message.Title, message.Message)
}

