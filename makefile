run_app:
	bash -c "make -j run_server -j run_client"
docker: ## Build and run development docker container
	docker build -t sarmerer/forum .
	docker run -rm -p 4433:4433 sarmerer/forum
prod: ## Build and run production docker container
	docker build -t sarmerer/forum .
	docker run -d -p 4433:4433 sarmerer/forum
go:
	cd api && bash -c  "go run main.go"
vue:
	cd ./ui && npm i && npm run serve
push:
	git add .
	git commit -m "$m"
	git push -u origin master