FROM alpine

RUN apk update && apk upgrade && apk add --no-cache -X http://dl-cdn.alpinelinux.org/alpine/edge/testing \
    curl git wget unzip iputils rsync openssh sshpass gnupg tar \
    mysql-client postgresql-client mongodb-tools ansible terraform helm kubectl

RUN wget https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.zip && \
    unzip google-cloud-sdk.zip && rm google-cloud-sdk.zip && \
    mv google-cloud-sdk /usr/local/bin/

CMD ["echo see the document on https://github.com/pepodev/super-utils-docker"]
ENTRYPOINT [ "/bin/sh", "-c" ]