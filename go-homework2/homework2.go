package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// ✅指针
/*题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
考察点 ：指针的使用、值传递与引用传递的区别。*/
func sum_ten(x *int) {
	fmt.Printf("修改前：%d\n", *x)
	*x = *x + 10
	fmt.Printf("修改后：%d\n", *x)
}

// 题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
// 考察点 ：指针运算、切片操作。
func int_slice(s *[]int) {
	fmt.Printf("修改前：%d\n", *s)
	for i := 0; i < len(*s); i++ {
		(*s)[i] = (*s)[i] * 2
	}
	fmt.Printf("修改后：%d\n", *s)
}

// ✅Goroutine
// 题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
// 考察点 ： go 关键字的使用、协程的并发执行。
func gr_Odd(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 10; i += 2 {
		fmt.Printf("奇数：%d\n", i)
	}
}

func gr_Even(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 2; i <= 10; i += 2 {
		fmt.Printf("偶数：%d\n", i)
	}
}

// 题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
// 考察点 ：协程原理、并发任务调度。
type Task struct {
	Name string
	Work func()
}

type TaskScheduler struct {
	tasks []Task
	wg    sync.WaitGroup
}

func (ts *TaskScheduler) AddTask(name string, work func()) {
	ts.tasks = append(ts.tasks, Task{Name: name, Work: work})
}

func (ts *TaskScheduler) RunTask(task Task) {
	defer ts.wg.Done()

	start := time.Now()
	fmt.Printf("[%s] 开始\n", task.Name)

	task.Work()

	duration := time.Since(start)
	fmt.Printf("[%s]完成,耗时:%v\n", task.Name, duration)
}

func (ts TaskScheduler) runAll() {
	for _, task := range ts.tasks {
		ts.wg.Add(1)
		go ts.RunTask(task)
	}
	ts.wg.Wait()
}

// ✅面向对象
// 题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle
// 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的
// Area() 和 Perimeter() 方法。
// 考察点 ：接口的定义与实现、面向对象编程风格。
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	width, height float64
}
type Circle struct {
	Radius float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// 题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee
// 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个
// PrintInfo() 方法，输出员工的信息。
// 考察点 ：组合的使用、方法接收者。
type Person struct {
	Name string
	Age  int
}
type Employee struct {
	Person
	EmployeeID int
}

func (e Employee) PrintInfo() {
	fmt.Printf("姓名：%s\n", e.Name)
	fmt.Printf("年龄：%d\n", e.Age)
	fmt.Printf("工号：%d\n", e.EmployeeID)
}

// ✅Channel
// 题目 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，
// 并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
// 考察点 ：通道的基本使用、协程间通信。
func producer(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 10; i++ {
		ch <- i
		fmt.Printf("发送：%d\n", i)
	}
	close(ch)
}
func consumer(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		fmt.Printf("接收到：%d\n", num)
	}
}

// 题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道
// 中接收这些整数并打印。
// 考察点 ：通道的缓冲机制。
func producer_1(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 100; i++ {
		ch <- i
		fmt.Printf("缓冲区：%d|%d\n", len(ch), cap(ch))
	}
	close(ch)
}
func consumer_1(id int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	count := 0
	for num := range ch {
		count++
		fmt.Printf("消费者%d 接收了：%d(第%d个)\n", id, num, count)
	}
	fmt.Printf("消费者%d 完成，共接收 %d 个数据\n", id, count)
}

// ✅锁机制
// 题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协
// 程对计数器进行1000次递增操作，最后输出计数器的值。
// 考察点 ： sync.Mutex 的使用、并发数据安全。

func increment(id int, counter *int, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	localCount := 0
	for i := 0; i < 1000; i++ {
		mu.Lock()
		*counter++
		mu.Unlock()
		localCount++
	}
	fmt.Printf("协程%d 完成，递增了 %d 次\n", id, localCount)
}

// 题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程
// 对计数器进行1000次递增操作，最后输出计数器的值。
// 考察点 ：原子操作、并发数据安全。
func increment_1(id int, counter *int64, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 1000; i++ {
		atomic.AddInt64(*&counter, 1)
	}
	fmt.Printf("协程%d,已完成\n", id)
}

func main() {
	// x := 5
	// sum_ten(&x)
	// fmt.Println("---------------------------------------------")

	// s := []int{1, 3, 5}
	// int_slice(&s)
	// fmt.Println("---------------------------------------------")

	// var wg sync.WaitGroup
	// wg.Add(2)
	// go gr_Odd(&wg)
	// go gr_Even(&wg)
	// wg.Wait()
	// fmt.Println("所有协程执行完毕")
	// fmt.Println("---------------------------------------------")

	// scheduler := &TaskScheduler{}
	// scheduler.AddTask("任务1", func() {
	// 	time.Sleep(1 * time.Second)
	// })
	// scheduler.AddTask("任务2", func() {
	// 	time.Sleep(3 * time.Second)
	// })
	// scheduler.runAll()
	// fmt.Println("---------------------------------------------")

	// rect := Rectangle{width: 5, height: 6}
	// fmt.Printf("矩形面积：%.2f\n", rect.Area())
	// fmt.Printf("矩形周长：%.2f\n", rect.Perimeter())
	// Circ := Circle{Radius: 5}
	// fmt.Printf("圆形面积：%.2f\n", Circ.Area())
	// fmt.Printf("圆形周长：%.2f\n", Circ.Perimeter())
	// fmt.Println("---------------------------------------------")

	// emp := Employee{
	// 	Person: Person{
	// 		Name: "张三",
	// 		Age:  18,
	// 	},
	// 	EmployeeID: 100001,
	// }
	// emp.PrintInfo()
	// fmt.Println("---------------------------------------------")

	// ch := make(chan int)
	// var wg sync.WaitGroup
	// wg.Add(2)
	// go producer(ch, &wg)
	// go consumer(ch, &wg)
	// wg.Wait()
	// fmt.Println("---------------------------------------------")

	// ch := make(chan int, 10)
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go producer_1(ch, &wg)
	// numConsumer := 3
	// wg.Add(numConsumer)
	// for i := 1; i <= numConsumer; i++ {
	// 	go consumer_1(i, ch, &wg)
	// }
	// wg.Wait()
	// fmt.Println("---------------------------------------------")

	// var wg sync.WaitGroup
	// var mu sync.Mutex
	// var counter = 0
	// for i := 1; i <= 10; i++ {
	// 	wg.Add(1)
	// 	go increment(i, &counter, &wg, &mu)
	// }
	// wg.Wait()
	// fmt.Printf("最终结果: %d\n", counter)
	// fmt.Println("---------------------------------------------")

	var counter int64
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go increment_1(i, &counter, &wg)
	}
	wg.Wait()
	finalValue := atomic.LoadInt64(&counter)
	fmt.Printf("最终结果：%d\n", finalValue)
	fmt.Println("---------------------------------------------")
}
