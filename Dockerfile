# Same version as our go.mod
FROM golang:1.23-alpine AS build

WORKDIR /svc

# Copy go.mod file across and download dependencies.
COPY go.* .
RUN go mod download

# Copy Go code into the container.
COPY cmd cmd
COPY internal internal

# Build the Go code.
RUN go build -o /bin/svc cmd/*.go

FROM alpine:latest AS prod

COPY --from=build /bin/svc /bin/svc

EXPOSE 8081

RUN adduser -S svcuser
RUN chown svcuser /bin/svc 
USER svcuser

ENTRYPOINT ["/bin/svc"]
