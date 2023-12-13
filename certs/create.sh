#!/bin/bash

DEST_DIR="./x509"

mkdir -p $DEST_DIR

# Create the server CA certs.
openssl req -x509                                     \
  -newkey rsa:4096                                    \
  -nodes                                              \
  -days 3650                                          \
  -keyout $DEST_DIR/ca_key.pem                        \
  -out $DEST_DIR/ca_cert.pem                          \
  -subj /C=US/ST=CA/L=SVL/O=gRPC/CN=grpc-server_ca/   \
  -config ./openssl.cnf                               \
  -extensions grpc_ca                                 \
  -sha256

# Generate a server cert.
openssl genrsa -out $DEST_DIR/server_key.pem 4096

openssl req -new                                        \
  -key $DEST_DIR/server_key.pem                         \
  -days 3650                                            \
  -out $DEST_DIR/server_csr.pem                         \
  -subj /C=US/ST=CA/L=SVL/O=gRPC/CN=grpc-server_moon/   \
  -config ./openssl.cnf                                 \
  -reqexts grpc_server

openssl x509 -req                     \
  -in $DEST_DIR/server_csr.pem        \
  -CAkey $DEST_DIR/ca_key.pem         \
  -CA $DEST_DIR/ca_cert.pem           \
  -days 3650                          \
  -set_serial 1000                    \
  -out $DEST_DIR/server_cert.pem      \
  -extfile ./openssl.cnf              \
  -extensions grpc_server             \
  -sha256

openssl verify -verbose -CAfile $DEST_DIR/ca_cert.pem  $DEST_DIR/server_cert.pem

rm $DEST_DIR/*_csr.pem