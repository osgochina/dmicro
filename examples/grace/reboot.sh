#!/usr/bin/env bash


pid=`cat ./server.pid`
echo $pid
kill -USR2 $pid