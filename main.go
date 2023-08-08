package main

import (
	"fmt"
	"sync"
	"time"
)

// The Dining Philosophers problem is well known in computer science circles.
// Five philosophers, numbered from 0 through 4, live in a house where the
// table is laid for them; each philosopher has their own place at the table.
// Their only difficulty – besides those of philosophy – is that the dish
// served is a very difficult kind of spaghetti which has to be eaten with
// two forks. There are two forks next to each plate, so that presents no
// difficulty. As a consequence, however, this means that no two neighbours
// may be eating simultaneously, since there are five philosophers and five forks.
//
// This is a simple implementation of Dijkstra's solution to the "Dining
// Philosophers" dilemma.

const (
  MEALS_NUMBER = 3
  EAT_TIME = 0 * time.Second
  THINK_TIME = 0 * time.Second
  SLEEP_TIME = 0 * time.Second
)

type Philosopher struct {
  name string
  leftFork int
  rightFork int
}

func main() {

  fmt.Println("Dining Philosopher Problem")
  fmt.Println("--------------------------")
  fmt.Println("The table is empty")


  philosophers := []Philosopher {
    { name: "Plato", leftFork: 4, rightFork: 0 },
    { name: "Socrates", leftFork: 0, rightFork: 1 },
    { name: "Aristotle", leftFork: 1, rightFork: 2 },
    { name: "Pascal", leftFork: 2, rightFork: 3 },
    { name: "Locke", leftFork: 3, rightFork: 4 },
  }

  ordersCh := make(chan string, len(philosophers))
  defer close(ordersCh)
  ordersCompleted := []string{}
  dine(philosophers, ordersCh)

  fmt.Println("The table is empty")

  for i := 0; i < len(philosophers); i++ {

    ordersCompleted = append(ordersCompleted, <-ordersCh)
  }
  fmt.Println("Orders completed: ", ordersCompleted)
  
}

func dine(philosophers []Philosopher, ordersCh chan string) {

  eatingWG := &sync.WaitGroup{}
  eatingWG.Add(len(philosophers))

  seatedWG := &sync.WaitGroup{}
  seatedWG.Add(len(philosophers))

  var forks = make(map[int]*sync.Mutex)

  for i := 0; i < len(philosophers); i++ {
    forks[i] = &sync.Mutex{}
  }

  for i := 0; i < len(philosophers); i++ {

    go eat(EatParameters{
      philosopher: philosophers[i],
      eatingWG: eatingWG,
      forks: forks,
      seatedWG: seatedWG,
      ordersCh: ordersCh,
    })
  }

  eatingWG.Wait()
}

type EatParameters struct {
  philosopher Philosopher
  eatingWG *sync.WaitGroup
  forks map[int]*sync.Mutex
  seatedWG *sync.WaitGroup
  ordersCh chan string
}

func eat(params EatParameters) {

  defer params.eatingWG.Done()

  fmt.Printf("%s is seated at the table \n", params.philosopher.name)
  params.seatedWG.Done()
  params.seatedWG.Wait()

  for i := MEALS_NUMBER; i > 0; i-- {

    if params.philosopher.leftFork < params.philosopher.rightFork {
      params.forks[params.philosopher.rightFork].Lock()
      fmt.Printf("\t%s takes the right fork \n", params.philosopher.name)
      params.forks[params.philosopher.leftFork].Lock()
      fmt.Printf("\t%s takes the left fork \n", params.philosopher.name)
    } else {
      params.forks[params.philosopher.leftFork].Lock()
      fmt.Printf("\t%s takes the left fork \n", params.philosopher.name)
      params.forks[params.philosopher.rightFork].Lock()
      fmt.Printf("\t%s takes the right fork \n", params.philosopher.name)
    }


    fmt.Printf("\t%s has both forks and is eating \n", params.philosopher.name)
    time.Sleep(EAT_TIME);

    fmt.Printf("\t%s is thinking \n", params.philosopher.name)
    time.Sleep(THINK_TIME);

    params.forks[params.philosopher.leftFork].Unlock()
    params.forks[params.philosopher.rightFork].Unlock()

    fmt.Printf("\t%s put down the forks \n", params.philosopher.name)
  }

  fmt.Printf("%s left the table \n", params.philosopher.name)
  params.ordersCh <- params.philosopher.name
}
