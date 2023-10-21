FROM golang:1.21
WORKDIR /brand
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o out .

FROM alpine:3.18
WORKDIR /brand
COPY --from=0 /brand/out /bin/brand

ENV VERSION="v0.1.0"
ARG TOKEN_DISCORD_APPLICATION
ENV TOKEN_DISCORD_APPLICATION=${TOKEN_DISCORD_APPLICATION}
ENV ENDPOINT_COSMOS="https://neko03cosmos.documents.azure.com:443/"

EXPOSE 80
CMD ["brand"]
