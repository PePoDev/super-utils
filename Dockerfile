FROM alpine:3.7

# common tools
RUN apk update && apk upgrade && apk add curl git wget unzip iputils rsync openssh

# database cli
RUN apk add mysql-client postgresql-client mongodb-tools

# gcloud-sdk
RUN wget https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.zip && \
    unzip google-cloud-sdk.zip && rm google-cloud-sdk.zip && \
    mv google-cloud-sdk /usr/local/bin/

# kubernetes
RUN curl -L -o /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.8.5/bin/linux/amd64/kubectl && \
    chmod +x /usr/bin/kubectl

# ansible
RUN apk add ansible

# terraform
RUN wget https://releases.hashicorp.com/terraform/0.12.28/terraform_0.12.28_linux_amd64.zip && \
    unzip terraform_0.12.28_linux_amd64.zip && rm terraform_0.12.28_linux_amd64.zip && \
    mv terraform /usr/local/bin/

CMD ["/bin/sh", "-c", "echo see the document on https://github.com/pepodev/super-utils-docker"]
