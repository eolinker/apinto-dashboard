#!/bin/sh


APP_NAME=apserver

ERR_FILE=../work/logs/stdout-apserver-$(date "+%Y%m%d-%H%M%S").log
pid=""
run() {

  echo "begin start"

  pid=`ps -ef|grep $APP_NAME|grep -v grep|awk '{print $2}' `
  if [ ! -z "${pid}" ]; then
    echo "${APP_NAME} is already running. pid=${pid} ."
  else
    mkdir -p ../work/logs
    nohup ./$APP_NAME > ${ERR_FILE} 2>&1 &

    #等待三秒在判断程序有没有成功启动
    sleep 3
    pid=`ps -ef|grep $APP_NAME|grep -v grep|awk '{print $2}' `
    if [ ! -z "${pid}" ]; then
      echo "${APP_NAME} start success"
    else
      #启动失败
      echo "failed to start, please check ${ERR_FILE}"
    fi

 fi

}

stop() {
    echo "begin stop"
    pid=`ps -ef|grep $APP_NAME|grep -v grep|awk '{print $2}' `
     if [ ! -z "${pid}" ]; then
        echo "kill -9 $pid"
        kill -9 $pid
        echo "stop success"
     else
        echo "${APP_NAME} is not running"
     fi
}


case "$1" in
    start)
        run
    ;;
    stop)
        stop
    ;;
    reload|restart|force-reload)
        stop
        run
        echo 222
    ;;
    *)
        echo "Usage: $0 {start|stop|reload|restart|force-reload} " 1>&2
    ;;
esac
