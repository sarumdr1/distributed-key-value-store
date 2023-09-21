
### Distributed Key-Value Store Readme

This readme provides information on how to use a distributed key-value store service. The service allows you to perform basic operations such as putting, getting, and deleting key-value pairs, as well as replicating data.

### Getting Started
To use the key-value store service, follow the instructions below:

### Prerequisites
Before you begin, make sure you have the following installed on your system:
go

Check go version 
go version

### Installation
1. Clone the repository or download the source code of the key-value store service to your local machine.

2, Navigate to the project directory.

3. Run the key-value store service:
go run main.go

This will start the service and ask for the options
Select an option:
1. Put
2. Get
3. Delete
4. Replicate
5. Quit
   
### Usage
1. Put
To add a key-value pair to the store, choose option 1
Enter key: name
Enter value:Alice

Replace "name" with your desired key.
Replace "Alice" with your desired value.

2. Get
To retrieve a value associated with a specific key, choose option 2

Enter key: name
Replace "name" with the key you want to retrieve.

3. Delete
To delete a key-value pair from the store,choose option 3

Enter key: name
Replace "name" with the key you want to delete.

4. Replication
To replicate data or perform other operations related to replication, choose option 4


### Testing

The unit test case is written on main_test.go.
To run the test case, use the following command:

go test -v
