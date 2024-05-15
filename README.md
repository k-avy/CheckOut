# CheckOut

Simple REST API application for managing orders.

## Problem Statement

Task: Building and Deploying an application for orders with rate limiting and authentication

1. Build a REST API based application in a language of your choice that has the following endpoints -

- GET /orders: Retrieve all orders from the database.
- POST /orders: Create a new order by accepting request data and saving it to the database.
- GET /orders/{orderId}: Retrieve a specific order by its ID.
- PUT /orders/{orderId}: Update an existing order with new data.
- DELETE /orders/{orderId}: Delete an order from the database.
Example of the order -

```json
{
  "order_id": 12345,
  "customer": "John Doe",
  "product_name": "Widget",
  "quantity": 2,
  "unit_price": 15.72,
  "order_date": "2023-06-09",
  "priority": "medium"
}
```

2. The application should store all the details in a Postgres or similar relational database.
3. The API endpoints should have global rate limiting implemented. For example, the service should only serve 5 requests across all replicas in a time bucket such as 1 minute. You can use a redis cache for this.
4. The API endpoints should have basic authentication implemented. Optionally, a token based authentication can be implemented.
5. Deploy the application on Kubernetes.
6. The database can be deployed on Kubernetes using an operator or a managed service.
