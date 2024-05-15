# CheckOut

Building and Deploying an application for orders with rate limiting and authentication

## Tasklist

- [x] Build a REST API based application in a language of your choice that has the following endpoints -

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

- [x] The application should store all the details in a Postgres or similar relational database.
- [x] The API endpoints should have global rate limiting implemented. For example, the service should only serve 5 requests across all replicas in a time bucket such as 1 minute. You can use a redis cache for this.
- [x] The API endpoints should have basic authentication implemented. Optionally, a token based authentication can be implemented.
- [x] Deploy the application on Kubernetes.
- [x] The database can be deployed on Kubernetes using an operator or a managed service.

## How to run this?

1. Install the following tools

- [Go](https://go.dev/doc/install)
- [Ko](https://ko.build/install/)
- [Kubectl](https://k8s-docs.netlify.app/en/docs/tasks/tools/install-kubectl/)
- [Curl](https://curl.se/download.html)
- [Kind](https://kind.sigs.k8s.io/) (or any kubernetes cluster)

2. Assuming using kind, create a cluster

```bash
kind create cluster --name checkout
```

3. Set environment variables

```bash
export KO_DOCKER_REPO='kind.local'
export KIND_CLUSTER_NAME='checkout'
```

4. Run the following command to deploy everything

```bash
kubectl kustomize config | ko apply -f -
```

5. Then run the following command and wait for everything to be ready

```bash
kubectl get all -n checkout
```

Once ready it will look similar to this:

```bash
NAME                            READY   STATUS    RESTARTS        AGE
pod/checkout-5b644b955b-49dwn   1/1     Running   4 (3m58s ago)   4m44s
pod/postgres-6f4fbcb784-fzh6t   1/1     Running   0               4m44s
pod/redis-6ffcf6865b-6qjlx      1/1     Running   0               4m44s

NAME                       TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
service/checkout           ClusterIP   10.96.89.103    <none>        8080/TCP         4m45s
service/postgres-service   NodePort    10.96.90.32     <none>        5432:31676/TCP   4m44s
service/redis-service      NodePort    10.96.210.246   <none>        6379:32533/TCP   4m44s

NAME                       READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/checkout   1/1     1            1           4m44s
deployment.apps/postgres   1/1     1            1           4m44s
deployment.apps/redis      1/1     1            1           4m44s

NAME                                  DESIRED   CURRENT   READY   AGE
replicaset.apps/checkout-5b644b955b   1         1         1       4m44s
replicaset.apps/postgres-6f4fbcb784   1         1         1       4m44s
replicaset.apps/redis-6ffcf6865b      1         1         1       4m44s
```

6. Port forward the checkout service

```bash
kubectl port-forward services/checkout 8080:8080 -n checkout
```

7. Now you can make API requests, to register:

```bash
curl -X POST "http://localhost:8080/register" \                                                   
-H "Accept: application/json" -H "Content-Type: application/json" \
-d '{ "username": "user", "password": "pass"}'

# Output: {"message":"user registered successfully"}
```

8. The project uses basic authantication, if your username is `user` and password is `pass` then run the following command to get token

```bash
echo "user:pass" | base64

# Output: dXNlcjpwYXNzCg==
```

9. You can then use this while authanticating:

```bash
curl -X POST "http://localhost:8080/api/orders" \
-H "Authorization: Basic dXNlcjpwYXNzCg==" \
-d '{
  "customer": "John Doe",
  "product_name": "Widget",
  "quantity": 2,
  "unit_price": 15.72,
  "order_date": "2023-06-09",
  "priority": "medium"
}'

# Output: {"order_id":1,"customer":"John Doe","product_name":"Widget","quantity":2,"unit_price":15.72,"order_date":"2023-06-09","priority":"medium"}
```

10. You can configure rate limiting (default 5/min) by changing `RATE` in `config/apiconfigmap.yaml`
