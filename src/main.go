package main

import (
	"fmt"
	"sync"
	"time"
)

//Sem is the semephore implementaion
var (
	log = fmt.Println
)

//Stack data structure not thread safe
type Stack struct {
	Top int32
	Stk []interface{}
	Max int32
}

func (stack *Stack) push(data interface{}) {

	if stack.Top == stack.Max-1 {

		return
	}

	stack.Top += 1

	stack.Stk[stack.Top] = data

}

func (stack *Stack) pop() (data interface{}, ok bool) {

	if stack.Top == -1 {
		data = nil
		ok = false

		return
	}

	data = stack.Stk[stack.Top]

	stack.Top -= 1

	ok = true

	return
}

type Pool = Stack

type Sem struct {
	inUse        bool
	lock         *sync.Cond
	counter      int32
	NumResources int32
	Pool         *Pool
}

func (sem *Sem) Init() {

	sem.inUse = true
	sem.counter = sem.NumResources
	sem.lock = sync.NewCond(&sync.Mutex{})

}

func (sem *Sem) GetResource() (resource interface{}) {

	if sem.inUse {

		sem.lock.L.Lock()

		for sem.counter == 0 {
			log("waiting")
			sem.lock.Wait()

		}

		resource, _ = sem.Pool.pop()

		sem.counter--

		defer sem.lock.L.Unlock()

		return
	}

	resource = "Initialize Sem instance (call Init()) before use."

	return
}

func (sem *Sem) ReleaseResource(resource interface{}) (ok bool) {

	sem.lock.L.Lock()

	if sem.counter == sem.NumResources {
		ok = false
		return
	}

	sem.Pool.push(resource)

	sem.counter += 1

	sem.lock.L.Unlock()

	sem.lock.Broadcast()

	ok = true

	return
}

type Res struct {
	data int
}

func Work(sem *Sem, wg *sync.WaitGroup) {

	resource, _ := sem.GetResource().(Res)

	for i := 100; i > 0; i-- {
		resource.data++
		time.Sleep(5 * time.Millisecond)
	}

	sem.ReleaseResource(resource)
	wg.Done()
	return
}

//NewSemopher  provide instance of new semeophore
func NewSemophore(res []interface{}) *Sem {

	//-------------------------------------------
	NUM_OF_RESOURCES := int32(len(res))

	var pool = &Pool{
		Stk: res[:],
		Top: NUM_OF_RESOURCES - 1,
		Max: NUM_OF_RESOURCES,
	}

	resourcePool := &Sem{
		NumResources: NUM_OF_RESOURCES,
		Pool:         pool,
	}

	resourcePool.Init()

	return resourcePool
}

func main() {

	var res = []interface{}{
		Res{data: 0},
		Res{data: 0},
		Res{data: 0},
		Res{data: 0},
	}

	wg := &sync.WaitGroup{}

	//allocate semeophore
	resourcePool := NewSemophore(res)

	num := 10

	wg.Add(num)

	for i := 0; i < num; i++ {
		go Work(resourcePool, wg)
	}

	wg.Wait()

	for _, val := range res {
		fmt.Println(val)
	}
}
