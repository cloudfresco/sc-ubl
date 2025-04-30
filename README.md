# sc-ubl

go get github.com/bufbuild/protovalidate-go

create buf.yaml

then run buf mod update

then create buf.gen.yaml 

then run buf generate - it will generate all proto files

-------------------------------------
Open Terminal

###
source sc-ubl.sh
go run cmd/main.go

Open another terminal
source sc-ubl.sh
chmod 755 sc-ublserver.sh
./sc-hubserver.sh

docker run -d -p 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one:latest

