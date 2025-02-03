FROM golang:alpine AS setup

RUN apk update \
    && apk add --no-cache --update git  ca-certificates  gcc musl-dev curl \
    && update-ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files to download dependencies
COPY go.mod go.sum ./

RUN go mod download

# Copy the source code into the container
COPY . .

FROM setup AS tester

ARG USER=tester

RUN go install github.com/onsi/ginkgo/v2/ginkgo@v2.22.1

RUN apk add --update sudo && \
    adduser -D $USER && \
    echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER && \
    chmod 0440 /etc/sudoers.d/$USER && \
    chown -R ${USER}:${USER} .

USER $USER

RUN ginkgo --v \
    --race -r \
    --junit-report=testreports/report.xml \
    --cover --coverprofile=coverage.out

# -------- builder stage -------- #
FROM setup AS builder

ARG TARGETOS
ARG TARGETARCH
ARG GIT_TAG

# build binary
RUN CGO_CFLAGS="-D_LARGEFILE64_SOURCE" GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=1 \
    go build -v \
    -ldflags="-X 'github.com/sks/kihocche/pkg/constants.Version=${GIT_TAG}'" \
    -o /bin/kihocche .


# -------- prod stage -------- #
FROM alpine:latest

WORKDIR /app/kihocche

# create non root user
RUN apk update && apk upgrade --no-cache libcrypto3 && \
    apk --no-cache add ca-certificates git && \
    addgroup --gid 101 kihocche && \
    adduser -S --uid 101 --ingroup kihocche kihocche

# run as non root user
USER 101

# copy kihocche binary from build
COPY --from=builder /bin/kihocche /bin/

ENTRYPOINT ["/bin/kihocche"]

WORKDIR /data

CMD [ "serve" ]
