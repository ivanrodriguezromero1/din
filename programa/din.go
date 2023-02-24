package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Host string
	Port int
}

const connStr = "mongodb://{{.Host}}:{{.Port}}"

func conectar() *mongo.Client {
	config := Config{
		Host: "localhost",
		Port: 27017,
	}
	var connBuffer bytes.Buffer
	tmpl, err1 := template.New("conn").Parse(connStr)
	if err1 != nil {
		log.Fatal(err1)
	}
	if err2 := tmpl.Execute(&connBuffer, config); err2 != nil {
		log.Fatal(err2)
	}
	connectionString := connBuffer.String()
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client
	// defer client.Disconnect(ctx)
	// collection := client.Database("din1").Collection("collectionDin1")
}
func operation(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	i := 1
	parametros := []string{}
	valores := []string{}
	for scanner.Scan() {
		if i%2 != 0 {
			parametros = append(parametros, scanner.Text())
		} else {
			valores = append(valores, scanner.Text())
		}
		i++
	}
	client := conectar()
	defer client.Disconnect(context.Background())
	var docs []bson.M
	for i := 0; i < len(parametros); i++ {
		doc := bson.M{
			"parametro": parametros[i],
			"valor":     valores[i],
			// "fechaDeCreacion": time.Now(),
		}
		docs = append(docs, doc)
	}
	// document := []interface{}{
	// 	bson.M{valores[0]: docs},
	// }
	// limpiar(client)
	mostrar(client)
	// insertar(client, document)
}
func limpiar(client *mongo.Client) {
	err := client.Database("din1").Drop(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Base de datos limpiada exitosamente")
}
func mostrar(client *mongo.Client) {
	collection := client.Database("din1").Collection("collectionDin1")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}
func insertar(client *mongo.Client, docs []interface{}) {
	collection := client.Database("din1").Collection("collectionDin1")
	_, err := collection.InsertMany(context.Background(), docs)
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	path := "archivos"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		operation(path + "/" + file.Name())
	}
}
