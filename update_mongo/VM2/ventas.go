package main

import (
	"context"
	"encoding/json"
	"fmt"
	"grpc-golang/pb"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

var mongo_Client *mongo.Client

type server struct {
	pb.OrderServiceServer
}

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

func (s *server) Order(ctx context.Context, req *pb.OrderServiceRequest) (*pb.OrderServiceReply, error) {
	receivedJSON, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	log.Printf("Recibido : %s", string(receivedJSON))
	var order Order
	err = json.Unmarshal(receivedJSON, &order)
	if err != nil {
		log.Fatal(err)
	}

	orderId := insertData(order)
	updateData(order)
	return &pb.OrderServiceReply{
		OrderResponse: fmt.Sprintf("You order id is : %s", orderId),
	}, nil
}

func connectToMongoDB() (*mongo.Client, error) {
	//URI := os.Getenv("CONNECTION_STRING")
	URI := "mongodb://admin:admin@127.0.0.1:27017/tarea2"
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(URI).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client, nil
}

func closeMongoDBConnection(client *mongo.Client) {
	if err := client.Disconnect(context.Background()); err != nil {
		fmt.Println("Error al desconectar de MongoDB:", err)
	}
}

func insertData(order Order) string {

	order.ID = primitive.NewObjectID()
	order.OrderID = order.ID.Hex()

	collection := mongo_Client.Database("tarea2").Collection("orders")

	resp, err := collection.InsertOne(context.Background(), order)
	if err != nil {
		fmt.Println("Error al insertar datos en MongoDB:", err)
		return ""
	}

	fmt.Println("Documento insertado con éxito, ID:", resp.InsertedID)
	myObjectId := resp.InsertedID.(primitive.ObjectID)
	return myObjectId.Hex()
}

func updateData(order Order) {
	collection := mongo_Client.Database("tarea2").Collection("orders")
	filter := bson.M{"_id": order.OrderID}
	updateDoc := bson.M{
		"$set": bson.M{
			"deliveries.0.trackingNumber":             "1234567890",
			"deliveries.0.shippingMethod":             "USPS",
			"deliveries.0.shippingAddress.name":       order.Customer.Name,
			"deliveries.0.shippingAddress.lastname":   order.Customer.Lastname,
			"deliveries.0.shippingAddress.address1":   order.Customer.Location.Address1,
			"deliveries.0.shippingAddress.address2":   order.Customer.Location.Address2,
			"deliveries.0.shippingAddress.city":       order.Customer.Location.City,
			"deliveries.0.shippingAddress.state":      order.Customer.Location.State,
			"deliveries.0.shippingAddress.postalCode": order.Customer.Location.PostalCode,
			"deliveries.0.shippingAddress.country":    order.Customer.Location.Country,
			"deliveries.0.shippingAddress.phone":      order.Customer.Phone,
		},
	}
	// Execute update one on orders collection
	result, err := collection.UpdateOne(context.TODO(), filter, updateDoc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v order\n", result.ModifiedCount)
}

func main() {
	var err error
	mongo_Client, err = connectToMongoDB()
	if err != nil {
		fmt.Println("Error al conectar a MongoDB:", err)
	}
	defer closeMongoDBConnection(mongo_Client)
	listener, err := net.Listen("tcp", ":8080")
	fmt.Println("Server is running on port 8080")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
