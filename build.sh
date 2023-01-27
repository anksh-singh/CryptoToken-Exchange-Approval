#!/bin/sh

#set -e
#COMMAND=$@
#
#echo ''
#maxTries=10
#while "" do  #TODO: Implement swap start logic with pre-checks
#    maxTries=$(($maxTries - 1))
#    sleep 3
#done
#echo
#if [ "$maxTries" -le 0 ]; then
#    echo >&2 'error: unable to contact mysql after 10 tries'
#    exit 1
#fi
#protoc command
cd pkg/grpc/proto
#
protoc --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative  *.proto


#Protoc by default generates the omitempty tag in json structure, but according to mobile client requirement , backend supposed to pass
#empty struct, for that the proto generated files need to be modified by removing omitempty tag
# the following command does the same as stated above
ls pb/*.pb.go | xargs -n1 -IX bash -c "sed -e  s/,omitempty// X > X.tmp && mv X{.tmp,}"

cd ../../..

export WEB_DATADOG_SERVICE="web-staging"
export DATADOG_SERVICE="bridge-allowance-Staging"
export NONEVM_CLUSTER="staging"
export NONEVM_DATADOG_SERVICE="nonevm-staging"
export PROXIES_ENDPOINT="https://staging.unifront.io"
export BRIDGE_DATADOG_SERVICE="bridge-staging"
export COSMOS_DATADOG_SERVICE="cosmos-staging"
export EVM_DATADOG_SERVICE="evm-staging"
export NONEVM_GRPC_ENDPOINT="localhost:8081"
export COSMOS_GRPC_ENDPOINT="localhost:8082"
export BRIDGE_GRPC_ENDPOINT="localhost:8085"
export EVM_GRPC_ENDPOINT="localhost:8083"
export LOG_ENCODING_FORMAT="console"


#to generate the swagger api docs
#`go env GOPATH`/bin/swag init -g  cmd/main.go
`go env GOPATH`/bin/swag init -g  cmd/main.go
#exec $COMMAND