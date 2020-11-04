FROM golang:latest
WORKDIR /app/api
COPY ./api .
RUN go mod download
RUN go build -o main .
EXPOSE 4433
CMD [ "./main" ]

FROM node:lts-alpine
RUN npm i -g http-server
WORKDIR /app/ui
COPY ./ui .
RUN npm i
RUN npm run build
EXPOSE 8080
CMD [ "http-server", "dist" ]

FROM scratch
WORKDIR /app 
COPY ./ui/dist . 
COPY ./api . 

EXPOSE 4000
CMD ["/app/server"]