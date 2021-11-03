/*
 * @Author: your name
 * @Date: 2021-10-03 02:59:09
 * @LastEditTime: 2021-10-09 21:13:28
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \src\test_1\test.go
 */
package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// test 判断是否回调
var test = true
var ticket sync.WaitGroup

// ticket
var total_ticket chan int //定义一个channel来作为火车票的总票数
//采用通道实现产品双向传递，总票数100
//通道的运作方式替代了伪代码中 empty 和 full 信号量的作用
//通道的另一个特性：任意时刻只能有一个协程能对 channel 中某一个item进行访问。替代了伪代码中 mutex 信号量的作用,这里就不需要锁的应用了。
// 四个取票口均共享 这100张票。

// goroutine 对应 CSP 中并发执行的实体，channel 也就对应着 CSP 中的 channel。
// 也就是说，CSP 描述这样一种并发模型：多个Process 使用一个 Channel 进行通信,
// 这个 Channel 连结的 Process 通常是匿名的，消息传递通常是同步的（有别于 Actor Model）。

// 当票数小于20的时候的时候，可以由取票口选择是否回调 到自己选定的位置。
func call_to_back() int {

	var n int
	fmt.Printf("Please input your tickets:")
	//fmt.Scanf("%d", &n)
	n = rand.Intn(50)
	fmt.Printf("总票数触发后将恢复到%v张\n", n)
	return n
}

func sell(i int, n *sync.WaitGroup) { //the params i represents the window number
	for {
		count := <-total_ticket
		if count > 0 {

			if count == 20 && test {
				var b int
				b = call_to_back()
				total_ticket <- b
				test = false
			} else {
				time.Sleep(time.Duration(rand.Intn(10)))
				// 售出每张票后， channel - 1
				total_ticket <- (count - 1)
				fmt.Printf("id:%v号的售窗口,售出了第%v张ticket\n", i, count)
			}
		} else {
			break
		}
	}
	// 每个 窗口开始运行
	n.Done()
}

func main() {
	runtime.GOMAXPROCS(2)
	total_ticket = make(chan int, 4)
	// 定义缓冲区为 100，即 火车票总数100
	total_ticket <- 100
	rand.Seed(time.Now().Unix())
	ticket.Add(3)
	//创造四个 售窗口
	for i := 0; i < 4; i++ {
		go sell(i, &ticket)
	}
	//主要功能停止，关闭channel
	time.Sleep(4e9)
	fmt.Println("所有票均已经售出")
	close(total_ticket)
	ticket.Wait()
}
