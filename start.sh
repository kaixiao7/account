#!/bin/bash

APP_NAME=account

usage() {
	echo "Usage: sh start.sh [start|stop|restart|status]"
	exit 1
}

is_exist() {
	pid=`ps -ef | grep $APP_NAME | grep -v grep | awk '{print $2}'`
	if [ -z "${pid}" ];then
		return 1
	else
		return 0
	fi
}

start() {
	is_exist
	if [ $? -eq "0" ];then
		echo "${APP_NAME} is already running. pid=${pid}."
	else
		#nohup ./${APP_NAME} --config account.yaml > /dev/null 2>&1 &
		nohup ./${APP_NAME} --config account.yaml > logs/account.log 2>&1 &
		echo "${APP_NAME} start success"
	fi
}

stop() {
	is_exist
	if [ $? -eq "0" ];then
		kill -9 ${pid}
	else
		echo "${APP_NAME} is not running."
	fi
}

status() {
	is_exist
	if [ $? -eq "0" ];then
		echo "${APP_NAME} is already running. pid=${pid}."
	else
		echo "${APP_NAME} is not running."
	fi
}

restart() {
	stop
	start
}

case "$1" in
	"start")
		start
		;;
	"stop")
		stop
		;;
	"status")
		status
		;;
	"restart")
		restart
		;;
	*)
		usage
		;;
esac
