#!/usr/bin/env bash

function installNodeModules() {
	echo
	if [ -d node_modules ]; then
		echo "============== node modules installed already ============="
	else
		echo "============== Installing node modules ============="
		npm install
	fi
	echo
}

function moveToDir() {
	echo
	if [ ! -d $1 ]; then
		echo "============== Folder $1 not present ============="
	else
		echo "============== Switching to $1 folder ============="
		cd $1
	fi
	echo
}

## Install node_modules
installNodeModules

## Moving to dir
moveToDir routes

## Install node_modules
# installNodeModules

## Clean wallet
echo "============== Removing existing ids in wallet ============="
echo
rm wallet/* -rf

## move back
echo "============== Moving back to medisot_sdk folder ============="
echo
cd ..

## Moving to dir
moveToDir fabtoken

## Install node_modules
installNodeModules

## Moving to dir
moveToDir ../performance-test

## Install node_modules
installNodeModules

# npm start

## move back
echo "============== Moving back to medisot_sdk folder ============="
echo
cd ..

#node app.js
#pm2 start app.js

