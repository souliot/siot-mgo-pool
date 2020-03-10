package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/souliot/siot-mgo-pool/pool"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	pool.RegisterMgoPool("default", "mongodb://192.168.50.200:27017")
	client, _ := pool.GetMgoClient("default")
	defer pool.PutMgoClient("default", client)

	collection := client.Database("yapi").Collection("group")

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			fmt.Println(err)
		}
		// do something with result....
		fmt.Println(result)
	}

}
