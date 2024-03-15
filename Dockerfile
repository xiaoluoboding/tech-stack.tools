# Build frontend dist.
FROM node:18.19.1-alpine3.18 AS frontend

WORKDIR /frontend-build

COPY ./frontend/ .

RUN yarn
RUN yarn generate

# Build backend exec file.
FROM golang:1.16.12-alpine3.18 AS backend
WORKDIR /backend-build

RUN apk --no-cache add gcc musl-dev

COPY . .

RUN go build \
    -o sha \
    ./backend/bin/server/main.go

# Make workspace with above generated files.
FROM alpine:3.18 AS monolithic
WORKDIR /usr/local/sha

COPY --from=backend /backend-build/sha /usr/local/sha/
COPY --from=frontend /frontend-build/.output/public /usr/local/sha/frontend/dist

# Directory to store the data, which can be referenced as the mounting point.
RUN mkdir -p /var/opt/sha

ENTRYPOINT ["./sha"]
