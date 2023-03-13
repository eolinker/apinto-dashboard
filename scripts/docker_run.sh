#!/bin/sh

set -e


redisAddr=${REDIS_ADDR}
OLD_IFS="$IFS"
IFS=","
arr=($redisAddr)
IFS="$OLD_IFS"


echo "" > config.yml

echo -e "mysql:" > config.yml
echo -e "  user_name: ${MYSQL_USER_NAME}" >> config.yml
echo -e "  password: ${MYSQL_PWD}" >> config.yml
echo -e "  ip: ${MYSQL_IP}" >> config.yml
echo -e "  port: ${MYSQL_PORT}" >> config.yml
echo -e "  db: ${MYSQL_DB}" >> config.yml
echo -e "redis:" >> config.yml
echo -e "  user_name: ${REDIS_USER_NAME}" >> config.yml
echo -e "  password: ${REDIS_PWD}" >> config.yml
echo -e "  addr: " >> config.yml
for s in ${arr[@]}
do
echo -e "    - $s" >> config.yml
done
echo -e "port: 8080" >> config.yml
echo -e "user_center_url: ${USER_CENTER_URL}" >> config.yml
echo -e "error_log:" >> config.yml
echo -e "  dir: ${ERROR_DIR}" >> config.yml
echo -e "  file_name: ${ERROR_FILE_NAME}" >> config.yml
echo -e "  log_level: ${ERROR_LOG_LEVEL}" >> config.yml
echo -e "  log_expire: ${ERROR_EXPIRE}" >> config.yml
echo -e "  log_period: ${ERROR_PERIOD}" >> config.yml

./apserver




