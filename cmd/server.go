package main

import (
	"distributedKV/kvraft"
	"distributedKV/net"
	"distributedKV/persist"
	"flag"
	"fmt"
	"net/http"
	"net/rpc"
	"strconv"
)

func main() {
	flag.Parse()
	args := flag.Args()
	var servers []net.ClientEnd
	for i := 0; i < len(args)-1; i++ {
		servers = append(servers, net.MakeNetClientEnd(args[i]))
	}
	me, _ := strconv.Atoi(args[len(args)-1])

	// 默认触发快照大小为 512M 后期改成可配置
	maxsaftstate := 1024 * 1024 * 512
	server := kvraft.StartKVServer(servers, me, persist.MakeFilePersister(me), maxsaftstate)
	rpc.Register(server) // 不注册也没关系? 不行: lab2不注册是因为在本地
	rpc.Register(server.Rf)
	fmt.Println("start server", me)
	rpc.HandleHTTP()
	http.ListenAndServe(args[me], nil)
}
