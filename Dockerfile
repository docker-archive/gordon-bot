FROM progrium/busybox
MAINTAINER Jessica Frazelle <jess@docker.com>

ADD https://jesss.s3.amazonaws.com/binaries/gordon-bot /usr/local/bin/gordon-bot

RUN chmod +x /usr/local/bin/gordon-bot

ENTRYPOINT [ "/usr/local/bin/gordon-bot" ]
