# argument for Go version
ARG GO_VERSION=1.14.5

# STAGE 1: building the executable
FROM golang:${GO_VERSION}-alpine AS build

# git required for go mod
RUN apk add --no-cache git

# certs
RUN apk --no-cache add ca-certificates

# Working directory will be created if it does not exist
WORKDIR /src

# We use go modules; copy go.mod and go.sum
COPY ./go.mod ./go.sum ./
RUN go mod download

# Import code
COPY ./ ./


# Build the executable
RUN CGO_ENABLED=0 go build \
	-installsuffix 'static' \
	-o /app .

# STAGE 2: build the container to run
FROM scratch AS final

# copy compiled app
COPY --from=build /app /app

# copy ca certs
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# run binary
ENTRYPOINT ["/app"]