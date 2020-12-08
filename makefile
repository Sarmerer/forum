git-status:
	@status=$$(git status --porcelain); \
	if [ ! -z "$${status}" ]; \
	then \
		echo "There are uncommited changes, commit them before deploy."; \
		exit 1; \
	fi
push:
	git add .
	git commit -m "$m"
	git push -u origin master
go:
	cd api && bash -c  "go run main.go"
go-deploy: git-status
	heroku git:remote -a forum-api-sarmerer
	git subtree push --prefix api heroku master

vue:
	cd ./ui && npm i && npm run serve
vue-build:
	cd ./ui && npm i && npm run build
vue-deploy: git-status
	heroku git:remote -a forum-sarmerer
	git subtree push --prefix ui heroku master
dockerize:
	docker-compose build && docker-compose up -d client