FROM golang:1.18.5

WORKDIR /app

COPY . .

EXPOSE 8080

RUN make build

RUN chmod +x ./server.out

CMD ./server.out
