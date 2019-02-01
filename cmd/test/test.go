package main

import (
	"fmt"

	"github.com/miekg/dns"
)

func main() {
	c := new(dns.Client)
	// c.DialTimeout = time.Millisecond * 5000
	// c.ReadTimeout = time.Millisecond * 5000
	// c.WriteTimeout = time.Millisecond * 5000
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("baidu.com"), dns.TypeAAAA)
	r, _, err := c.Exchange(m, "[2001:4860:4860::8844]:53")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(r.Answer)
	// r, _, err := c.Exchange(m, "114.114.114.114:53")
	// if err != nil {
	// 	fmt.Print("a" + err.Error())
	// }
	// fmt.Println(r.Answer)
	// var l sync.Mutex
	// c := sync.NewCond(&l)
	// f := func() {
	// 	c.L.Lock()
	// 	defer c.L.Unlock()
	// 	c.Wait()
	// 	fmt.Println("aaa")
	// }
	// go f()
	// go f()
	// time.Sleep(3 * time.Second)
	// c.L.Lock()
	// c.Broadcast()
	// c.L.Unlock()
	// time.Sleep(3 * time.Second)
}
