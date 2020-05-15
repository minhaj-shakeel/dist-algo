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
                  reqTimeStamp : 0,
                  requested : false,
                  pendingReplies : make([] Pair,0),
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
      go func(i int) {      reply := <-ch2[i]
                            p.mux.Lock()
                            p.clock = max(p.clock,reply.TimeStamp)+1
                            p.mux.Unlock()
                            fmt.Println("reply received on ",curr,"from ",i,reply)
                        wg.Done() } (i)
    }
  }
  wg.Wait()
  fmt.Println("recv reply",curr)
}

func (p *Process) SendPendingReplies(ch2 [] chan Pair){
  curr := p.pid
  fmt.Println(curr)
  num := len(p.pendingReplies)
  p.mux.Lock()
  p.clock++
  p.mux.Unlock()
  fmt.Println(p.pendingReplies)
  for i:=0;i<num;i++{
    //Change to concurrent send
    //req:=p.pendingReplies[0]
    ch2[curr] <-Pair{curr,p.clock}
    p.pendingReplies = p.pendingReplies[1:]
  }

}

func (p *Process) DoSomething(){
  fmt.Printf("Process %d Entered\n",p.pid)
  time.Sleep(10*time.Second)
  fmt.Printf("Process %d Exited\n",p.pid)
}


func (p *Process) EnterCS(ch1,ch2 [] chan Pair, wg *sync.WaitGroup){
  p.SendRequest(ch1)
  fmt.Println("request",p.pid)
  p.RecvReply(ch2)
  fmt.Println("received reply")
  p.requested=false
  p.inCS = true
  p.DoSomething()
  p.inCS = false
  p.SendPendingReplies(ch2)
  wg.Done()
  fmt.Println("process 0 completed")
}

func main(){
  N:=5 //Number of Processes

  //Initialising Channels for Request Receiving and Sending
  ch1 := make([] chan Pair,N) 
  for i,_:= range ch1{
    ch1[i] = make(chan Pair)
  }
  //Initialising Channels for token Receiving and Sending
  ch2 := make([] chan Pair,N)
  for i,_ := range ch2{
    ch2[i] = make(chan Pair)
  }
  
  //Initialising List of Processes
  PList := make([] *Process,N)
  for  i:= 0;i < N ; i++{
    PList[i] = NewProcess(i)
  }

  for i:=0;i<N;i++{
    go PList[i].RecvRequest(ch1,ch2)
  }
  
  var wg1 sync.WaitGroup
  wg1.Add(1)
  go PList[0].EnterCS(ch1,ch2,&wg1)
  time.Sleep(time.Second)
  wg1.Add(1)
  go PList[1].EnterCS(ch1,ch2,&wg1)
  time.Sleep(time.Second)
  //wg1.Add(1)
  //go PList[2].EnterCS(ch1,ch2,&wg1)
  wg1.Wait()
}

