# Versión de Go
ARG GO_VERSION=1.27

# BUILD
FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /src

# Descarga dependencias
COPY go.mod go.sum ./
RUN go mod download

# copair codigo
COPY . .

#compilar
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# RUNTIME
FROM alpine:3

## creando usuario sin privilegios
RUN addgroup -S app && adduser -S -G app app
## copia del ejecutable
COPY --from=build /out/server /usr/local/bin/server

USER app
ENV PORT=3000
EXPOSE 3000
ENTRYPOINT ["server"]