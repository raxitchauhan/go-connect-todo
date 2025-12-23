apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-connect-todo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: go-connect-todo
  template:
    metadata:
      labels:
        app: go-connect-todo
    spec:
      containers:
        - name: go-api
          image: go-connect-todo-server:${IMAGE_TAG}
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 8080
          env:
            - name: PORT
              value: "8080"
