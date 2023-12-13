package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Item represents a generic item in our collection
type Item struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func initFirebase() *firestore.Client {
	opt := option.WithCredentialsFile("service.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error getting Firestore client: %v", err)
	}

	return client
}

func createItem(client *firestore.Client, item Item) (*firestore.DocumentRef, error) {
	docRef, _, err := client.Collection("items").Add(context.Background(), item)
	return docRef, err
}

func getItem(client *firestore.Client, id string) (*Item, error) {
	doc, err := client.Collection("items").Doc(id).Get(context.Background())
	if err != nil {
		return nil, err
	}
	var item Item
	doc.DataTo(&item)
	item.ID = doc.Ref.ID
	return &item, nil
}

func deleteItem(client *firestore.Client, id string) error {
	_, err := client.Collection("items").Doc(id).Delete(context.Background())
	return err
}

func updateItem(client *firestore.Client, id string, item Item) error {
	_, err := client.Collection("items").Doc(id).Set(context.Background(), item)
	return err
}

func getAllItems(client *firestore.Client) ([]Item, error) {
	var items []Item
	iter := client.Collection("items").Documents(context.Background())
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var item Item
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		items = append(items, item)
	}

	return items, nil
}

func handleCreateItem(client *firestore.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var item Item
		json.NewDecoder(r.Body).Decode(&item)
		docRef, err := createItem(client, item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		item.ID = docRef.ID
		json.NewEncoder(w).Encode(item)
	}
}

func handleGetItem(client *firestore.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		item, err := getItem(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(item)
	}
}

func handleUpdateItem(client *firestore.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var item Item
		json.NewDecoder(r.Body).Decode(&item)
		err := updateItem(client, id, item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		item.ID = id
		json.NewEncoder(w).Encode(item)
	}
}

func handleDeleteItem(client *firestore.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		err := deleteItem(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handleGetAllItems(client *firestore.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := getAllItems(client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(items)
	}
}

func main() {
	client := initFirebase()
	defer client.Close()

	r := mux.NewRouter()
	r.HandleFunc("/item", handleCreateItem(client)).Methods("POST")
	r.HandleFunc("/item/{id}", handleGetItem(client)).Methods("GET")
	r.HandleFunc("/item", handleGetAllItems(client)).Methods("GET")
	r.HandleFunc("/item/{id}", handleUpdateItem(client)).Methods("PUT")
	r.HandleFunc("/item/{id}", handleDeleteItem(client)).Methods("DELETE")

	http.Handle("/", r)
	log.Println("Server is starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
