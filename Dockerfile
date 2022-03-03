FROM golang:1.17.7-alpine as build

# Install dependencies
WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download -x

# Build artifacts
COPY main.go main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/sample main.go

FROM scratch

# Copy our static executable
COPY --from=build /go/bin/sample /go/bin/sample

# Run the binary.
ENTRYPOINT ["/go/bin/sample"]
