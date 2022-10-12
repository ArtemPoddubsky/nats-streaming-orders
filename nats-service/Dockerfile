FROM golang:1.19 AS builder
COPY . /nats-service
WORKDIR /nats-service
RUN make

FROM golang:1.19
COPY --from=builder /nats-service ./
CMD ./app