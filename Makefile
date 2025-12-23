APP_NAME := go-connect-todo-server

TAG ?= dev
COUNTER_FILE := counter.txt
COUNTER := $(shell cat $(COUNTER_FILE))
# IMAGE_TAG := $(TAG)-$(COUNTER)
IMAGE_TAG := latest
IMAGE := $(APP_NAME):$(IMAGE_TAG)

KUBE_NAMESPACE ?= default
KUBE_CONTEXT ?= desktop

DOCKER := docker
KUBECTL := kubectl
BUF := buf

export IMAGE_TAG

.PHONY: init show increment gen

# Initialize counter
init:
	@echo "0" > $(COUNTER_FILE)
	@echo "Counter initialized to 0"

# Show current counter
show:
	@cat $(COUNTER_FILE)

# Increment counter
increment:
	@current=$$(cat $(COUNTER_FILE)); \
	next=$$((current + 1)); \
	echo $$next > $(COUNTER_FILE); \
	echo "Counter incremented: $$next"

dep:
	go mod tidy
	go mod vendor

fmt:
	gofumpt -l -w .

build: increment
	@echo "Building Docker image $(IMAGE)..."
# 	$(DOCKER) build -t go-connect-todo-server:$(IMAGE_TAG) .
	$(DOCKER) build -t go-connect-todo-server:latest .

# load image to k8s cluster registry
load-image: build
	kind load docker-image go-connect-todo-server:$(IMAGE_TAG) --name $(KUBE_CONTEXT)

# apply deployment and service to k8s cluster
deploy: load-image
	@echo "Deploying $(IMAGE) to Kubernetes..."
	IMAGE_TAG=$(IMAGE_TAG) envsubst < ./k8s/deployment.yaml.tpl | $(KUBECTL) apply -f - && $(KUBECTL) apply -f ./k8s/service.yaml

# port-forward service to local port
port-forward:
	$(KUBECTL) port-forward service/go-connect-todo 8080:80

# rollout deployment to k8s cluster
restart: load-image
	$(KUBECTL) rollout restart deployment/go-connect-todo

boot: build
# 	$(DOCKER) run --rm -p 8080:8080 --name go-connect-todo-test $(IMAGE)
	docker compose up -d go-app

# boot totally new infra from scratch
boot-new: load-image deploy port-forward

# boot and redeploy infra from scratch
boot-restart: load-image restart port-forward

gen:
	$(BUF) dep update
	$(BUF) lint
	$(BUF) generate