package main

import "fmt"

func quicksort(a []int) []int {

	if len(a) < 2 {
		return a
	}

	left, right := 0, len(a)-1
	pivot := right
	right--
	for left <= right {
		if a[left] <= a[pivot] {
			left++
			continue
		}

		if a[right] > a[pivot] {
			right--
			continue
		}
		a[left], a[right] = a[right], a[left]
	}

	a[left], a[pivot] = a[pivot], a[left]

	quicksort(a[:left])
	quicksort(a[left+1:])

	return a
}

func main() {
	nums := []int{6, 3}
	fmt.Println(quicksort(nums))
}

//
//// 项目运行初始配置
//func init() {
//	initialize.Init()
//}
//
//func main() {
//	//开启socket服务
//	//go service.SocketServer()
//	//engine := router.GetEngine()
//	//if err := engine.Run(":7001"); err != nil {
//	//	panic(err)
//	//}
//}
