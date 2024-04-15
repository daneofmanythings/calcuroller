#!/usr/bin/env bash

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# wipe *ALL* local certs
rm ./*.pem

# Generating certs for a local Content Authority and Server

country="US"
state="California"
local="Sunnyvale"

# Local Content Authority metadata.
ca_org="nunya"
ca_org_unit="business"
ca_org_domain="www.nunyabusiness.com"
ca_org_email="nunya@business.com"

# local Server Metadata
serv_org="calcuroller"
serv_unit="calc_server"
serv_domain="*.calcuroller.com"
serv_email="roller@calcuroller.com"

# 1. Generate CA's private key and self-signed certificate
openssl \
  req \
  -x509 \
  -newkey rsa:4096 \
  -days 365 \
  -nodes \
  -keyout "$CURRENT_DIR"/ca-key.pem \
  -out "$CURRENT_DIR"/ca-cert.pem \
  -subj "/C=$country/ST=$state/L=$local/O=$ca_org/OU=$ca_org_unit/CN=$ca_org_domain/emailAddress=$ca_org_email"

echo "CA's self-signed certificate"
openssl x509 -in ca-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl \
  req \
  -newkey rsa:4096 \
  -nodes \
  -keyout "$CURRENT_DIR"/server-key.pem \
  -out "$CURRENT_DIR"/server-req.pem \
  -subj "/C=$country/ST=$state/L=$local/O=$serv_org/OU=$serv_unit/CN=$serv_domain/emailAddress=$serv_email"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl \
  x509 \
  -req \
  -in "$CURRENT_DIR"/server-req.pem \
  -days 60 \
  -CA ca-cert.pem \
  -CAkey ca-key.pem \
  -CAcreateserial \
  -out "$CURRENT_DIR"/server-cert.pem \
  -extfile server-ext.cnf

echo "Server's signed certificate"
openssl x509 -in "$CURRENT_DIR"/server-cert.pem -noout -text
