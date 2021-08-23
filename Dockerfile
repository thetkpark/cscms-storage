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
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=client-builder /app/build ./client/build
COPY --from=server-builder /app/main ./
CMD ["/app/main"]
