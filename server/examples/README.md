# gRPC client examples for `bhlindex`


## Ruby

Ruby client

```bash
cd ./ruby
bundle
grpc_tools_ruby_protoc -I ../../../protob --ruby_out=lib --grpc_out=lib \
../../../protob/protob.proto
ruby bin/bhl_client
```
