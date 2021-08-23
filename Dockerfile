FROM node:14-alpine as client-builder
WORKDIR /app
COPY ./client/package.json ./
COPY ./client/yarn.lock ./
RUN yarn
COPY ./client ./
RUN yarn build


FROM golang:alpine as server-builder
RUN apk add build-base
WORKDIR /app
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o ./server ./cmd/server/main.go
RUN go build -o ./cleaner ./cmd/cleaner/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=client-builder /app/build ./client/build
COPY --from=server-builder /app/server ./
COPY --from=server-builder /app/cleaner ./
CMD ["/app/main"]
