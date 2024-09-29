#!/bin/sh
openssl genrsa -out ../ssl/cert.key 2048
openssl req -new -key ../ssl/cert.key -subj "/CN=$1" -sha256 | openssl x509 -req -days 365 -CA ../ssl/ca.crt -CAkey ../ssl/ca.key -set_serial "0x`openssl rand -hex 8`" > ../ssl/nck.crt
#sudo chmod 777 ssl/*