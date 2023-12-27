FROM golang:1.20-alpine as build

WORKDIR /app
COPY . .

RUN go mod tidy

FROM golang:1.20-alpine as run

WORKDIR /app

COPY --from=build /app/github-graph-drawer ./run

EXPOSE 8080

CMD ["./run"]
