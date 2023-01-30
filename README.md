Go plugins for the protocol compiler:

Protobuf

    $ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

Grpc

    $ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

Update your PATH so that the protoc compiler can find the plugins:

    $ export PATH="$PATH:$(go env GOPATH)/bin"

Run the `build.sh` file for executing proto files, generating swagger docs and for updating the environment variables

    $ source build.sh

There are two parts:-

Start the web server by running the following command

    $ go run cmd/main.go  web   

Open a new terminal,start the respected adapter based on your requirement by running

    $ go run cmd/main.go <adapter command-name>

Adapter Command Names:-
1. For NONEVM: `    nonevm`
2. For EVM: `    evm`


    Supported EVM based chains: arbitrum,fantom,ethereum,harmony,polygon,celo,optimism,xinfin,metis,avalanche,aurora,bsc,boba,heco,gnosis,moonriver,moonbeam,klaytn,bttc,iotex,tomochain,evmos,fuse,cronos,astar

3. For COSMOS:`   cosmos`


    Supported COSMOS based chains: axelar, akash, bandchain, cosmoshub, crescent, cryptoorgchain, injective, juno, kujira, kava, osmosis, secretnetwork, sifchain, umee, regen, stargaze, sentinel, persistence, irisnet, agoric, shentu, impacthub, emoney, sommelier, bostrom, gravitybridge, stride, assetmantle



