package main

import "github.com/BoynChan/GopherProxy/pkg"

// Author:Boyn
// Date:2020/9/8

func main() {
	initPkg()

}

func initPkg() {
	pkg.InitDB()
	pkg.InitRedisCli()
}
