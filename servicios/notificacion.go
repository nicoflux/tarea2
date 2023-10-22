package main

import (
	"encoding/json"
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
	q, err := ch.QueueDeclare("notificacion-cola", false, false, false, false, nil)
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
		var orderData map[string]interface{}
		err := json.Unmarshal(msg.Body, &orderData)
		if err != nil {
			log.Printf("Error al decodificar el mensaje: %s", err)
			continue
		}

		// Imprime el mensaje recibido en VM3
		fmt.Println("Mensaje recibido en VM3:")
		fmt.Printf("%v\n", orderData)
	}
	//Resto de la lógica que hay que implementar de aws
}
