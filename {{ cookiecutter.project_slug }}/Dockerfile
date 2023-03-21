############################
# STEP 1 build the binary
############################
FROM golang:1.19 AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPROXY=https://goproxy.cn,direct

# Move to working directory /build
WORKDIR /build

# # Copy and download dependency using go mod
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application and debugger
RUN make build

############################
# STEP 2 build a small image
############################
FROM debian:stable-slim

ENV TZ=Asia/Shanghai

RUN sed -i 's/deb.debian.org/mirrors.tencentyun.com/g' /etc/apt/sources.list \
  && sed -i 's/security.debian.org/mirrors.tencentyun.com/g' /etc/apt/sources.list \
  && apt-get update && apt-get install --no-install-recommends -y ca-certificates tzdata \
  && apt-get autoremove -y && apt-get autoclean -y && rm -rf /var/lib/apt/lists/*

COPY --from=builder /build/bin/{{ cookiecutter.project_slug }} /bin/{{ cookiecutter.project_slug }}-cli

# Command to run the executable
CMD ["{{ cookiecutter.project_slug }}-cli", "run"]
