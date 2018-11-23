# product-service
Example REST service for managing products, written in Go using Chi.

## Configuration

### Environment Variables

##### PRODUCT_SERVICE_ENVIRONMENT

Controls log levels and configuration defaults. 

* DEV
* PRE_PROD
* PROD
 
##### PRODUCT_REPO_TYPE

* IN_MEMORY

##### PRODUCT_SERVICE_TIMEOUT

Incoming request timeout value in seconds

## Run

```go run main.go```
