# Base images follow RFD 0070: AppCo Go as the builder, BCI nano for the final image.
FROM --platform=$BUILDPLATFORM registry.rancher.com/rancher/appco-go:1.26.10 AS builder

ARG TARGETOS TARGETARCH

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -o llm-mock ./cmd/main.go

FROM registry.suse.com/bci/bci-nano:16.0

WORKDIR /app
COPY --from=builder /app/llm-mock .
EXPOSE 8083
ENTRYPOINT ["./llm-mock"]
