FROM golang:1.14
ENV MYSQLTMPROOT urlooker.pass
ENV MYSQL_MAJOR 5.7

LABEL maintainer="710leo@gmail.com"

WORKDIR /app

RUN apt-key adv --keyserver ha.pool.sks-keyservers.net --recv-keys A4A9406876FCBD3C456770C88C718D3B5072E1F5 \
&& echo "deb http://repo.mysql.com/apt/debian/ jessie mysql-${MYSQL_MAJOR}" > /etc/apt/sources.list.d/mysql.list

RUN echo mysql-community-server mysql-community-server/root-pass password $MYSQLTMPROOT | debconf-set-selections \
&& echo mysql-community-server mysql-community-server/re-root-pass password $MYSQLTMPROOT | debconf-set-selections \
&& apt-get update && apt-get install -y mysql-server

COPY . .

RUN ./control build

VOLUME ["/var/lib/mysql"]

EXPOSE 1984

ENTRYPOINT ["./docker-entrypoint.sh"]