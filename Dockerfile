# base
FROM golang:alpine AS base
RUN mkdir -p /opt/app
WORKDIR /opt/app
COPY . .
RUN go mod download

# build
FROM base AS build
RUN apk add --no-cache build-base
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o bookmark-service cmd/api/main.go

# test
FROM base AS test-exec
ARG _outputdir="/tmp/coverage"
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
ARG _outputdir="/tmp/coverage"
COPY --from=test-exec ${_outputdir}/coverage.out /
COPY --from=test-exec ${_outputdir}/coverage.html /

# final
FROM alpine:latest AS final
RUN addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=build /opt/app/bookmark-service .
COPY --from=build /opt/app/docs ./docs
USER app
CMD ["./bookmark-service"]