#!/bin/bash

CWD=$(cd $(dirname $0)/; pwd)
cd $CWD

usage()
{
	echo $"Usage: $0 {start|stop|restart|status|build|pack} <module>"
	exit 0
}

start_all()
{
	# http: 1984
	test -x urlooker-web && start web
	test -x urlooker-alarm && start alarm
	test -x urlooker-agent && start agent
}

start()
{
	mod=$1
	if [ "x${mod}" = "x" ]; then
		usage
		return
	fi

	if [ "x${mod}" = "xall" ]; then
		start_all
		return
	fi

	binfile=urlooker-${mod}

	if [ ! -f $binfile ]; then
		echo "file[$binfile] not found"
		exit 1
	fi


	if [ $(ps aux|grep -v grep|grep -v control|grep "$binfile" -c) -gt 0 ]; then
		echo "${mod} already started"
		return
	fi

	mkdir -p logs/$mod
	nohup $CWD/$binfile &> logs/${mod}/stdout.log &

	for((i=1;i<=15;i++)); do
		if [ $(ps aux|grep -v grep|grep -v control|grep "$binfile" -c) -gt 0 ]; then
			echo "${mod} started"
			return
		fi
		sleep 0.2
	done

	echo "cannot start ${mod}"
	exit 1
}

stop_all()
{

	test -x urlooker-web && stop web
	test -x urlooker-alarm && stop alarm
	test -x urlooker-agent && stop agent
}

stop()
{
	mod=$1
	if [ "x${mod}" = "x" ]; then
		usage
		return
	fi

	if [ "x${mod}" = "xall" ]; then
		stop_all
		return
	fi

	binfile=urlooker-${mod}

	if [ $(ps aux|grep -v grep|grep -v control|grep "$binfile" -c) -eq 0 ]; then
		echo "${mod} already stopped"
		return
	fi

	ps aux|grep -v grep|grep -v control|grep "$binfile"|awk '{print $2}'|xargs kill
	for((i=1;i<=15;i++)); do
		if [ $(ps aux|grep -v grep|grep -v control|grep "$binfile" -c) -eq 0 ]; then
			echo "${mod} stopped"
			return
		fi
		sleep 0.2
	done

	echo "cannot stop $mod"
	exit 1
}

restart()
{
	mod=$1
	if [ "x${mod}" = "x" ]; then
		usage
		return
	fi

	if [ "x${mod}" = "xall" ]; then
		stop_all
		start_all
		return
	fi

	stop $mod
	start $mod

	status
}

status()
{
	ps aux|grep -v grep|grep "urlooker"
}

build_one()
{
	mod=$1
	go build -o urlooker-${mod} modules/${mod}/${mod}.go
}

build()
{
	mod=$1
	if [ "x${mod}" = "x" ]; then
		build_one web
		build_one alarm
		build_one agent
		return
	fi

	build_one $mod
}

reload()
{
	mod=$1
	if [ "x${mod}" = "x" ]; then
		echo "arg: <mod> is necessary"
		return
	fi
	
	build_one $mod
	restart $mod
}

pack()
{
	v=$(date +%Y-%m-%d-%H-%M-%S)
	tar zcvf urlooker-$v.tar.gz control \
	supervisord.d \
	urlooker-web configs/web.yml \
	urlooker-alarm configs/alarm.yml \
	urlooker-agent configs/agent.yml
}

case "$1" in
	start)
		start $2
		;;
	stop)
		stop $2
		;;
	restart)
		restart $2
		;;
	status)
		status
		;;
	build)
		build $2
		;;
	reload)
		reload $2
		;;
	pack)
		pack
		;;
	*)
		usage
esac
