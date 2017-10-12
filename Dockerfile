FROM scratch
MAINTAINER Christian Sullivan <cs@bleve.io>

COPY build/kobs /kobs

CMD ["/kobs"]
