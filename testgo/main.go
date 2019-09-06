package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
)

const (
	MAXIII = 98
)

func main() {
	//defer_call()

	// a := 1
	// b := 2
	// defer calc("1", a, calc("10", a, b)) // defer 关键字修饰的函数逆序执行，貌似再往里的逻辑顺序执行的 所以是10-20-2-1
	// a = 0
	// defer calc("2", a, calc("20", a, b))
	// b = 1
	// ccccc()
	//ddddd()

	// runtest()

	// child := Child{}
	// child.MethodA() //ma mb

	// u := &UserAges{}
	// u.Add("zhang", 22)
	// u.Get("zhang")

	//jsontest()

	// test7()

	//test8()
	// test9()
	// u := User{1, "OK", 12}
	// testReflect(u)

	// str1 := "abcd"
	// str2 := "d"
	// str3 := test(str1, str2)
	// fmt.Println(str3)

	//pase_student()
	// cloudTest1()
	//mysql()
	//TestArrayAndSlice()
	temp, _ := strconv.Atoi("123")
	fmt.Println(temp)
}

// 畅聊面试一
func defer_call() {
	defer func() {
		fmt.Println("before")
	}()
	defer func() {
		recover()
	}()
	defer func() {
		fmt.Println("after")
	}()
	panic("panic info")
}

// 畅聊面试二
func calc(index string, a, b int) int {
	ret := a + b
	fmt.Println(index, a, b, ret)
	return ret
}

// defer的测试
func ccccc() {
	i := 0
	defer fmt.Println(i) //传递i=0进去
	i++
	return
}

func ddddd() {
	for i := 0; i < 3; i++ {
		//defer fmt.Println(i) //2 1 0
		defer func(i int) {
			fmt.Println(i) // 2 1 0 作为参数，拷贝
		}(i)
		// defer func() {
		// 	fmt.Println(i) //3 3 3 作为引用
		// }()
	}
}

// 畅聊面试三
func runtest() {
	runtime.GOMAXPROCS(1)
	// for {
	int_chan := make(chan int, 1)
	string_chan := make(chan string, 1)
	int_chan <- 1
	string_chan <- "hello"
	select {
	case value := <-int_chan:
		//fmt.Println(value + " out put") // 面试题：编译就不会通过的 int类型 + string 类型
		fmt.Printf("%d %s", value, " out put")
	case value := <-string_chan:
		panic(value + " out put")
		//fmt.Println(value)
	}
	// }
}

// 畅聊面试四
type Parent struct{}

func (p *Parent) MethodB() {
	fmt.Println("mb from p")
}

func (p *Parent) MethodA() {
	fmt.Println("ma from p")
	p.MethodB()
}

type Child struct {
	Parent
}

func (c *Child) MethodB() {
	fmt.Println("cb from c")
}

// 畅聊面试五
type UserAges struct {
	ages map[string]int
	sync.Mutex
}

func (u *UserAges) Add(name string, age int) {
	u.Lock()
	defer u.Unlock()
	u.ages[name] = age
}

func (u *UserAges) Get(name string) int {
	if age, ok := u.ages[name]; ok {
		return age
	}
	return -1
}

// 畅聊面试六
func jsontest() {
	json_str := []byte(`{"age":1}`)
	var value map[string]interface{}
	json.Unmarshal(json_str, &value)
	age := value["age"]
	fmt.Println(reflect.TypeOf(age))
}

// 畅聊面试七
func test7() {
	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println("i ", i)
			wg.Done()
		}()
	}
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("j ", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

// 畅聊面试八
func test8888(a, b int64) int64 {
	fmt.Println(a)
	fmt.Println(b)
	max := math.Max(float64(a), float64(b))
	fmt.Println(max)
	return int64(max)
}

func test8() {
	//fmt.Println(test8888(1, 2))
	fmt.Println(math.MaxInt64)
	fmt.Println(test8888(math.MaxInt64-2, math.MaxInt64-1), math.MaxInt64-1)

	//float64 转为 int64 会产生开销，且转换值太大的数会发生截断
}

func test9() {
	child := Child{}
	fmt.Println(reflect.TypeOf(child))
	fmt.Println(reflect.ValueOf(child))
	selfMap := make(map[string]interface{})
	selfMap["name"] = 3
	m := selfMap["name"]
	fmt.Println(reflect.TypeOf(m))
}

type User struct {
	Id   int
	Name string
	Age  int
}

func (u User) Hello() {
	fmt.Println("hello world")
}

func testReflect(o interface{}) {
	t := reflect.TypeOf(o)
	fmt.Println(t)                 //main.User
	fmt.Println("Type:", t.Name()) //Type:User

	if k := t.Kind(); k != reflect.Struct {
		fmt.Println("not struct")
		return
	}

	v := reflect.ValueOf(o)
	fmt.Println(v)
	fmt.Println("Fields:")

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Println(f)
		val := v.Field(i).Interface()
		fmt.Printf("%6s: %v = %v\n", f.Name, f.Type, val)
	}

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		fmt.Println(m)
		fmt.Printf("%6s: %v\n", m.Name, m.Type)
	}
}

func test(str1, str2 string) string {
	str1B := []byte(str1)
	str2B := []byte(str2)
	temp := make([]byte, 10)
	for _, v1 := range str1B {
		isExist := false
		for _, v2 := range str2B {
			if v1 == v2 {
				isExist = true //存在
			}
		}
		if !isExist {
			temp = append(temp, v1)
		}
	}
	return string(temp)
}

// cSphere - 希云科技 做docker那家
type kvStore map[string]string

func (s kvStore) set(k, v string) {
	s[k] = v
}

func (s kvStore) get(k string) string {
	return s[k]
}

func cSphere() {
	listener, _ := net.Listen("tcp", ":2333")
	store := make(kvStore)
	for {
		conn, _ := listener.Accept()
		buf := make([]byte, 1024)
		conn.Read(buf)
		cmd := strings.Split(string(buf), " ")
		switch cmd[0] {
		case "set":
			store.set(cmd[1], strings.Join(cmd[2:], " "))
		case "get":
			conn.Write([]byte(store.get(cmd[1])))
		}
		conn.Close()
	}
}

// 时速云科技
// 面试题一
type student struct {
	Name string
	Age  int
}

func pase_student() {
	m := make(map[string]*student)
	stus := []student{
		{Name: "zhou", Age: 24},
		{Name: "li", Age: 23},
		{Name: "wang", Age: 22},
	}

	for index, stu := range stus {
		//m[stu.Name] = &stu //地址都是值拷贝（同一个临时变量）的地址
		m[stu.Name] = &stus[index]
	}
	fmt.Println(m)
}

// 面试题二
func cloudTest1() {
	s := make([]int, 5)
	s = append(s, 1, 2, 3)
	fmt.Println(s) // 0 0 0 0 0 1 2 3
}

// 造大量假数据
//结构体
type Job struct {
	db    *sql.DB
	ch    chan int
	total int
	n     int
}

func mysql() {
	db, err := sql.Open("mysql", "root:zhanglijun@/start_test?charset=utf8")
	if err != nil {
		fmt.Println("访问数据库出错", err)
		return
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Second * 500) //设置连接超时500秒
	db.SetMaxOpenConns(100)                  //设置最大连接数

	total := 25000 //每次插入2.5万条数据
	gonum := 400   //启用400个协程
	fmt.Println("====start=====")
	start := time.Now()
	// 测试插入数据库的功能,每次最多同步20个工作协程
	jobChan := make(chan Job, 20) //任务队列
	go worker(jobChan)
	//统计使用次数
	ch := make(chan int, gonum)
	for n := 0; n < gonum; n++ {
		job := Job{
			db:    db,    //数据库连接
			ch:    ch,    //使用次数
			total: total, //每次插入数据量
			n:     n,     //第几次连接
		}
		jobChan <- job // 每塞满20个便阻塞
	}
	ii := 0
	for {
		<-ch
		ii++
		if ii >= gonum {
			break
		}
	}

	end := time.Now()
	curr := end.Sub(start)
	fmt.Println("run time:", curr)
}

func worker(jobChan <-chan Job) {
	for job := range jobChan {
		go sqlExec(job)
	}
}

func sqlExec(job Job) {
	buf := make([]byte, 0, job.total)
	buf = append(buf, " insert into student(user_name) values "...)
	for i := 0; i < job.total; i++ {
		myid, _ := uuid.NewV4()
		userName := myid.String()
		if i == job.total-1 {
			buf = append(buf, "('"+userName+"');"...)
		} else {
			buf = append(buf, "('"+userName+"'),"...)
		}
	}
	ss := string(buf)
	fmt.Println("第" + strconv.Itoa(job.n) + "次插入2.5万条数据！")
	_, err := job.db.Exec(ss)
	checkErr(err)
	fmt.Println("完成---" + strconv.Itoa(job.n) + "次插入2.5万条数据！")
	job.ch <- 1
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func TestArrayAndSlice() {
	s1 := []int{1, 2, 3}
	s2 := s1[1:]
	for i := range s2 {
		s2[i] += 10
	}
	fmt.Println(s2)
	s2 = append(s2, 4)
	for i := range s2 {
		s2[i] += 10
	}
	fmt.Println(s2)
}
