#!/bin/bash

# start mysql
mysqld_safe &
sleep 5s

mysql -h127.0.0.1 -uroot -purlooker.pass < sql/schema.sql

# start web
./control start web

# start alarm
./control start alarm

# start agent
./control start agent

# keep script in foreground
tail -f logs/*/stdout.log