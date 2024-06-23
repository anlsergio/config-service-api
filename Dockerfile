#####################
### BUILDER STAGE ###
#####################
FROM golang:1.22.4-alpine AS builder
RUN apk update && \
    apk add --no-cache make
# TODO: install delve and create dev stage?
WORKDIR /src
# Copy Go dependencies definitions separately to take
# advantage of image layer caching
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENV GOOS=linux
ENV GOARCH=amd64
RUN make build

##################
### BASE STAGE ###
##################
FROM alpine:3.20.1 AS base
EXPOSE 80
WORKDIR /app
COPY --from=builder /src/build/* ./
COPY --from=builder /src/*.env ./
COPY --from=builder /src/entrypoint.sh .
ENTRYPOINT ["./entrypoint.sh"]

###################
### FINAL STAGE ###
###################
FROM base AS final
RUN addgroup -S config-service && \
    adduser -S config-service -G config-service
USER config-service
