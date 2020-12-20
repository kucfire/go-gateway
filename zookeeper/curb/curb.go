package main

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var (
	addr = "127.0.0.1:2181"
	host = []string{"127.0.0.1:2181"}
)

func main() {
	// 链接zookeeper
	conn, _, err := zk.Connect(host, 5*time.Second)
	if err != nil {
		panic(err)
	}
	// fmt.Println(conn)

	// 增
	if _, err := conn.Create("/test_tree2", []byte("tree_connect"), 0, zk.WorldACL(zk.PermAll)); err != nil {
		// panic(err)
		fmt.Println("create error : ", err)
	}

	// 查
	nodeValue, dStat, err := conn.Get("/test_tree2")
	if err != nil {
		fmt.Println("get error : ", err)
		return
	}
	fmt.Println("nodeValue : ", string(nodeValue))

	// 改
	if _, err = conn.Set("/test_tree2", []byte("tree_connect_set"), dStat.Version); err != nil {
		fmt.Println("update error : ", err)
	}

	// 删除
	_, dStat, _ = conn.Get("/test_tree2")
	if err = conn.Delete("/test_tree2", dStat.Version); err != nil {
		fmt.Println("Delete error : ", err)
	}

	// 验证存在
	hasNode, _, err := conn.Exists("/test_tree2")
	if err != nil {
		fmt.Println("Exist error : ", err)
	}
	fmt.Println("Node exsit : ", hasNode)

	// 增加
	if _, err := conn.Create("/test_tree2", []byte("tree_connect"), 0, zk.WorldACL(zk.PermAll)); err != nil {
		// panic(err)
		fmt.Println("create error : ", err)
	}

	// 设置子节点
	if _, err = conn.Create("/test_tree2/subnode", []byte("node_content_children"), 0, zk.WorldACL(zk.PermAll)); err != nil {
		fmt.Println("create children node error : ", err)
	}

	// 获取子节点列表
	childNodes, _, err := conn.Children("/test_tree2")
	if err != nil {
		fmt.Println("Children error : ", err)
	}
	fmt.Println("childNode : ", childNodes)
}
