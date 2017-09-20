#!/bin/bash
set -e

if [ "$APPEND" == "true" ]; then
    echo $SEARCH >> /etc/resolv.conf
else
    echo -e "$SEARCH\n$(cat /etc/resolv.conf)" > /etc/resolv.conf
fi
