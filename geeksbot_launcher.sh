#!/bin/bash

until python /home/dusty/bin/geeksbot/geeksbot.py; do
	echo "Geeksbot shutdown with error: $?. Restarting..." >&2
	sleep 1
done
