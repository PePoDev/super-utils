FROM alpine:3.7

# common tools
RUN apk add -y

# database tools
RUN apk add --no-cache mysql-client

# install gcloud tools


# kubernetes


# devops tool


# terraform

CMD ["echo https://github.com/pepodev/super-utils-docker"]
ENTRYPOINT [ "/bin/sh", "-c"]
