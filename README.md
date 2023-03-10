# Parcel-Simulator

## HTTP Requests
Here are some demo request of the server:

### POST requests:

- Place order for customer with ID 11 (Authorized Customer)
```bash
curl localhost:8090/api/v1/parcel/11 -X POST -H "Content-Type: application/json" -d '{"receiver_name":"xyz", "receiver_number":"0131234131", "pickup_latitude":37.7749,"pickup_longitude":-122.4313,"dropoff_latitude":37.7886,"dropoff_longitude":-122.4314}'  
```
Sample response:
`{"order_id":1678341228942000002,"success":true}`

- Place order for customer with ID 16 (Unauthorized)
```bash
curl localhost:8090/api/v1/parcel/16 -X POST -H "Content-Type: application/json" -d '{"receiver_name":"xyz", "receiver_number":"0131234131", "pickup_latitude":37.7749,"pickup_longitude":-122.4313,"dropoff_latitude":37.7886,"dropoff_longitude":-122.4314}' 
```
Sample response:
`No authorized customer found with id: 16`

### Get requests:

- Ping the homepage
```bash
curl localhost:8090
```
Sample response:
`Welcome to the Homepage of parcel simulator!`

- Get the current status of the order with ID 1678341228942000002 (Valid Order Id)
```bash
curl localhost:8090/api/v1/order/status/1678341228942000002
```
Sample response:
`{"order_id":1678345319223000002,"order_status":"Accepted"}`



## System Design

![System Design](./media/system-design.png)

### High level system design flow

1. The API will receive a new order request and response with a unique order ID
2. The API will send events to the order processing task queue and Order status tracking queue
3. The Order processing service will process the order requests asynchronously
4. The Order processing service will find a nearest rider for the request based on availability
5. If there is no available rider, then send another requeue event to the order processing task
6. Rider service will send a request to the Location service to track available Riders location
7. The Location service will update the current location to the cache
8. Order Processing service will update the oder status in the cache
9. Order Status Queue will get updated status from the cache periodically
10. Order Status Queue will send event to the notification service
11. The Notification service will send notification to the users

## DB Schema

![DB Schema](./media/parcel.png)