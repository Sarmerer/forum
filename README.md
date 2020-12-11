[Live demo](https://forum-sarmerer.herokuapp.com)
---
- [Live demo](#live-demo)
- [Getting started](#getting-started)
    - [Requirements](#requirements)
- [License](#license)

## Getting started
#### Requirements
* [Go](https://golang.org) >= v1.12
* [Node JS](https://nodejs.org/)
* [Vue CLI](https://cli.vuejs.org/guide/installation.html)

Once you have fulfilled all the requirements:

```shell
make go
# in separate terminal
make vue
```
or without make:
```shell
cd api && bash -c  "go run main.go"
# in separate terminal
cd ui && npm i && npm run serve
```
Then wait for vue to finish the build, and open the link vue cli has provided

License
---
MIT - [Sarmerer](https://github.com/sarmerer), [sarmai](https://github.com/sarmai)
