package utils

import "fmt"

func Try(userFunc func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("程序执行发生异常，错误抛出!", err)
		}
	}()

	userFunc()
}
