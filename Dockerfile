FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${GOARCH} go build -o /out/comms-service

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /out/comms-service /comms-service

EXPOSE 8000

USER nonroot

ENTRYPOINT ["/comms-service"]
