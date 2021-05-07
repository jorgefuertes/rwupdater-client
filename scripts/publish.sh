#!/usr/bin/env bash

SCRIPTS=$(dirname $0)
SERVER_IP="updater.retrowiki.es"
SV_HOME="/home/retroserver"
SV_SERVER_NAME="retroserver"
CLIENT_HOME="${SV_HOME}/files/client"
LOCAL_SERVER_DIR="${SCRIPTS}/../../retroupdater-server/files/client"

$SCRIPTS/build.sh
if [[ $? -ne 0 ]]
then
	echo "Building error!"
	exit 1
fi

# Local
mkdir -p $LOCAL_SERVER_DIR/bin
mkdir -p $LOCAL_SERVER_DIR/dist
rm -f $LOCAL_SERVER_DIR/bin/*
rm -f $LOCAL_SERVER_DIR/dist/*

cp $SCRIPTS/../bin/* $LOCAL_SERVER_DIR/bin/.
cp $SCRIPTS/../dist/* $LOCAL_SERVER_DIR/dist/.

# Remote
ssh root@$SERVER_IP <<-CMD
	mkdir -p $CLIENT_HOME
	rm -Rf $CLIENT_HOME/*
	mkdir $CLIENT_HOME/bin
	mkdir $CLIENT_HOME/dist
CMD

rsync -avz --delete --delete-excluded --exclude=".DS_Store" \
	$SCRIPTS/../dist/* root@$SERVER_IP:$CLIENT_HOME/dist/.
rsync -avz --delete --delete-excluded --exclude=".DS_Store" \
	$SCRIPTS/../bin/* root@$SERVER_IP:$CLIENT_HOME/bin/.

ssh root@$SERVER_IP <<-CMD
	chown -R retroserver:retroserver $CLIENT_HOME
CMD

exit 0
