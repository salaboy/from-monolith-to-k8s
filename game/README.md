# Knative Functions Game 

This tutorial shows you how to deploy an instance of the Knative Functions quiz game. 
This game is built using different functions to evaluate and score the answers provided by the players. 
There is also a function to obtain a leaderboard based on the data recorded from each game session.

## Architecture

This application is based on functions, but it relies on a `frontend` component that serves as an API Gateway to connect the User Interface to the functions, so functions are not exposed outside the cluster. 

![game-architecture.png](game-architecture.png)
  
Each Game level is a different function that evaluates the submitted answers and score the results. 
All the state is kept into a Redis database and each function access Redis by creating a connection to it.
  
  

## Prerequisites

- Kubernetes Cluster
- Install Knative Serving
- (Optional for the second part) Install Knative Eventing
- Install Knative `func` CLI

## Installing the game functions and components
  
Once you have Knative Serving installed you can installing the FrontEnd Service
  
```
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: game-frontend
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/minScale: "1"
    spec:
      containers:
        - image: salaboy/fmtok8s-game-frontend:0.1.0
          imagePullPolicy: Always
          livenessProbe:
            httpGet:
              path: /actuator/health
          ports:
            - containerPort: 8080
              name: http1
              protocol: TCP
          readinessProbe:
            httpGet:
              path: /actuator/health
```


