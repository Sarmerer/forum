FROM golang:latest
RUN mkdir /app
WORKDIR /app
COPY . /app
RUN go build -o main .
EXPOSE 4433
CMD ["./main"]