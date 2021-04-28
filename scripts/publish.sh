#!/usr/bin/env bash

SCRIPTS=$(dirname $0)
SERVER_IP="server.abadiaretro.com"
SV_HOME="/home/retroserver"
SV_SERVER_NAME="retroserver"
CLIENT_HOME="${SV_HOME}/files/client"

$SCRIPTS/build.sh
if [[ $? -ne 0 ]]
then
	echo "Building error!"
	exit 1
fi

ssh root@$SERVER_IP <<-CMD
	mkdir -p $CLIENT_HOME
	rm -f $CLIENT_HOME/*
CMD

scp $SCRIPTS/../bin/* root@$SERVER_IP:$CLIENT_HOME/.

ssh root@$SERVER_IP <<-CMD
	chown -R retroserver:retroserver $CLIENT_HOME
CMD

exit 0
