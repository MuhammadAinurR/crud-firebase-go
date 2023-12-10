package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	firebase "firebase.google.com/auth"

	"google.golang.org/api/option"
)

var (
	firebaseApp *firebase.App
)

type Todo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	// Initialize Firebase Admin SDK
	opt := option.WithCredentialsFile("path/to/your/serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
	}

	firebaseApp = app

	// Initialize Firebase Realtime Database client
	client, err := app.DatabaseURL(context.Background())
	if err != nil {
		log.Fatalf("Error initializing database client: %v", err)
	}
	defer client.Close()

	router := mux.NewRouter()

	// Define API endpoints
	router.HandleFunc("/todos", createTodo).Methods("POST")
	router.HandleFunc("/todos", getTodos).Methods("GET")
	router.HandleFunc("/todos/{id}", getTodo).Methods("GET")
	router.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)

	// Get a reference to the Firebase Realtime Database
	ref := firebaseApp.DatabaseURL(context.Background())

	// Push the new todo to the database
	newTodoRef := ref.NewRef("/todos").Push(context.Background(), &todo)

	// Return the ID of the newly created todo
	todo.ID = newTodoRef.Key()
	json.NewEncoder(w).Encode(todo)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	// Get a reference to the Firebase Realtime Database
	ref := firebaseApp.DatabaseURL(context.Background())

	// Get all todos from the database
	todosRef := ref.NewRef("/todos")
	snapshot, err := todosRef.GetOrdered(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert snapshot to a slice of todos
	var todos []Todo
	err = snapshot.ForEach(func(snap db.DataSnapshot) error {
		var todo Todo
		if err := snap.Unmarshal(&todo); err != nil {
			return err
		}
		todo.ID = snap.Key()
		todos = append(todos, todo)
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the list of todos
	json.NewEncoder(w).Encode(todos)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	todoID := params["id"]

	// Get a reference to the Firebase Realtime Database
	ref := firebaseApp.DatabaseURL(context.Background())

	// Get the todo by ID from the database
	todosRef := ref.NewRef("/todos")
	snapshot, err := todosRef.Child(todoID).Get(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the todo exists
	if !snapshot.Exists() {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	// Convert snapshot to a todo
	var todo Todo
	if err := snapshot.Unmarshal(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	todo.ID = snapshot.Key()

	// Return the todo
	json.NewEncoder(w).Encode(todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	todoID := params["id"]

	var updatedTodo Todo
	_ = json.NewDecoder(r.Body).Decode(&updatedTodo)

	// Get a reference to the Firebase Realtime Database
	ref := firebaseApp.DatabaseURL(context.Background())

	// Update the todo by ID in the database
	todosRef := ref.NewRef("/todos")
	if err := todosRef.Child(todoID).Set(context.Background(), &updatedTodo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the updated todo
	updatedTodo.ID = todoID
	json.NewEncoder(w).Encode(updatedTodo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	todoID := params["id"]

	// Get a reference to the Firebase Realtime Database
	ref := firebaseApp.DatabaseURL(context.Background())

	// Delete the todo by ID from the database
	todosRef := ref.NewRef("/todos")
	if err := todosRef.Child(todoID).Delete(context.Background()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success message
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo deleted successfully"})
}
