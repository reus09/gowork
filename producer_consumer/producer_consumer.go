/*
 * @Author: your name
 * @Date: 2021-10-02 18:49:17
 * @LastEditTime: 2021-10-09 21:13:58
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \src\test\test.go
 */
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Product struct {
	prod_num int //生产该产品的生产者序号
	value    int
}

//采用通道实现产品双向传递，缓冲区设为10
//通道的运作方式替代了伪代码中 empty 和 full 信号量的作用
//通道的另一个特性：任意时刻只能有一个协程能对 channel 中某一个item进行访问。替代了伪代码中 mutex 信号量的作用,这里就不需要锁的应用了。

// goroutine 对应 CSP 中并发执行的实体，channel 也就对应着 CSP 中的 channel。
// 也就是说，CSP 描述这样一种并发模型：多个Process 使用一个 Channel 进行通信,
// 这个 Channel 连结的 Process 通常是匿名的，消息传递通常是同步的（有别于 Actor Model）.
// 因此 Channel 的使用 解决了 进程通信

var ch = make(chan Product, 10)

//终止生产者生产的信号，否则无限执行下去
var stop = false

func Producer(wgp *sync.WaitGroup, wgc *sync.WaitGroup, p_num int) {
	for !stop { //不断生成，直到主线程执行到终止信号
		p := Product{prod_num: p_num, value: rand.Int()}
		// 判断如果 缓冲区大于等于10 ，消费者工作，生产者不工作，则call to back
		if len(ch) == 10 {
			call_to_back()
		}
		ch <- p
		fmt.Printf("sum product:%v:producer %v produce a product: %#v\n", len(ch), p_num, p)
		time.Sleep(time.Duration(200+rand.Intn(1000)) * time.Millisecond) //延长执行时间
	}
	wgp.Done()
}

func Consumer(wgp *sync.WaitGroup, wgc *sync.WaitGroup, c_num int) {
	//通道里没有产品了就停止消费

	// 判断如果 缓冲区等于0 ，此时生产者工作，消费者停止，同样call to back
	if len(ch) == 0 {
		fmt.Println("begin call to back")
		call_to_back()
	}
	for p := range ch {
		fmt.Printf("sum product:%v:consumer %v consume a product: %#v\n", len(ch), c_num, p)
		time.Sleep(time.Duration(200+rand.Intn(1000)) * time.Millisecond)
	}
	wgc.Done()
}

// 当触发到缓冲区为 0 或者 10的情况时候，由系统自身恢复缓冲区的存取的产品任意个数。
func call_to_back() chan Product {
	id := rand.Intn(10)
	var dh = make(chan Product, 10)
	for i := 0; i <= id; i++ {
		p := Product{prod_num: i, value: rand.Int()}
		dh <- p
	}
	return dh
}

func main() { //主线程
	var wgp sync.WaitGroup
	var wgc sync.WaitGroup
	wgp.Add(5)
	wgc.Add(5)
	//设5个生产者、5个消费者
	for i := 0; i < 5; i++ {
		go Producer(&wgp, &wgc, i)
		go Consumer(&wgp, &wgc, i)
	}
	time.Sleep(time.Duration(1) * time.Second)

	stop = true
	wgp.Wait()
	// 关闭ch 即关闭缓冲区，不工作。
	close(ch)

	fmt.Printf("sum:%v,消费完毕！", len(ch))
	wgc.Wait()
}
