# Build stage
FROM golang:alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations for target platform
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o super-utils main.go

# Final stage
FROM alpine:latest

ARG TARGETOS
ARG TARGETARCH

LABEL maintainer="PePoDev"
LABEL description="Super Utils - A comprehensive DevOps utilities container"

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/super-utils .

# Install base packages and tools in a single layer
RUN apk --no-cache add --update \
    # Base utilities
    bash curl wget git unzip tar gzip jq vim \
    # Network tools
    bind-tools iproute2 net-tools iputils openssh-client \
    netcat-openbsd tcpdump nmap mtr iperf3 socat \
    # Python and pip
    python3 py3-pip \
    # Database clients
    mysql-client postgresql-client mongodb-tools redis \
    # DevOps tools
    ansible helm kubectl aws-cli \
    # Additional utilities
    rsync sshpass gnupg apache2-utils coreutils \
    # Advanced network tools
    mii-tool tcptraceroute traceroute tshark \
    lftp cpio \
    bird bridge-utils busybox-extras conntrack-tools \
    drill file fping httpie iftop ipset iptables iptraf-ng \
    ipvsadm libc6-compat net-snmp-tools nftables ngrep \
    openssl scapy strace util-linux websocat \
    # Additional requested tools
    ctop dhcping ethtool iperf liboping nmap-nping termshark && \
    # Install calicoctl separately (binary download) - multi-arch support
    wget -O /usr/local/bin/calicoctl https://github.com/projectcalico/calico/releases/latest/download/calicoctl-linux-${TARGETARCH} && \
    chmod +x /usr/local/bin/calicoctl && \
    # Clean up
    rm -rf /var/cache/apk/* /tmp/*

# Install Terraform
ARG TERRAFORM_VERSION=1.13.4
RUN wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_${TARGETARCH}.zip && \
    unzip terraform_${TERRAFORM_VERSION}_linux_${TARGETARCH}.zip && \
    rm terraform_${TERRAFORM_VERSION}_linux_${TARGETARCH}.zip && \
    mv terraform /usr/bin/terraform

# Install Google Cloud SDK
RUN wget -q https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.tar.gz && \
    tar -xzf google-cloud-sdk.tar.gz && \
    rm google-cloud-sdk.tar.gz && \
    mv google-cloud-sdk /usr/local/ && \
    /usr/local/google-cloud-sdk/install.sh --quiet --usage-reporting=false --path-update=true && \
    ln -s /usr/local/google-cloud-sdk/bin/* /usr/local/bin/ && \
    rm -rf /root/.cache

# Install MSSQL Tools (only for amd64 architecture as ARM64 is not officially supported)
ARG MSSQL_VERSION=17.5.2.1-1
ENV MSSQL_VERSION=${MSSQL_VERSION}
ENV PATH=$PATH:/opt/mssql-tools/bin
RUN if [ "$TARGETARCH" = "amd64" ]; then \
    curl -sSL https://packages.microsoft.com/keys/microsoft.asc | gpg --import - && \
    curl -sSL -O https://download.microsoft.com/download/e/4/e/e4e67866-dffd-428c-aac7-8d28ddafb39b/msodbcsql17_${MSSQL_VERSION}_amd64.apk && \
    curl -sSL -O https://download.microsoft.com/download/e/4/e/e4e67866-dffd-428c-aac7-8d28ddafb39b/mssql-tools_${MSSQL_VERSION}_amd64.apk && \
    curl -sSL -O https://download.microsoft.com/download/e/4/e/e4e67866-dffd-428c-aac7-8d28ddafb39b/msodbcsql17_${MSSQL_VERSION}_amd64.sig && \
    curl -sSL -O https://download.microsoft.com/download/e/4/e/e4e67866-dffd-428c-aac7-8d28ddafb39b/mssql-tools_${MSSQL_VERSION}_amd64.sig && \
    gpg --verify msodbcsql17_${MSSQL_VERSION}_amd64.sig msodbcsql17_${MSSQL_VERSION}_amd64.apk && \
    gpg --verify mssql-tools_${MSSQL_VERSION}_amd64.sig mssql-tools_${MSSQL_VERSION}_amd64.apk && \
    echo y | apk add --allow-untrusted msodbcsql17_${MSSQL_VERSION}_amd64.apk mssql-tools_${MSSQL_VERSION}_amd64.apk && \
    rm -f msodbcsql*.sig msodbcsql*.apk mssql-tools*.sig mssql-tools*.apk; \
    else \
    echo "MSSQL Tools not available for $TARGETARCH architecture"; \
    fi

EXPOSE 8080

# Run the application
ENTRYPOINT [ "/bin/sh", "-c" ]
CMD ["/app/super-utils"]
