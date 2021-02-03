git-status:
	@status=$$(git status --porcelain); \
	if [ ! -z "$${status}" ]; \
	then \
		echo "There are uncommited changes, commit them before deploy."; \
		exit 1; \
	fi

go:
	cd api && bash -c  "go run main.go"
goexec:
	cd api && ./main
go-build:
	cd api && bash -c  "go build -o main"
go-build-dbauth:
	cd api && bash -c "go build -o main --tags \"sqlite_userauth\""

vue:
	cd ./ui && npm i && npm run serve
vue-build:
	cd ./ui && npm i && npm run build


dockerize:
	docker-compose build && docker-compose up -d client