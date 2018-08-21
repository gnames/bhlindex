# gRPC client examples for `bhlindex`


## Ruby

Ruby client

```bash
cd ./ruby
bundle
grpc_tools_ruby_protoc -I ../../../protob --ruby_out=lib --grpc_out=lib \
../../../protob/protob.proto

# To get version of the server
ruby bin/version

# To collect names statistics
ruby bin/pages

# To collect most prevalent clades for volumes
ruby bin/context

# To collect most prevalent class for found names
ruby bin/classes
```
