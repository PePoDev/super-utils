# Super Utils

Super utils is the utility containerized for multiple purpose

## Getting Started

Get this image in local machine

```bash
docker pull quay.io/pepodev/super-utils

# mysqldump
docker run --rm quay.io/pepodev/super-utils "mysqldump --help"

# gsutil // copy local file to google cloud storage
# working directory is /opt/ mount folder in this path to easy use as volume
docker run --rm -v ./service-account:./service-account quay.io/pepodev/super-utils "gcloud activate ./service-account/gcp && gsutils -m copy -r ./dir/ gs://some-bucket"
```

## List of utility

- curl, wget, unzip, git, ping, rsync, ssh, scp, tar, gzip
- mysql-client
- mysqldump
- postgresql-client
- sqlcmd (MSSQL client)
- mongodb-client
- redis
- gcloud, gsutil
- aws-cli
- kubectl
- helm
- ansible
- terraform
- python, pip
