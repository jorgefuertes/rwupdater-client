#!/usr/bin/env bash

if [[ "$0" != *scripts/*.sh ]]
then
	echo "Please, execute from project's root directory"
	exit 1
fi

source scripts/build_common.inc.sh

mkdir -p ./bin
rm -Rf ./bin/*

IFS='/'
echo $VER > ./bin/version.txt
for b in "${BUILD_LIST[@]}"
do
    echo "Building ${b}"
	read -ra THIS <<< "$b"
	OS=${THIS[0]}
	ARCH=${THIS[1]}
	DEST="./bin/${OS}/${ARCH}"
	mkdir -p "${DEST}"
	EXT=""
	[[ $OS == "windows" ]] && EXT=".exe"
    GOOS=$OS GOARCH=$ARCH go build -ldflags "${FLAGS}" \
		-o "${DEST}/${EXE_NAME}${EXT}" \
		retroupdater-client.go
    if [[ $? -ne 0 ]]
    then
        echo "Compilation error!"
        exit 1
    fi
done
IFS=' '

mkdir -p ./dist
rm -f ./dist/*
echo -n "Compressing releases..."
IFS='/'
for b in "${BUILD_LIST[@]}"
do
	read -ra THIS <<< "$b"
	OS=${THIS[0]}
	ARCH=${THIS[1]}
	SRC="./bin/${OS}/${ARCH}"
	FILE="${EXE_NAME}"
	if [[ $OS == "windows" ]]
	then
		zip -q "dist/${FILE}_${VER}-${OS}_${ARCH}.zip" "${SRC}/${FILE}.exe"
	else
		cp "${SRC}/${FILE}" ./dist/
		pushd dist &> /dev/null
		chmod +x "${FILE}"
		tar czf "${FILE}_${VER}-${OS}_${ARCH}.tar.gz" "${FILE}"
		rm "${FILE}"
		popd &> /dev/null
	fi

    if [[ $? -ne 0 ]]
    then
        echo "Compression error!"
        exit 1
    fi
done
IFS=' '
chmod -x dist/*
echo "OK"

echo -n "Publishing to local server..."
DEST="../retroupdater-server/files/client"
rm -Rf "${DEST}"
mkdir -p "${DEST}"
cp -a ./bin "${DEST}/."
cp -a ./dist "${DEST}/."
echo "OK"
echo -n "Publishing to remote server..."
SERVER_IP="updater.retrowiki.es"
SV_HOME="/home/retroserver"
rsync -az --delete --delete-excluded --exclude=".DS_Store" \
	$DEST/* $SERVER_IP:$SV_HOME/files/client/.
echo "OK"
