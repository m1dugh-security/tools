#!/usr/bin/env bash

set -eu

printInfo () {
    echo -e "[\033[0;36mINFO\033[0m]\t$1"
}

ROOT_FOLDER="./recon"
SETTINGS=""
if [ $# -eq 2 ]; then
    ROOT_FOLDER="$1"
    SETTINGS="$2"
elif [ $# -eq 1 ]; then
    ROOT_FOLDER="$1"
else
    echo "Usage: recon <folder> [settings]"
    exit -1
fi

printInfo "starting program browser"
if [ -n "$SETTINGS" ]; then 
    filesetup -o "$ROOT_FOLDER" -settings "$SETTINGS"
else
    filesetup -o "$ROOT_FOLDER"
fi


for platform in `ls $ROOT_FOLDER`; do
    for program in `ls $ROOT_FOLDER/$platform`; do
        echo -e "[\033[0;36mSTARTING\033[0m]\t$platform-$program"
        (
            path="$ROOT_FOLDER/$platform/$program/"
            {
                cat "$path/domains.txt"
                subfinder -dL "$path/domains.txt"
            } 2> /dev/null | tee "$path/subdomains.txt" | while read domain; do
                    printInfo "$program:SUBDOMAIN $domain"
                done && \
                sort -u "$path/subdomains.txt" | tee "$path/subdomains.txt" &> /dev/null
            {
                cat "$path/urls.txt"
                cat "$path/subdomains.txt" | httprobe 
            } 2> /dev/null | tee "$path/urls.txt" | while read url; do
                    printInfo "$program:URL $url"
                done && \
                sort -u "$path/urls.txt" | tee "$path/urls.txt" &> /dev/null
        ) &
    done
done

wait

