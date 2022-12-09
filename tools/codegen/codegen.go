package main

import "bytes"

func main() {

}

// 生成器对象，用来控制生成文档信息
type Generator struct {
	buf bytes.Buffer // 生成文档缓存对象
	pkg *Package     // 需要扫描的包
}

type Package struct {
}
