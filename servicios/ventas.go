package main

import (
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
	Customer struct {
		Name     string `json:"name"`
		Lastname string `json:"lastname"`
		Email    string `json:"email"`
		Location struct {
			Address1   string `json:"address1"`
			Address2   string `json:"address2"`
			City       string `json:"city"`
			State      string `json:"state"`
			PostalCode string `json:"postalCode"`
			Country    string `json:"country"`
		} `json:"location"`
		Phone string `json:"phone"`
	}
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

	// Declara tres colas para los servicios de despacho, inventario y notificación
	despachoQueue, err := ch.QueueDeclare("despacho-cola", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	inventarioQueue, err := ch.QueueDeclare("inventario-cola", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	notificacionQueue, err := ch.QueueDeclare("notificacion-cola", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Define una estructura de orden de ejemplo
	orderData := `{
		"products": [
			{
				"title": "The Lord of the Rings",
				"quantity": 2,
				"price": 20
			}
			// Agrega el resto de los productos aquí...
		],
		"customer": {
			"name": "John",
			"lastname": "Doe",
			"email": "hleytondiaz@gmail.com",
			"location": {
				"address1": "123 Main Street",
				"address2": "Apt. 1",
				"city": "Anytown",
				"state": "CA",
				"postalCode": "91234",
				"country": "USA"
			},
			"phone": "555-555-5555"
		}
	}`

	// Publica la estructura de orden en las colas de los servicios
	err = ch.Publish("", despachoQueue.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(orderData),
	})
	if err != nil {
		log.Fatal(err)
	}

	err = ch.Publish("", inventarioQueue.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(orderData),
	})
	if err != nil {
		log.Fatal(err)
	}

	err = ch.Publish("", notificacionQueue.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(orderData),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Estructura de orden publicada en RabbitMQ para despacho, inventario y notificación")
}
