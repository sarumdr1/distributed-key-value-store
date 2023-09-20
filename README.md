
### Distributed Key-Value Store Readme

This readme provides information on how to use a distributed key-value store service. The service allows you to perform basic operations such as putting, getting, and deleting key-value pairs, as well as replicating data.

### Getting Started
To use the key-value store service, follow the instructions below:

### Prerequisites
Before you begin, make sure you have the following installed on your system:

curl: Command-line tool for making HTTP requests. You can download it.

### Installation
1. Clone the repository or download the source code of the key-value store service to your local machine.

2, Navigate to the project directory.

3. Run the key-value store service:
go run main.go

This will start the service on http://localhost:8081.

### Usage
1. Put
To add a key-value pair to the store, use the following curl command:

curl -X POST -d '{"key": "name", "value": "Alice"}' http://localhost:8081/put

Replace "name" with your desired key.
Replace "Alice" with your desired value.

2. Get
To retrieve a value associated with a specific key, use the following curl command:

curl http://localhost:8081/get?key=name

Replace "name" with the key you want to retrieve.

3. Delete
To delete a key-value pair from the store, use the following curl command:

curl -X DELETE http://localhost:8081/delete?key=name

Replace "name" with the key you want to delete.

4. Replication
To replicate data or perform other operations related to replication, use the following curl command:

curl http://localhost:8081/replicate


### Testing

The unit test case is written on main_test.go.
To run the test case, use the following command:

go test -v
