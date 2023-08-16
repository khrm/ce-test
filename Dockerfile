FROM golang as builder

WORKDIR /go/src/github.com/khrm/ce-test
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o ce-test

FROM alpine

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/khrm/ce-test/ce-test /ce-test

# Run the web service on container startup.
CMD ["/ce-test"]
