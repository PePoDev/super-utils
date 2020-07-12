FROM debian:stretch-slim

# common tools
RUN apt update && apt install curl wget gnupg2 inetutils-ping rsync apt-transport-https -y

# database tools
RUN apt install mysql-client -y && \
    wget -qO - https://www.mongodb.org/static/pgp/server-4.2.asc | apt-key add - && \
    echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu bionic/mongodb-org/4.2 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-4.2.list && \
    apt update && apt install mongodb-org-tools -y

# install gcloud tools
RUN echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] http://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add - && \
    apt update && apt install google-cloud-sdk -y

# kubernetes
RUN curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - && \
    echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | tee -a /etc/apt/sources.list.d/kubernetes.list && \
    apt update && apt install -y kubectl

# devops tool
RUN echo "deb http://ppa.launchpad.net/ansible/ansible/ubuntu trusty main" | tee /etc/apt/sources.list.d/ansible.list && \
    apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 93C4A3FD7BB9C367 && \
    apt update && apt install ansible -y

# terraform
RUN wget https://releases.hashicorp.com/terraform/0.12.28/terraform_0.12.28_linux_amd64.zip && \
    unzip terraform_0.12.28_linux_amd64.zip && rm terraform_0.12.28_linux_amd64.zip && \
    mv terraform /usr/local/bin/

CMD ["echo https://github.com/pepodev/super-utils-docker"]
ENTRYPOINT [ "/bin/sh", "-c"]
