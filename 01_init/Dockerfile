FROM golang:1.23.0-alpine3.19 AS builder

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o ./main.exe ./cmd/app/main.go

FROM alpine:latest AS runner

COPY --from=builder /app/main.exe /main.exe

CMD [ "/main.exe" ]