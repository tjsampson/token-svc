##########################################
# STEP 1 build binary in Build Stage Image
##########################################
FROM golang:alpine AS builder

LABEL maintainer="Troy Sampson <troy.j.sampson@gmail.com>"

# Build ARGS
ARG VERSION=0.0.0
ARG GIT_BRANCH=""
ARG GIT_COMMIT_AUTHOR=""
ARG BUILD_HOST=""
ARG BUILDER=""
ARG GIT_COMMIT=""
ARG BUILD_DATE=""
ARG PRERELEASE=""

# ENVIRONMENT VARIABLES
# set go modules on
ENV GO111MODULE=on

# Golang buildtime ldflags
ENV LDFLAGS=" -X main.BuildHost=${BUILD_HOST} \
    -X main.GitBranch=${GIT_BRANCH} \
    -X main.Builder=${BUILDER} \
    -X main.Version=${VERSION} \
    -X main.BuildDate=${BUILD_DATE} \
    -X main.GitCommit=${GIT_COMMIT} \
    -X main.GitCommitAuthor=${GIT_COMMIT_AUTHOR} \
    -X main.VersionPrerelease=${PRERELEASE} "

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
# tzdata is for timezone data
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# create appuser.
RUN adduser -D -g '' appuser

# set app working dir
WORKDIR $GOPATH/token-svc

# Copy The source assets from the CWD (project root) into the container WORKDIR ($GOPATH/token-svc)
COPY . .

# Verify the Go Modules (we are vendoring the go modules)
RUN go mod verify

# Build the golang binary
RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -ldflags "${LDFLAGS}" \
    -a -installsuffix cgo \
    -mod vendor \
    -o /go/bin/token-svc ./cmd/token-svc/


RUN touch /tmp/tokensvc.logs && chown appuser:appuser /tmp/tokensvc.logs


######################################
# STEP 2 build a smaller runtime image
######################################
FROM scratch

LABEL maintainer="Troy Sampson <troy.j.sampson@gmail.com>"

# Import assets from the build stage image
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/token-svc /go/bin/token-svc
COPY --from=builder /go/token-svc/migrations /go/bin/migrations
COPY --from=builder /tmp/tokensvc.logs /tmp/tokensvc.logs

# Use the unprivileged user (created in the build stage image)
USER appuser

WORKDIR /go/bin


# Set the entrypoint to the golang executable binary
ENTRYPOINT ["/go/bin/token-svc"]

# Expose the service ports (4000 for app and 4001 for metrics)
EXPOSE 4000
EXPOSE 4001

# Setup Container HealthCheck
HEALTHCHECK --interval=1m --timeout=2s --start-period=10s \
    CMD curl -f http://localhost:4000/health/ping || exit 1
