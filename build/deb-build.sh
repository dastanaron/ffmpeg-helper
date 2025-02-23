#!/usr/bin/env bash

PROJECT_NAME="ffmpeg-helper"
PROJECT_VERSION="1.0"
PROJECT_BUILD_VERSION="1"

PACKAGE_NAME="ffmpeg-helper"

ARCHITECTURE="amd64"

INSTALL_PATH="/usr/local/bin"

MAINTAINER="dastanaron dastanaron@daserver.ru"
DESCRIPTION="FFmpeg Helper is a command-line tool that simplifies the use of FFmpeg by providing a user-friendly interface for executing predefined FFmpeg commands. It allows users to select from a list of commands, modify them, and execute them with ease."

BINARY_FILE_PATH=./ffmpeg-helper

PROJECT_PATH=${PROJECT_NAME}_${PROJECT_VERSION}-${PROJECT_BUILD_VERSION}_${ARCHITECTURE}


mkdir -p ${PROJECT_PATH}${INSTALL_PATH}

cp ${BINARY_FILE_PATH} ${PROJECT_PATH}${INSTALL_PATH}

mkdir ${PROJECT_PATH}/DEBIAN

size=$(du -ks ./ffmpeg-helper | grep -v DEBIAN | cut -f1 | xargs | sed -e 's/\ /+/g' | bc)

touch ${PROJECT_PATH}/DEBIAN/control

echo "Package: $PACKAGE_NAME" >> ${PROJECT_PATH}/DEBIAN/control
echo "Version: $PROJECT_VERSION" >> ${PROJECT_PATH}/DEBIAN/control
echo "Architecture: $ARCHITECTURE" >> ${PROJECT_PATH}/DEBIAN/control
echo "Maintainer: $MAINTAINER" >> ${PROJECT_PATH}/DEBIAN/control
echo "Description: $DESCRIPTION" >> ${PROJECT_PATH}/DEBIAN/control
echo "Installed-Size: $size" >> ${PROJECT_PATH}/DEBIAN/control
dpkg-deb --build --root-owner-group ${PROJECT_PATH}
