package main

import(
  "fmt"
  "sync"
  "time"
)

var wg sync.WaitGroup

func max(a,b int) int {
  var ret int
  switch a>b{
    case true:
      ret=a
    case false:
      ret=b
  }
  return ret
}

type Pair struct {
  pid int
  TimeStamp int
}

type Process struct{
  pid int
  mux sync.Mutex
  clock int
  inCS bool
  reqTimeStamp int
  requested bool
  pendingReplies [] Pair
}

func NewProcess(i int) *Process{
  return &Process{pid :i,
                  clock :0,
                  inCS : false,
                }
}

func (p *Process) SendRequest(ch1 [] chan Pair){
  curr := p.pid
  p.mux.Lock()
  p.clock++
  p.mux.Unlock()
  p.requested = true
  p.reqTimeStamp = p.clock
  for i:=0 ;i<len(ch1) ;i++{
    if i!=curr{
      go func(i int){ch1[i]<-Pair{curr,p.reqTimeStamp} }(i)
    }
  }
}

func (p *Process) RecvRequest(ch1 []chan Pair , ch2 [] chan Pair){
  curr :=p.pid
  for{
    req:=<-ch1[curr]
    p.mux.Lock()
    p.clock = max(p.clock,req.TimeStamp)+1
    p.mux.Unlock()
    if p.inCS == false{
      if ((p.requested == false) || ((p.requested == true) &&(p.reqTimeStamp > req.TimeStamp))){
        p.mux.Lock()
        p.clock++
        p.mux.Unlock()
        go func(){ ch2[curr] <- Pair{curr,p.clock} }()
      } else{
        p.pendingReplies = append(p.pendingReplies,req)
      }
    } else{
      p.pendingReplies = append(p.pendingReplies,req)
    }
  }


}

func (p *Process )RecvReply(ch2 [] chan Pair){
  curr :=p.pid
  for i:=0 ; i <len(ch2) ; i++{
    if i!=curr{
      wg.Add(1)
      go func(i int){   for {
                            reply := <-ch2[i]
                            p.mux.Lock()
                            p.clock = max(p.clock,reply.TimeStamp)+1
                            p.mux.Unlock()
                            if reply.TimeStamp>p.reqTimeStamp { break }
                          }
                        wg.Done() }(i)
    }
  }
  wg.Wait()
}

func (p *Process) DoSomething(){
  fmt.Printf("Process %d Entered\n",p.pid)
  time.Sleep(5*time.Second)
  fmt.Printf("Process %d Exited\n",p.pid)
}


func (p *Process) EnterCS(ch1,ch2 [] chan Pair){
  p.SendRequest(ch1)
  p.RecvReply(ch2)
  p.DoSomething()


}

func main(){



}
