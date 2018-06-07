#!/bin/bash

until python -m /home/dusty/bin/geeksbot_dev; do
	echo "Geeksbot shutdown with error: $?. Restarting..." >&2
	sleep 1
done
