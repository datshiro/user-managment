# user-managment

A simple API server handle user management including register and login written in Golang with Clean Architect applied


# Requirements

```
Go 1.21.4
Docker installed
```


# How to run 

To init Database and Redis, run:


```
make docker/up
```

To start API server, run: 


```
make run
```

> Note: API server can both load environments from .env file or flag argument


To stop server gracefully press `Ctr-C`

To shutdown all infrastructure, run:


```
make docker/down
```
