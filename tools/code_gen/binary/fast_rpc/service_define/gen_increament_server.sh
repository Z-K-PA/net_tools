#!/bin/bash

rm -rf ../output/*
mkdir -p ../output/main
mkdir -p ../output/message
../gen/gen_msg --template="../msg/obj_define.tpml" --in="./inc_server_msg.yaml" --out="../output/message/obj_define.go"
../gen/gen_msg --template="../msg/msg_define.tpml" --in="./inc_server_msg.yaml" --out="../output/message/msg_define.go"
../gen/gen_msg --template="../msg/msg_pack_unpack.tpml" --in="./inc_server_msg.yaml" --out="../output/message/msg_pack_unpack.go"
../gen/gen_msg --template="../msg/msg_parse.tpml" --in="./inc_server_msg.yaml" --out="../output/message/msg_parse.go"

../gen/gen_msg_api --template="../msg/msg_api.tpml" --in="./inc_server_api.yaml" --out="../output/message/msg_api.go"

go fmt ../output/message/obj_define.go
go fmt ../output/message/msg_define.go
go fmt ../output/message/msg_pack_unpack.go
go fmt ../output/message/msg_parse.go

go fmt ../output/message/msg_api.go

../gen/gen_server --template="../main/main.tpml" --in="./inc_server_app.yaml" --out="../output/main/main.go"
cp ../main/server_run.go ../output/main/server_run.go
cp ../main/service_init.go ../output/main/service_init.go

go fmt ../output/main/main.go
go fmt ../output/main/server_run.go
go fmt ../output/main/service_init.go