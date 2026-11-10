# base
FROM golang:alpine AS base
RUN mkdir -p /opt/app
WORKDIR /opt/app
# 1. Copy ONLY go.mod and go.sum first
COPY go.mod go.sum ./

# 2. Download dependencies (this layer will now be cached!)
RUN go mod download

# 3. Copy the rest of the source code afterwards
COPY . .

# build
FROM base AS build
RUN apk add --no-cache build-base
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o bookmark-service cmd/api/main.go

# test
FROM base AS test-exec
ARG _outputdir="/test-output"
ARG COVERAGE_EXCLUDE
RUN mkdir -p ${_outputdir} && \
    go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1 && \
    if [ -n "${COVERAGE_EXCLUDE}" ]; then \
        grep -vE "${COVERAGE_EXCLUDE}" coverage.tmp > ${_outputdir}/coverage.out; \
    else \
        cp coverage.tmp ${_outputdir}/coverage.out; \
    fi && \
    go tool cover -html=${_outputdir}/coverage.out -o ${_outputdir}/coverage.html

# run test
FROM scratch AS test
ARG _outputdir="/test-output"
COPY --from=test-exec ${_outputdir}/coverage.out /
COPY --from=test-exec ${_outputdir}/coverage.html /

# final
FROM alpine:latest AS final
RUN addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=build /opt/app/bookmark-service .
COPY --from=build /opt/app/docs ./docs
COPY --from=build /opt/app/internal/migrations ./internal/migrations
USER app
CMD ["./bookmark-service"]