run_app:
	bash -c "make -j run_server -j run_client"
docker: ## Build and run development docker container
	docker build -t sarmerer/forum .
	docker run -rm -p 4433:4433 sarmerer/forum
prod: ## Build and run production docker container
	docker build -t sarmerer/forum .
	docker run -d -p 4433:4433 sarmerer/forum
run_server:
	bash -c "go run main.go"
run_client:
	cd ./ui/app && ng serve --open