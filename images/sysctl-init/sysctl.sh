#!/bin/bash
set -e

IFS=',' read -ra settings <<< "$SYSCTL"
for setting in "${settings[@]}"; do
    IFS='=' read -ra setting_split <<< "$setting"
    key=${setting_split[0]}
    value=${setting_split[1]}
    IFS='.' read -ra key_split <<< "$key"
    path=/proc/sys
    for i in "${key_split[@]}"; do
        path=$path/$i
    done
    echo $value > $path
done
