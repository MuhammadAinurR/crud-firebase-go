# Firestroe CRUD Rest-Api with Go and Gorilla Mux
This project demonstrates a simple CRUD (Create, Read, Update, Delete) Rest-Api using Google Firestore as the backend database and Go (Golang) for the server implementation. The Api Allows you to manage items with basic information such as ID, Name, and Description.
<br>
<br>

## Prerequisites
Before running the application, make sure you have the following prerequisites:
- Google Firestore Project: Create Google Firebase Project and enable the Firestore API.
- Go: Install the Go programming language on your machine. you can download it from the <a href="https://go.dev/dl/">Official Go website</a>
<br>
<br>

### SETUP
### Clone the repository

    git clone https://github.com/MuhammadAinurR/crud-firebase-go
   
### Navigate to the project directory

    cd crud-firebase-go
    
### Install the dependencies

    go get -u ./...

<br>
<br>

## Configuration
Ensure that the 'service.json' file containing your Google Firebase service account key is in the project directory.

<br>
<br>

## Usage
### Run the application
    go run main.go
###Access the API using tools like `curl`, Postman, or any other API testing tools

<br>
<br>

## API Endpoints
### Create Item
Endpoint: `POST /item`
Request Body:

    {
      "name": "Sample Item",
      "description": "This is a sample item."
    }
![img_post](https://github.com/MuhammadAinurR/crud-firebase-go/blob/main/img/Screenshot%202023-12-13%20at%2009.40.23.png?raw=true)
![img_getAllItems_afterPost](https://github.com/MuhammadAinurR/crud-firebase-go/blob/main/img/Screenshot%202023-12-13%20at%2009.40.41.png?raw=true)

    
### Get All Item
Endpoint: `GET /item`
Response:

    {
      "id": "itemID1",
      "name": "Sample Item 1",
      "description": "This is the first sample item."
    },
    {
      "id": "itemID2",
      "name": "Sample Item 2",
      "description": "This is the second sample item."
    }
![img_getAllItems](https://github.com/MuhammadAinurR/crud-firebase-go/blob/main/img/Screenshot%202023-12-13%20at%2009.37.37.png?raw=true)


### Get Specific Item
Endpoint: `GET /item/{id}`
Response:

    {
      "id": "itemID",
      "name": "Sample Item",
      "description": "This is a sample item."
    }
![img_getSpecificItem](https://github.com/MuhammadAinurR/crud-firebase-go/blob/main/img/Screenshot%202023-12-13%20at%2009.37.51.png?raw=true)
### Update Item
Endpoint: `PUT /item/{id}`
Request Body:

    {
      "name": "Updated Item",
      "description": "This item has been updated."
    }
![img_update](https://github.com/MuhammadAinurR/crud-firebase-go/blob/main/img/Screenshot%202023-12-13%20at%2009.38.12.png?raw=true)
![img_getAllItem_update](https://github.com/MuhammadAinurR/crud-firebase-go/blob/main/img/Screenshot%202023-12-13%20at%2009.38.26.png?raw=true)
### Delete Item
Endpoint: `DELETE /item/{id}`
![img_delete](https://github.com/MuhammadAinurR/crud-firebase-go/blob/main/img/Screenshot%202023-12-13%20at%2009.38.44.png?raw=true)
![img_getAllItems_delete](https://github.com/MuhammadAinurR/crud-firebase-go/blob/main/img/Screenshot%202023-12-13%20at%2009.38.55.png?raw=true)
<br>
<br>

## Closing the Application
After you finish using the application, stop the server by pressing `ctrl+c` in the terminal
