# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server \
	&& CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/importer ./cmd/importer

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /out/server /usr/local/bin/server
COPY --from=build /out/importer /usr/local/bin/importer
EXPOSE 8081
CMD ["/usr/local/bin/server"]
