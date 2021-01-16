FROM golang:1.15.2 AS build

LABEL maintainer="Oladiy"

ADD . /opt/app_forum/
WORKDIR /opt/app_forum/
RUN go build ./cmd/main.go

FROM ubuntu:20.04 AS release

RUN apt-get -y update && apt-get install -y locales gnupg2
RUN locale-gen en_US.UTF-8
RUN update-locale LANG=en_US.UTF-8

ENV PGVER 12
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update -y && apt-get install -y postgresql-$PGVER postgresql-contrib

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -E UTF8 -O docker Forum &&\
    /etc/init.d/postgresql stop

# PostgreSQL port
EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

COPY --from=build /opt/app_forum/main /usr/bin/

# Server port
EXPOSE 5000

CMD service postgresql start && psql -h localhost -U docker -d Forum -p 5432 -a -q -f ./init/init_postgre.sql && main
