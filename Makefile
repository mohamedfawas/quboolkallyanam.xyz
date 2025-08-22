# List of services (matches folder names in services/)
SERVICES := auth user chat payment notification gateway

# Output dir for built binaries (relative to repo root)
ARTIFACTS_DIR := deploy/docker/artifacts

# Common Go build settings
GO_BUILD_FLAGS := -ldflags="-s -w"
GO_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Docker compose invocation used by targets below
DOCKER_COMPOSE := docker compose --env-file .env -f deploy/compose/docker-compose.yml

.PHONY: all build clean $(SERVICES) docker-build docker-build-% rebuild

# Default
all: build

# Build all services' artifacts
build: $(ARTIFACTS_DIR) $(SERVICES)

# ensure artifacts dir exists
$(ARTIFACTS_DIR):
	@echo ">>> Creating artifacts dir $(ARTIFACTS_DIR)"
	mkdir -p $(ARTIFACTS_DIR)

# Generic rule: build each service into artifacts dir
$(SERVICES): $(ARTIFACTS_DIR)
	@echo ">>> Building $@ service"
	$(GO_ENV) go build -C services/$@ $(GO_BUILD_FLAGS) -o ../../$(ARTIFACTS_DIR)/$@-service ./cmd

# Clean build artifacts
clean:
	@echo ">>> Cleaning build artifacts"
	rm -rf $(ARTIFACTS_DIR)/*

# Build images for all services via docker compose (ensures artifacts exist)
docker-build: build
	@echo ">>> Running docker compose build for all services"
	$(DOCKER_COMPOSE) build

# Build image for single service (build artifact first then compose build that service)
# usage: make docker-build-chat
docker-build-%: %
	@echo ">>> Building artifact for $* and building docker image for $*"
	@$(DOCKER_COMPOSE) build $*

# Full rebuild: clean, build artifacts, build images
rebuild: clean docker-build
	@echo ">>> Rebuild complete"

# ----------------------------
# Artifact Registry settings
# ----------------------------
REGION     ?= asia-south1
PROJECT_ID ?= qubool-kallyanam
REPO       ?= qubool-kallyanam-docker
TAG        ?= latest

# Push all services to Artifact Registry
push: docker-build
	@for svc in $(SERVICES); do \
	  echo ">>> Pushing $$svc image"; \
	  docker tag qubool-kallyanam-$$svc:latest \
	    $(REGION)-docker.pkg.dev/$(PROJECT_ID)/$(REPO)/$$svc:$(TAG); \
	  docker push $(REGION)-docker.pkg.dev/$(PROJECT_ID)/$(REPO)/$$svc:$(TAG); \
	done

# Push only one service (usage: make push-auth)
push-%: docker-build-%
	@echo ">>> Pushing $* service image"
	docker tag qubool-kallyanam-$*:latest \
	  $(REGION)-docker.pkg.dev/$(PROJECT_ID)/$(REPO)/$*:$(TAG)
	docker push $(REGION)-docker.pkg.dev/$(PROJECT_ID)/$(REPO)/$*:$(TAG)
