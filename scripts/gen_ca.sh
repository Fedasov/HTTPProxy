#!/bin/sh

openssl genrsa -out ../ssl/ca.key 2048
openssl req -new -x509 -days 365 -key ../ssl/ca.key -out ../ssl/ca.crt -subj "/CN=$1 proxy CA"
#sudo chmod 777 ssl/*