# 生成rpc
###
protoc --proto_path=. --micro_out=./pbGo --go_out=./pbGo demo.proto

# 注入
#protoc-go-inject-tag -input=./demo.pb.go

# 生成swagger
protoc --swagger_out=logtostderr=true:./pbJson demo.proto
# 去掉omitempty
ls pbGo/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'

# 合并
swagger mixin \
pbJson/demo.swagger.json \
-o swagger.json
