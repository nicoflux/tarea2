package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	Title       string  `bson:"title"`
	Author      string  `bson:"author"`
	Genre       string  `bson:"genre"`
	Pages       int     `bson:"pages"`
	Publication string  `bson:"publication"`
	Quantity    int     `bson:"quantity"`
	Price       float64 `bson:"price"`
}

func main() {
	// Configura la conexión a la base de datos MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://admin:admin@10.10.11.221:27017/tarea2")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Obtiene una referencia a la base de datos "tarea2"
	db := client.Database("tarea2")

	// Obtiene una referencia a la colección "products"
	collection := db.Collection("products")

	// Crea un documento de producto
	product := Product{
		Title:       "The Lord of the Rings",
		Author:      "J.R.R. Tolkien",
		Genre:       "Fantasy",
		Pages:       1224,
		Publication: "1954",
		Quantity:    98,
		Price:       20.00,
	}

	// Inserta el documento en la colección
	_, err = collection.InsertOne(context.Background(), product)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Documento insertado en la colección 'products'")
}
