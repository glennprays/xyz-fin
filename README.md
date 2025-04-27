# XYZ Multifinance

## Get Started 
### Prerequisites
- Go 1.24 or later
- Docker 
- Make 
- Git 

### Running development mode 
Development mode is running PostgreSQL database and Swagger API documentation using Docker containers.
To run the development mode, use the following command:
```bash
make run-dev
```
To stop the development mode, use:
```bash
make stop-dev
```

### Database Migration
Ensure that you have installed [go-migrate](https://github.com/golang-migrate/migrate). Before migrating the database, create a database in your PostgreSQL.

To install go-migrate:
```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

To create migrations file:
```
migrate create -ext sql -dir migrations -seq <migrator_name>
```  
To run the database migrations:
- UP Migration
  ```
  migrate -database ${POSTGRESQL_URL} -path migrations up
  ```
- DOWN Migration
  ```
  migrate -database ${POSTGRESQL_URL} -path migrations -verbose down
  ```
> Note: in your local computer (without using docker) you need to add POSTGRESQL_URL as enviroment variable.
 ```
export POSTGRESQL_URL='postgres://postgres:example@localhost:5432/dev_xf?sslmode=disable'
 ```

### Swagger documentation
If you updated the swagger documentation, you need to refresh the Swagger UI. You can do this by running:
```bash
make swagger
```
To access the Swagger UI, open your browser and go to:
``` 
http://localhost:8080/
```

### Running code in development mode
To run the code in development mode, use the following command:
```bash
go run cmd/api/main.go
```

### Run the program in build container 
```bash
$ cd misc/production 
$ docker compose build 
$ docker compose up -d
```
> Note: make your you already run the database in dev mode

### Dummy Account 
- Account 1 (Budi)
Phone Number: `081234567890`
Password: `Password@123`
NIK: `1111`
- Account 2 (Annisa)
Phone Number: `087654321098`
Password: `Password@123`
NIK: `2222`

### Run Unit Test 
To run the unit tests, use the following command:
```bash
go test ./...
```
