FROM alpine:3.7

# common tools
RUN apk add && apk install update curl git wget unzip

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
RUN wget https://releases.hashicorp.com/terraform/0.12.28/terraform_0.12.28_linux_amd64.zip && \
    unzip terraform_0.12.28_linux_amd64.zip && rm terraform_0.12.28_linux_amd64.zip && \
    mv terraform /usr/local/bin/

CMD ["echo https://github.com/pepodev/super-utils-docker"]
ENTRYPOINT [ "/bin/sh", "-c"]