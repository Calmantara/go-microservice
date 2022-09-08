# Project Title

## Table of Contents

- [About](#about)

## About <a name = "about"></a>

this is project with data driven and window function query with rolling period

## Services
- go-gateway: service to route the request
- go-emitter: service to receive the payload to send to message broker. All messaage that wants to proceed asynchronously, should be sent through this service
- go-wallet: service to topup the balance for particular wallet, and get ballance. All wallet and balance related should be proceeded in this service

## On Development
- go-auth: to handle authentication and give access token to client

## Optional
- next development, current architecture should be change. we need to develop one more service called **go-forwarder**, this service is generic service to consume all payload from message broker and forward the payload to particular services. 
This service should has requirement:
    - minimal configured parameter:
        - throttling amount
        - topic consumed
    - throttling message, to control how many data to procced at once
    - easy to replicate without any duplication process
    - should insert/update the status of proceed payload for any host/service

## How to Run
this project is running on docker compose (can be v1 or v2)
1. to up the infrastructure services, run command below:
```
make infra
```
    1. postgres
    2. redis
    3. zookeeper (2)
    4. kafka broker
2. make sure kafka broker up and running
3. because kafka's topic is not automatically created, we should create it manually by run command below
```
make topic
```
4. to run the service, run command below
```
make up
```
5. go-gateway is running on port 8080 by default
## Architecture
![Alt text](./wallet.png?raw=true "wallet_architecture")
