# gohotel

## How to run
1. Download and install `Go` and `Docker Compose`
2. Clone project via `git clone https://github.com/DirectDuck/gohotel.git`
3. Run `go mod download` to install dependencies
4. Create `.env` file (by example in `.env.example`)
5. Run `docker-compose up -d` to start database
6. Run `make run`

## Modules
- **types**
    - Describes types (models) of the project
- **db**
    - Handles database connection
    - Implemens basic database operations
- **controllers**
    - Ties up types and database
    - Implements CRUD and all other business logic
- **api**
    - Handles HTTP requests to server
    - Serializes data from request to defined types
- **services**
    - Stores different microservices
    - **roomprices**
        - Calculates prices for rooms
        - Acts as gRPC server

## Side note
This project was created for learning purposes only.
Thus, some features like tests, documentation, logging were only partially implemented, or completly excluded.
