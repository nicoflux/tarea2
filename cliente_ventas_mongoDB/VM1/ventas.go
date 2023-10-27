package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func insertData(order Order) string {

	// Connect to remote MongoDB server
	clientOptions := options.Client().ApplyURI("mongodb://10.10.11.221:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	// Insert JSON data into collection
	collection := client.Database("tarea2").Collection("orders")
	insertResult, err := collection.InsertOne(context.TODO(), order)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted document with ID:", insertResult.InsertedID)

}

func main() {
	file, _ := os.Open("products.json")
	defer file.Close()
	var order Order
	json.NewDecoder(file).Decode(&order)
	insertData(order)
}
