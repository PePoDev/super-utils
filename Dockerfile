# Start from the official Golang image
FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/super-utils .
RUN apk --update add --no-cache -X http://dl-cdn.alpinelinux.org/alpine/edge/testing \
    curl git wget unzip iputils rsync openssh sshpass gnupg tar python3 py3-pip gzip jq cmake \
    mysql-client postgresql-client mongodb-tools redis ansible terraform helm kubectl && \
    pip3 install --upgrade pip

# Google Cloud SDK
RUN wget https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.zip && \
    unzip google-cloud-sdk.zip && rm google-cloud-sdk.zip && \
    mv google-cloud-sdk /usr/local/bin/

# AWS SDK and Lambda Runtime API
RUN pip3 install --no-cache-dir awscli && rm -rf /var/cache/apk/*

# MSSQL Tools
ARG MSSQL_VERSION=17.5.2.1-1
ENV MSSQL_VERSION=${MSSQL_VERSION}
RUN curl -O https://download.microsoft.com/download/e/4/e/e4e67866-dffd-428c-aac7-8d28ddafb39b/msodbcsql17_${MSSQL_VERSION}_amd64.apk && \
    curl -O https://download.microsoft.com/download/e/4/e/e4e67866-dffd-428c-aac7-8d28ddafb39b/mssql-tools_${MSSQL_VERSION}_amd64.apk && \
    curl -O https://download.microsoft.com/download/e/4/e/e4e67866-dffd-428c-aac7-8d28ddafb39b/msodbcsql17_${MSSQL_VERSION}_amd64.sig && \
    curl -O https://download.microsoft.com/download/e/4/e/e4e67866-dffd-428c-aac7-8d28ddafb39b/mssql-tools_${MSSQL_VERSION}_amd64.sig && \
    curl https://packages.microsoft.com/keys/microsoft.asc  | gpg --import - && \
    gpg --verify msodbcsql17_${MSSQL_VERSION}_amd64.sig msodbcsql17_${MSSQL_VERSION}_amd64.apk && \
    gpg --verify mssql-tools_${MSSQL_VERSION}_amd64.sig mssql-tools_${MSSQL_VERSION}_amd64.apk && \
    echo y | apk add --allow-untrusted msodbcsql17_${MSSQL_VERSION}_amd64.apk mssql-tools_${MSSQL_VERSION}_amd64.apk && \
    rm -f msodbcsql*.sig msodbcsql*.apk mssql-tools*.sig mssql-tools*.apk
ENV PATH=$PATH:/opt/mssql-tools/bin

RUN apk add --no-cache \
	arping busybox mii-tool tcpdump tcptraceroute traceroute tshark \
	awk cut diff find grep sed vim coreutils \
	curl wget \
	bind-tools iproute2 net-tools mtr iputils iperf3 ethtool nmap \
	lftp rsync openssh-client socat netcat-openbsd apache2-utils \
	mysql-client postgresql-client git gzip cpio tar \
	bash bird bridge-utils busybox-extras calicoctl conntrack-tools ctop dhcping drill file fping httpie iftop iperf ipset iptables iptraf-ng ipvsadm jq libc6-compat liboping net-snmp-tools netgen nftables ngrep nmap-nping openssl py-crypto py2-virtualenv python2 scapy strace termshark util-linux websocat

    EXPOSE 8080

CMD ["echo see the document on https://github.com/PePoDev/super-utils"]
ENTRYPOINT [ "/bin/sh", "-c" ]
