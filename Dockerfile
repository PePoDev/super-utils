FROM alpine:3

RUN apk --update add --no-cache -X http://dl-cdn.alpinelinux.org/alpine/edge/testing \
    curl git wget unzip iputils rsync openssh sshpass gnupg tar python3 py3-pip \
    mysql-client postgresql-client mongodb-tools redis ansible terraform helm kubectl && \
    pip3 install --upgrade pip

# Google Cloud SDK
RUN wget https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.zip && \
    unzip google-cloud-sdk.zip && rm google-cloud-sdk.zip && \
    mv google-cloud-sdk /usr/local/bin/

# AWS S
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

CMD ["echo see the document on https://github.com/PePoDev/super-utils"]
ENTRYPOINT [ "/bin/sh", "-c" ]
