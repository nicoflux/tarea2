package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/streadway/amqp"
)

type Order struct {
	OrderID  string `json:"orderID"`
	GroupID  string `json:"groupID"`
	Products []struct {
		Title       string  `json:"title"`
		Author      string  `json:"author"`
		Genre       string  `json:"genre"`
		Pages       int     `json:"pages"`
		Publication string  `json:"publication"`
		Quantity    int     `json:"quantity"`
		Price       float64 `json:"price"`
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
	} `json:"customer"`
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
		var order Order
		err := json.Unmarshal(msg.Body, &order)
		if err != nil {
			log.Printf("Error al decodificar el mensaje: %s", err)
			continue
		}
		order.GroupID = "G4b!7S9k#3"
		// Imprime el mensaje recibido en VM3
		fmt.Println("Orden n° ", order.OrderID)

		// Envía el mensaje a la API Gateway
		client := &http.Client{}
		requestByte, _ := json.Marshal(order)
		data := bytes.NewReader(requestByte)
		req, err := http.NewRequest("POST", "https://sjwc0tz9e4.execute-api.us-east-2.amazonaws.com/Prod", data)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Respuesta de la API Gateway:", resp.StatusCode)
	}
}
