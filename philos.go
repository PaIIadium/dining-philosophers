package main

import (
	"fmt"
	"sync"
)

type chopStick struct {
	sync.Mutex
}

type host struct {
	eatingPhiloNumbers []int
	mut                sync.Mutex
}

func (h *host) isEating(philoNum int) bool {
	for _, num := range h.eatingPhiloNumbers {
		if num == philoNum {
			return true
		}
	}
	return false
}

func (h *host) wantToEat(ph *philo) bool {
	if len(h.eatingPhiloNumbers) == 2 {
		return false
	}
	var numNeighbor1 int
	var numNeighbor2 int
	if ph.number == 1 {
		numNeighbor1 = 2
		numNeighbor2 = 5
	} else if ph.number == 5 {
		numNeighbor1 = 1
		numNeighbor2 = 4
	} else {
		numNeighbor1 = ph.number + 1
		numNeighbor2 = ph.number - 1
	}
	h.mut.Lock()
	defer h.mut.Unlock()
	neighborsNotEat := !h.isEating(numNeighbor1) && !h.isEating(numNeighbor2)
	if neighborsNotEat {
		h.eatingPhiloNumbers = append(h.eatingPhiloNumbers, ph.number)
		return true
	}
	return false
}

func (h *host) endEat(philoNum int) {
	h.mut.Lock()
	for i, num := range h.eatingPhiloNumbers {
		if num == philoNum {
			arr := h.eatingPhiloNumbers
			h.eatingPhiloNumbers = append(arr[:i], arr[i+1:]...)
		}
	}
	h.mut.Unlock()
}

type philo struct {
	host      *host
	number    int
	leftCP    *chopStick
	rightCP   *chopStick
	countLeft int
	isEating  bool
}

func (ph *philo) tryToEat() {
	needToEat := !ph.isEating && ph.countLeft > 0
	if !needToEat {
		return
	}
	canEat := ph.host.wantToEat(ph)
	if canEat {
		ph.leftCP.Lock()
		ph.rightCP.Lock()
		fmt.Println("starting to eat", ph.number)
		ph.isEating = true
		ph.countLeft--
		// time.Sleep(100)
		// you can uncomment the line above to make sure
		// other philos can start or finish eating at this point
		ph.host.endEat(ph.number)
		ph.isEating = false
		fmt.Println("finishing eating", ph.number)
		ph.leftCP.Unlock()
		ph.rightCP.Unlock()
	}
}

func (ph *philo) eat(wg *sync.WaitGroup) {
	for ph.countLeft != 0 {
		ph.tryToEat()
	}
	wg.Done()
}

func main() {
	h := &host{make([]int, 0), sync.Mutex{}}

	chopSticks := make([]*chopStick, 5)
	for i := 0; i < 5; i++ {
		chopSticks[i] = &chopStick{}
	}

	philos := make([]*philo, 5)
	for i := 1; i < 6; i++ {
		ph := &philo{
			host:      h,
			number:    i,
			leftCP:    chopSticks[(i-1)%5],
			rightCP:   chopSticks[(i)%5],
			countLeft: 3,
			isEating:  false,
		}
		philos[i-1] = ph
	}

	wg := &sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go philos[i].eat(wg)
	}
	wg.Wait()
}
