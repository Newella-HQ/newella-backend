# newella-backend
Newella microservices backends

Put .env file in the root directory of the project:
```dotenv
POSTGRES_HOST=postgres-db
POSTGRES_PORT=5432
POSTGRES_USERNAME=postgres
POSTGRES_PASSWORD=
POSTGRES_NAME=newella
POSTGRES_SSLMODE=disable

SERVER_HOST=localhost
AUTH_SERVER_PORT=8081
STATIC_SERVER_PORT=8082
USER_SERVICE_GRPC_PORT=50051

GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=

JWT_SIGNING_KEY=

LOG_LEVEL=debug
```