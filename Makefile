go:	## Run Go API
	bash -c "go run main.go"
tests: ## Run Go tests
	bash -c "go test ./api/auth"
docker: ## Build and run development docker container
	docker build -t sarmerer/forum .
	docker run -rm -p 4433:4433 sarmerer/forum
prod: ## Build and run production docker container
	docker build -t sarmerer/forum .
	docker run -d -p 4433:4433 sarmerer/forum