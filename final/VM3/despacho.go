package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Order struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	OrderID  string             `json:"orderID"`
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
	Deliveries []struct {
		ShippingAddress struct {
			Name       string `json:"name"`
			Lastname   string `json:"lastname"`
			Address1   string `json:"address1"`
			Address2   string `json:"address2"`
			City       string `json:"city"`
			State      string `json:"state"`
			PostalCode string `json:"postalCode"`
			Country    string `json:"country"`
			Phone      string `json:"phone"`
		} `json:"shippingAddress"`
		ShippingMethod string `json:"shippingMethod"`
		TrackingNumber string `json:"trackingNumber"`
	} `json:"deliveries"`
}
type Delivery struct {
	ShippingAddress struct {
		Name       string `json:"name"`
		Lastname   string `json:"lastname"`
		Address1   string `json:"address1"`
		Address2   string `json:"address2"`
		City       string `json:"city"`
		State      string `json:"state"`
		PostalCode string `json:"postalCode"`
		Country    string `json:"country"`
		Phone      string `json:"phone"`
	} `json:"shippingAddress"`
	ShippingMethod string `json:"shippingMethod"`
	TrackingNumber string `json:"trackingNumber"`
}

func updateData(order Order) {
	URI := "mongodb://admin:admin@10.10.11.221:27017/tarea2"
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(URI).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connected to MongoDB!")
	}

	collection := client.Database("tarea2").Collection("orders")
	filter := bson.M{"orderid": order.OrderID}
	var delivery Delivery
	delivery.ShippingAddress.Name = order.Customer.Name
	delivery.ShippingAddress.Lastname = order.Customer.Lastname
	delivery.ShippingAddress.Address1 = order.Customer.Location.Address1
	delivery.ShippingAddress.Address2 = order.Customer.Location.Address2
	delivery.ShippingAddress.City = order.Customer.Location.City
	delivery.ShippingAddress.State = order.Customer.Location.State
	delivery.ShippingAddress.PostalCode = order.Customer.Location.PostalCode
	delivery.ShippingAddress.Country = order.Customer.Location.Country
	delivery.ShippingAddress.Phone = order.Customer.Phone
	delivery.ShippingMethod = "USPS"
	delivery.TrackingNumber = "1234567890"
	updateDoc := bson.M{
		"$set": bson.M{
			"deliveries.0": delivery,
		},
	}
	fmt.Println(delivery)
	// Execute update one on orders collection
	result, err := collection.UpdateOne(context.TODO(), filter, updateDoc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v order\n", result.ModifiedCount)
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
		//fmt.Printf("Mensaje recibido en VM3: %s\n", msg.Body)
		var order Order
		err = json.Unmarshal(msg.Body, &order)
		if err != nil {
			log.Printf("Error al decodificar el mensaje: %s", err)
			continue
		}
		updateData(order)
	}
}
