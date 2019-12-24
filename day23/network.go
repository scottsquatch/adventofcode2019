package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/scottsquatch/adventofcode2019/utils/fileutils"

	"github.com/scottsquatch/adventofcode2019/common"
)

type sendaction func(message)

type message struct {
	dst, x, y int64
}
type irouter interface {
	poweron()
	receive(msg message)
	getport() chan message
	queuelength() int
}

type router struct {
	nic         *common.Computer
	nicsoftware *common.IntCodeProgram
	in, out     chan int64
	msgqueue    []int64
	msglock     *sync.Mutex
	address     int
	port        chan message
}

func newRouter(software *common.IntCodeProgram, address int) *router {
	in, out := make(chan int64), make(chan int64)
	nic := common.NewComputer(in, out)
	r := &router{nic: nic, nicsoftware: software, msgqueue: []int64{}, msglock: &sync.Mutex{}, address: address, in: in, out: out, port: make(chan message)}

	return r
}

func (r *router) poweron() {
	go r.nic.Run(r.nicsoftware)

	r.in <- int64(r.address)
	r.in <- -1

	go func() {
		for {
			msg := message{dst: <-r.out, x: <-r.out, y: <-r.out}
			// fmt.Printf("%d: send msg %v\n", r.address, msg)
			r.port <- msg
		}
	}()

	go func() {
		for {
			r.in <- r.consume()
		}
	}()
}

func (r *router) consume() int64 {
	msg := int64(-1)
	r.msglock.Lock()
	if len(r.msgqueue) > 0 {
		msg, r.msgqueue = r.msgqueue[0], r.msgqueue[1:]
		// fmt.Printf("%d: consume %d msgqueue: %v\n", r.address, msg, r.msgqueue)
	}
	r.msglock.Unlock()

	return msg
}

func (r *router) receive(msg message) {
	r.msglock.Lock()
	r.msgqueue = append(r.msgqueue, msg.x, msg.y)
	r.msglock.Unlock()
}

func (r *router) getport() chan message {
	return r.port
}

func (r *router) queuelength() int {
	r.msglock.Lock()
	l := len(r.msgqueue)
	r.msglock.Unlock()

	return l
}

type nat struct {
	lastmsg   message
	msglock   *sync.Mutex
	port      chan message
	routermap map[int]irouter
}

func newNat(routermap map[int]irouter) *nat {
	return &nat{msglock: &sync.Mutex{}, port: make(chan message), routermap: routermap}
}

func (n *nat) poweron() {
	// not great to have polling in asynchronous system, but works for now
	go func() {
		for {
			time.Sleep(5 * time.Second)
			// fmt.Println("Start network idle check")
			networkidle := true
			for addr, r := range n.routermap {
				if addr != 255 && r.queuelength() > 0 {
					networkidle = false
					break
				}
			}

			if networkidle {
				// fmt.Printf("Network idle, sending %v to %d\n", n.lastmsg, 0)
				n.port <- n.lastmsg
				n.routermap[0].receive(n.lastmsg)
			}
			// fmt.Println("End idle check")
		}
	}()
}

func (n *nat) queuelength() int {
	return 0
}

func (n *nat) receive(msg message) {
	n.msglock.Lock()
	n.lastmsg = msg
	// fmt.Printf("NAT received %v\n", msg)
	n.msglock.Unlock()
}

func (n *nat) getport() chan message {
	return n.port
}

func runPartA() {
	software := common.NewIntCodeProgram(fileutils.ReadFile(os.Args[1]))
	done := make(chan bool)
	routermap := make([]irouter, 50)
	sendlistner := func(r irouter) {
		for {
			m := <-r.getport()
			if m.dst == 255 {
				fmt.Printf("Message sent to 255 %v\n", m)
				done <- true
				return
			}
			routermap[m.dst].receive(m)
		}
	}
	for i := 0; i < 50; i++ {
		r := newRouter(software, i)
		routermap[i] = r
		go sendlistner(r)
	}
	for _, r := range routermap {
		go r.poweron()
	}

	<-done
}

func runPartB() {
	software := common.NewIntCodeProgram(fileutils.ReadFile(os.Args[1]))
	done := make(chan bool)
	routermap := make(map[int]irouter)
	natmsgyhistory := make(map[int64]bool)
	sendlistner := func(r irouter) {
		for {
			m := <-r.getport()
			routermap[int(m.dst)].receive(m)
		}
	}

	for i := 0; i < 50; i++ {
		r := newRouter(software, i)
		routermap[i] = r
		go sendlistner(r)
	}

	routermap[255] = newNat(routermap)
	go func(r irouter) {
		for {
			m := <-r.getport()
			dup := natmsgyhistory[m.y]
			// fmt.Printf("Message %v, messageys %v\n", m, natmsgyhistory)
			if dup {
				fmt.Printf("First duplicate Y value sent to 0 from Nat: %d\n", m.y)
				done <- true
			} else {
				natmsgyhistory[m.y] = true
			}
		}
	}(routermap[255]) // Special send listner from nat to intercept part b completion condition
	for _, r := range routermap {
		go r.poweron()
	}

	<-done
}

func main() {
	runPartA()
	runPartB()
}
