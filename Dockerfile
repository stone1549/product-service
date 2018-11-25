FROM golang:1.8

ENV PRODUCT_SERVICE_ENVIRONMENT=DEV
ENV	PRODUCT_SERVICE_REPO_TYPE=IN_MEMORY
ENV PRODUCT_SERVICE_TIMEOUT=60
ENV PRODUCT_SERVICE_PORT=8080
ENV PRODUCT_SERVICE_PG_URL=NA
ENV PRODUCT_SERVICE_INIT_DATASET=SMALL

WORKDIR /go/src/github.com/stone1549/product-service/
COPY . .

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["poduct-service"]
