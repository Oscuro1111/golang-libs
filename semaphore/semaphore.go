package semaphore

import (
	"sync"
)



//Stack data structure not thread safe
type semStack struct {
	Top int32
	Stk []interface{}
	Max int32
}

func (stack *semStack) push(data interface{}) {

	if stack.Top == stack.Max-1 {

		return
	}

	//Making space
	stack.Top += 1

	stack.Stk[stack.Top] = data

}

func (stack *semStack) pop() (data interface{}, ok bool) {

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

type Pool = semStack

type sSem struct {
	inUse        bool
	lock         *sync.Cond
	counter      int32
	NumResources int32
	Pool         *Pool
}

func (sem *sSem) Init() {

	sem.inUse = true
	sem.counter = sem.NumResources
	sem.lock = sync.NewCond(&sync.Mutex{})

}

func (sem *sSem) GetResource() (resource interface{}) {

	if sem.inUse {

		sem.lock.L.Lock()
		//waiting
		for sem.counter == 0 {
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

func (sem *sSem) ReleaseResource(resource interface{}) (ok bool) {

	sem.lock.L.Lock()

	if sem.counter == sem.NumResources {
		ok = false
		return
	}

	//Put back the resource in pool
	sem.Pool.push(resource)

	sem.counter += 1

	sem.lock.L.Unlock()

	//Make message for sleeping threads for availity of resource in pool
	sem.lock.Broadcast()

	ok = true

	return
}

type Semaphore = sSem

//NewSemapher  provide instance of new semaphore
func NewSemaphore(res []interface{}) *Semaphore {

	//-------------------------------------------
	NUM_OF_RESOURCES := int32(len(res))

	var pool = &Pool{
		Stk: res[:],
		Top: NUM_OF_RESOURCES - 1,
		Max: NUM_OF_RESOURCES,
	}

	resourcePool := &Semaphore{
		NumResources: NUM_OF_RESOURCES,
		Pool:         pool,
	}

	resourcePool.Init()

	return resourcePool
}
