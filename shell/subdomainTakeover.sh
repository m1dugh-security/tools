#!/usr/bin/env bash

printErr() {
    echo -e "[\033[0;31mFAILED\033[0m]\t$1" > /dev/stderr
}

printNotFound () {
    echo -e "[\033[0;33mNOT FOUND\033[0m]\t$1" > /dev/stderr
}

printSuccess () {
    echo -e "[\033[0;36mFOUND\033[0m]\t$1"
}

function extractDomain() {
    read cname
    domain=$(echo "$cname" | grep -Eo "\w+\.[a-z]{2,10}\.$" | sed -E "s/\.$//g")
    echo "$domain"
}

function checkAvailable() {
    if [ $# -le 0 ];then
        return 127
    fi
    domain=$1
    nslookup "$domain" &> /dev/null
    if [ $? -ne 0 ]; then
        return 0
    else
        return 1
    fi

}

# functions expecting the CNAME record
function checkGithub() {
    cname=$1
    body=$(curl -g "https://$cname" 2> /dev/null)
    if [[ "$body" =~ .*"There isn't a GitHub Pages site here.".* ]];then
        return 0
    else
        return 1
    fi
}

function checkAws () {
    cname=$1
    body=$(curl -g "https://$cname" 2> /dev/null)
    if [[ "$body" =~ .*"The specified bucket does not exist".* ]]; then
        return 0
    else
        return 1
    fi
}

function checkWordpress () {
    cname=$1
    body=$(curl -g "https://$cname" 2> /dev/null)
    if [[ "$body" =~ .*"Do you want to register".* ]]; then
        return 0
    else
        return 1
    fi
}

while read sub;do
    cname=$(dig $sub cname | grep "CNAME" | grep -v ";" | grep -Eo "(\w+\.)+[a-z]{2,10}\.$")
    if [ `printf "$cname" | wc -c` -gt 0 ]; then
        domain=$(echo "$cname" | extractDomain)
        checkAvailable $domain
        if [ $? -eq 0 ];then
            printSuccess "subdomain $sub linked to $cname can be taken over"
        elif [[ "$cname" = "github.io" ]] && (checkGithub $cname); then
            printSuccess "github: $sub -> $cname"
        elif [[ "$cname" =~ .*"s3.amazonaws.com" ]] && (checkAws $cname); then
            printSuccess "aws: $sub -> $cname"
        elif [[ "$cname" =~ .*"wordpress.com" ]] && (checkWordpress $cname); then
            printSuccess "wordpress: $sub -> $cname"
        else
            printErr $sub
        fi
    else
        printNotFound $sub
    fi
done
