#!/bin/sh

openssl req -config openssl.cnf -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
