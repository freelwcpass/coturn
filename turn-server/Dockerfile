FROM golang:1.21-alpine AS builder
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/turn-server ./main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /out/turn-server /turn-server
COPY .env /.env
EXPOSE 3478/udp
EXPOSE 49160-49200/udp
ENTRYPOINT ["/turn-server"]