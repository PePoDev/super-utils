FROM alpine:3.7

# common tools
RUN apk add && apk install update curl git

# database tools
RUN apk add --no-cache mysql-client
RUN apk add --no-cache mongodb

# install gcloud tools

# kubernetes
ADD https://storage.googleapis.com/kubernetes-release/release/v1.6.4/bin/linux/amd64/kubectl /usr/local/bin/kubectl
RUN apk add --no-cache curl ca-certificates && \
    chmod +x /usr/local/bin/kubectl
# devops tool

# terraform

CMD ["echo https://github.com/pepodev/super-utils-docker"]
ENTRYPOINT [ "/bin/sh", "-c"]