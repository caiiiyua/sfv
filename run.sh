#!/bin/bash

PROGRAM="sfver"

while true; do
	sleep 20
	PRO_NOW=`ps aux | grep $PROGRAM | grep -v grep | wc -l`
	if [[ $PRO_NOW -lt 1 ]]; then
		./$PROGRAM 2>/dev/null 1>&2 &
		date >> ./timeinfo.log 
		echo "sfver start" >> ./timeinfo.log
	fi

	PRO_STAT=`ps aux | grep $PROGRAM | grep T | grep -v grep | wc -l`

	if [[ $PRO_STAT -gt 0 ]]; then
		killall -9 $PROGRAM
		sleep 5
		./$PROGRAM 2>/dev/null 1>&2 &
		date >> ./timeinfo.log 
		echo "sfver start ..." >> ./timeinfo.log
	fi

done

exit 0