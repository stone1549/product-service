# product-service
Example REST service for managing products, written in Go using Chi.

## Configuration

### Environment Variables

##### PRODUCT_SERVICE_ENVIRONMENT

Controls log levels and configuration defaults. 

* DEV
* PRE_PROD
* PROD
 
##### PRODUCT_SERVICE_REPO_TYPE

* IN_MEMORY
* POSTGRESQL
    * PRODUCT_SERVICE_PG_URL - Full connection string for PG

##### PRODUCT_SERVICE_TIMEOUT

Incoming request timeout value in seconds.

##### PRODUCT_SERVICE_PORT

Port to run service on.


## Run

```go run main.go```
