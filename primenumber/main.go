package main

import (
	"fmt"
	"strconv"
	"net/http"
)

//生成器，2，3，4，5，6，7，8....
func generator() chan int {
	ch := make(chan int)
	go func(){
		for i := 2; ; i++ {
			ch <- i
		}
	}()
	return ch
}

//过滤器，将通道能被num整除的数据过滤掉
func filter(ch chan int, num int) chan int{
	out := make(chan int)
	go func(){
		for {
			cur := <- ch
			if cur % num != 0 {
				out <- cur
			}
		}
	}()
	return out
}

//通道里是阈值内的所有素数
func primeNum(threshold int) chan int{
	prime := make(chan int)
	gen := generator()
	num := <- gen
	go func(){
		for num <= threshold {
			prime <- num
			gen = filter(gen, num)
			num = <- gen
		}
		close(prime)
	}()
	return prime
}

func ServeHTTP(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	input := r.FormValue("threshold")
	if len(input) == 0 {
		fmt.Fprintf(w,"请再地址栏输入阈值：?threshold=")
		return
	}
	threshold, err := strconv.Atoi(input)
	if err != nil {
		fmt.Fprintf(w,"非法输入！请重试")
		return
	}

	prime := primeNum(threshold)
	result := make([]int, 0)
	for v := range prime {
		result = append(result, v)
	}
	fmt.Fprintf(w,"output: %v",result)
}

func main(){
	http.HandleFunc("/primenumber",ServeHTTP)
	http.ListenAndServe(":9090",nil)
}