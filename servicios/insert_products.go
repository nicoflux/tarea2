package main

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
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
	// Establece la sesión de MongoDB
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Obtiene una referencia a la base de datos "tarea2"
	db := session.DB("tarea2")

	// Obtiene una referencia a la colección "products"
	collection := db.C("products")

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
	err = collection.Insert(product)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Documento insertado en la colección 'products'")
}
