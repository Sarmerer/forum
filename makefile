push:
	git add .
	git commit -m "$m"
	git push -u origin master

go:
	cd api && bash -c  "go run main.go"
go-deploy:
	heroku git:remote -a forum-api-sarmerer
	git add .
	git commit -m "go deploy"
	git subtree push --prefix api heroku master

vue:
	cd ./ui && npm i && npm run serve
vue-build:
	cd ./ui && npm i && npm run build
vue-deploy:
	heroku git:remote -a forum-sarmerer
	git add .
	git commit -m "vue deploy"
	git subtree push --prefix ui heroku master