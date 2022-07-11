# Introduction

This is a rough example on how to build an audit log. I expect to expand on it in the future.

## Building

Here is the script I use:

```
#!/bin/bash

echo "Building TiDB";
cd ~/go/src/github.com/morgo/tidb
make server

cd plugin/s3audit
echo "Building Plugin"
pluginpkg -pkg-dir . -out-dir .
cd ../..

killall -9 tidb-server pd-server tikv-server etcd
rm -rf /tmp/tidb*

./bin/tidb-server -plugin-dir plugin/s3audit -plugin-load s3audit-1
```

The plugin must be compiled with the same source code as the server, which is why it's easier to always do these 2 steps together.

You can get the required pluginpkg with:
```
cd ~/go/src/github.com/morgo/tidb/cmd/pluginpkg
go build 
```