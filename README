# Super Utils Docker Image

Super utils docker images is the containerize with many utility for multiple purpose

## Getting Start

Get this image in local machine

```bash
docker pull pepodev/super-utils

# mysqldump
docker run --rm pepodev/super-utils "mysqldump --help"

# gsutil // copy local file to google cloud storage
# working directory is /opt/ mount folder in this path to easy use as volume
docker run --rm -v ./service-account:./service-account pepodev/super-utils "gcloud activate ./service-account/gcp && gsutils -m copy -r ./dir/ gs://some-bucket"
```

## List of avaliables utility

- ping, curl, wget, rync, scp, ssh
- mysql-client (mysql, mysqldump)
- gcloud, gsutil
- kubectl
- ansible
- terraform
