package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

type Product struct {
	Title    string `json:"title"`
	Quantity int    `json:"quantity"`
}

func main() {
	// Establece una conexión a MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://admin:admin@10.10.11.221:27017/tarea2"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("tarea2")
	collection := db.Collection("products")

	// Conectarse a RabbitMQ
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

	q, err := ch.QueueDeclare("inventario-cola", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Escuchar mensajes de la cola
	for msg := range msgs {
		var order Order
		err := json.Unmarshal(msg.Body, &order)
		if err != nil {
			log.Println("Error al decodificar el mensaje:", err)
			continue
		}

		// Restar la cantidad de productos en la base de datos
		for _, product := range order.Products {
			filter := bson.M{"title": product.Title}
			update := bson.M{"$inc": bson.M{"quantity": -product.Quantity}}

			_, err := collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				log.Println("Error al actualizar la base de datos:", err)
			}
		}

		fmt.Printf("Se actualizó el stock: %s\n", order.OrderID)
	}
}
