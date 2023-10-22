package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// Estructura de una orden de compra
type Order struct {
	Products []struct {
		Title    string  `json:"title"`
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	} `json:"products"`
}

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
	q, err := ch.QueueDeclare("inventario-cola", false, false, false, false, nil)
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
		var order Order
		err := json.Unmarshal(msg.Body, &order)
		if err != nil {
			log.Printf("Error al decodificar el mensaje: %s", err)
			continue
		}

		// Simplemente muestra el mensaje
		fmt.Printf("Mensaje recibido en inventario.go: %+v\n", order)
	}
}
