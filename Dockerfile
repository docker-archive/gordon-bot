FROM debian:jessie
MAINTAINER Jessica Frazelle <jess@docker.com>

RUN apt-get update && apt-get install -y ca-certificates

ADD https://jesss.s3.amazonaws.com/binaries/gordon-bot /usr/local/bin/gordon-bot

RUN chmod +x /usr/local/bin/gordon-bot

ENTRYPOINT [ "/usr/local/bin/gordon-bot" ]
