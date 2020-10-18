# Generating Certficates for webrtc clients and turn server
# * Make sure the CA:true is set for certificate authority
# * Set DNS and IP for the alternative subject name for the turn server and client certs
# * Keep the serials unique for each certificate

# Server.
CA_NAME='CA'
openssl ecparam -name prime256v1 -genkey -noout -out "${CA_NAME}.pem"
openssl req -config "/etc/ssl/openssl.cnf" -key "${CA_NAME}.pem" -new -sha256 -subj '/C=NL' -out "${CA_NAME}.csr"
openssl x509 -req -in "${CA_NAME}.csr" -days 365 -signkey "${CA_NAME}.pem" -extfile "extcnfCA.ext" -sha256 -out "${CA_NAME}.pub.pem"

# Server.
SERVER_NAME='server'
openssl ecparam -name prime256v1 -genkey -noout -out "${SERVER_NAME}.pem"
openssl req -config "/etc/ssl/openssl.cnf" -key "${SERVER_NAME}.pem" -new -sha256 -subj '/C=NL' -out "${SERVER_NAME}.csr"
openssl x509 -req -in "${SERVER_NAME}.csr" -days 365 -CA "${CA_NAME}.pub.pem" -extfile "extcnf.ext" -CAkey "${CA_NAME}.pem" -set_serial '0xabcf' -sha256 -out "${SERVER_NAME}.pub.pem"

# Client.
CLIENT_NAME='clients'
openssl ecparam -name prime256v1 -genkey -noout -out "${CLIENT_NAME}.pem"
openssl req -config "/etc/ssl/openssl.cnf" -key "${CLIENT_NAME}.pem" -new -sha256 -subj '/C=NL' -out "${CLIENT_NAME}.csr"
openssl x509 -req -in "${CLIENT_NAME}.csr" -days 365 -CA "${CA_NAME}.pub.pem" -extfile "extcnf.ext" -CAkey "${CA_NAME}.pem" -set_serial '0xabcd' -sha256 -out "${CLIENT_NAME}.pub.pem"

CLIENT_NAME='clientc'
openssl ecparam -name prime256v1 -genkey -noout -out "${CLIENT_NAME}.pem"
openssl req -config "/etc/ssl/openssl.cnf" -key "${CLIENT_NAME}.pem" -new -sha256 -subj '/C=NL' -out "${CLIENT_NAME}.csr"
openssl x509 -req -in "${CLIENT_NAME}.csr" -days 365 -CA "${CA_NAME}.pub.pem" -extfile "extcnf.ext" -CAkey "${CA_NAME}.pem" -set_serial '0xabce' -sha256 -out "${CLIENT_NAME}.pub.pem"