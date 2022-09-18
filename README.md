# Welcome to BookStore (Hosted to GCP)
[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/golang.org/x/example)

Made with :heart: from Jeyakaran Karnan to EDIA ANZ.

# About

This repo contains the Go source code for a simple web app that can be deployed to Google Cloud Run. It is a demonstration of how to connect to a MySQL instance in Cloud SQL. The simple application uses a single table that stores all the details of Book and the demonstration was made to connect to that table in Cloud SQL and retrieve the data, Insert the data, and Update the data via API deployed in Google Cloud Run.

# Before you begin

## Setting up Cloud environment

1. Using the personal google account, initate the Google Cloud console, and the trial period of 3 months will be given and credits of $300 will be provided to play around.
2. Open the console, and create a project to host all the cloud tools.
3. Initate the Cloud SQL Instance. For this project, MySQL Database Engine has been used. Provide all the basic information, like instance name, password, zone, and Click on Create Instance. The instance will be up and running in few minutes. For extra reference, look at these [instructions](https://cloud.google.com/sql/docs/mysql/create-instance).
4. Create a Database called `books` in the created instance.
5. Create a new service under Cloud Run tool in Google Cloud Platform. Step by Step instruction given [here](https://cloud.google.com/run/docs/quickstarts/deploy-container?hl=en_US).
6. Create a new Github Repository and connect the sanme repository in the Cloud Run to enable the CI & CD using Cloudbuild.
7. ***Most important point to remember***: Appropriate IAM previleges should be provided to the service account used in Cloud Run. Make sure the account has the following previleges,

`cloud sql admin`

`cloud sql client`

## Setting up local environment

1. If you haven't already, set up a Go Development Environment by following the [Go setup guide](https://cloud.google.com/go/docs/setup) and
[create a project](https://cloud.google.com/resource-manager/docs/creating-managing-projects#creating_a_project).

2. Download and install the `cloud_sql_proxy` by
following the instructions
[here](https://cloud.google.com/sql/docs/mysql/sql-proxy#install). 

3. Download and install the ultimate testing framework of GoLang, `Ginkgo and Gomega` by following the instructions [here](https://onsi.github.io/ginkgo/)


## Running locally

1. To run this application locally, after writing all the handlers for the endpoints, just run `make run`. This will launch the application server, listening to 8080

2. To run the cloud_sql_proxy, run the below command,

```bash
./cloud_sql_proxy -instances=bookstore-362511:australia-southeast2:bookstore=tcp:5432  &
chmod +x /build/cloud_sql_proxy
```


## Why cloud_sql_proxy???

Cloud_sql_proxy is an utility for ensuring connections to your Google Cloud SQL. In this example, no TLS connections were made, so we are not using any certificates. But to ensure the connections is secured, it is always advisable to connect proxy via TLS connection. This is very much usedful to test your endpoint locally, Just sping up the proxy and the application, and test the application seamlessly before deploying to Google Cloud Platform.

# Deploying to Cloud 

## Containerisation
Docker has been used for image creation. Refer to `Dockerfile` in the repository to know all the steps.

## CI and CD 

1. Cloud Build is used in this project for CI and CD. When the Cloud Run service is created on Google Cloud Platform, the service is configured to trigger the build when the commit is pushed into develop branch. 
2. cloudbuild.yaml is used to create the docker image, push the image into Google container registry, and deploy into Cloud Run automatically.


# Important Information

## Database Schema
Sample below: `Primary key - id` (Auto incremented)
```json
{
        "id": "integer",
         "name": "string",
        "author": "string",
        "publishedDate": "string(YYYY-MM-DDD)",
        "price": "float64",
        "inStock": "boolean",
        "timeAdded": "time"
}
```

## End points served by service

1. Create a book: `POST /v1/books`

Sample Request Body: 
```json
{
        "name": "string",
        "author": "string",
        "publishedDate": "string(YYYY-MM-DDD)",
        "price": "float64",
        "inStock": "boolean",
        "timeAdded": "time"
}
```

Sample Response: 
```json
{
    "code": "8200",
    "description": "Success"
}
```


2. Get all books: `GET /v1/books/`

Sample Response: 
```json
[
    {
        "id": "integer",
         "name": "string",
        "author": "string",
        "publishedDate": "string(YYYY-MM-DDD)",
        "price": "float64",
        "inStock": "boolean",
        "timeAdded": "time"
    }
]
```
3. Get a specific book by ID: `GET /v1/books/{bookID}`

Sample Response Body:
```json
{
    "name": "string",
    "author": "string",
    "publishedDate": "string(YYYY-MM-DDD)",
    "price": "float64",
    "inStock": "boolean",
    "timeAdded": "time"
}
```
4. Update book: `PUT /v1/books/{bookID}`

Sample Request Body: 
```json
{
        "name": "string",
        "author": "string",
        "publishedDate": "string(YYYY-MM-DDD)",
        "price": "float64",
        "inStock": "boolean",
        "timeAdded": "time"
}
```

Sample Request Body: 
```json
{
    "code": "8200",
    "description": "Updated successfully"
}
```
5. Get Book by Wildcard: `GET /v1/books?name=*{some_character}`

Sample Request Body: 
```json
[
    {
        "id": "integer",
        "name": "string",
        "author": "string",
        "publishedDate": "Date",
        "price": "float64",
        "inStock": "Bool",
        "timeAdded": "Time"
    }
]
```
