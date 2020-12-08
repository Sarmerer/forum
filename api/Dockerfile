FROM golang:latest
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o forum-api .
EXPOSE 4433
CMD [ "./forum-api" ]