FROM alpine:3.7

# common tools
RUN apk upgrade && apk update && apk add curl git wget unzip

# database tools
RUN apk add  mysql-client && \
    apk add  mongodb && \
    apk add postgresql-client

# install gcloud tools
RUN wget https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.zip &&\
    unzip google-cloud-sdk.zip &&\
    rm google-cloud-sdk.zip
# kubernetes
RUN curl -L -o /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.8.5/bin/linux/amd64/kubectl && \
    chmod +x /usr/bin/kubectl

# devops tool
RUN apk add ansible
# terraform
RUN wget https://releases.hashicorp.com/terraform/0.12.28/terraform_0.12.28_linux_amd64.zip && \
    unzip terraform_0.12.28_linux_amd64.zip && rm terraform_0.12.28_linux_amd64.zip && \
    mv terraform /usr/local/bin/

CMD ["echo https://github.com/pepodev/super-utils-docker"]
ENTRYPOINT [ "/bin/sh", "-c"]