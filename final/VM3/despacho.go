package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	// Configura la conexión a RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	// Declara una cola
	q, err := ch.QueueDeclare("despacho-cola", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Consume mensajes de la cola
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Procesa los mensajes
	for msg := range msgs {
		fmt.Printf("Mensaje recibido en VM3: %s\n", msg.Body)
	}
}
