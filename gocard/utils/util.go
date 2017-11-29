//工具包，封装一些公用的函数
package utils

import (
	"os"
	"log"
)

//如果错误不为空的话，先log fatal，再退出程序
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}