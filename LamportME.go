package main 

import( "fmt"
        "sort"
        "time"
        "sync"
)

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

type ByTimeStamp []Pair

func (q ByTimeStamp) Len() int {
  return len(q)
}

func (q ByTimeStamp) Swap(i,j int) {
  q[i],q[j]=q[j],q[i]
}

func (q ByTimeStamp) Less(i,j int) bool{
  if (q[i].TimeStamp == q[j].TimeStamp){
    return q[i].pid < q[j].pid
  }
  return q[i].TimeStamp< q[j].TimeStamp
}

type Process struct{
  pid int
  clock int
  mux sync.Mutex
  reqTimeStamp int
  inCS bool
  ReqQueue []Pair
}

func NewProcess(i int) *Process{
  return &Process{pid :i,
                  clock:0,
                  inCS : false,
                  ReqQueue : make([]Pair,0),
                  }
}

var wg sync.WaitGroup

func (p *Process) SendRequest(ch1[]  chan Pair){
  curr :=p.pid
  p.mux.Lock()
  p.clock++
  p.mux.Unlock()
  p.reqTimeStamp = p.clock
  p.ReqQueue = append(p.ReqQueue,Pair{curr,p.clock})
  for i := 0 ; i < len(ch1) ; i++{
    if i!=curr{
      go func(i int){ ch1[i]<-Pair{curr,p.reqTimeStamp} }(i)
    }
  }
}

func (p *Process) RecvRequest(ch1 []chan Pair , ch2 [] chan Pair){
  curr :=p.pid
  for{
    req:= <-ch1[curr]
    p.mux.Lock()
    p.clock = max(p.clock,req.TimeStamp)+1
    p.mux.Unlock()
    p.ReqQueue=append(p.ReqQueue,req)
    sort.Sort(ByTimeStamp(p.ReqQueue))
    p.mux.Lock()
    p.clock++
    p.mux.Unlock()
    ch2[curr] <- Pair{curr,p.clock}
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

func (p *Process) SendRelease(ch3 []chan Pair){
  curr := p.pid
  p.mux.Lock()
  p.clock++
  p.mux.Unlock()
  for i,_ :=range ch3{
    if i!=curr{
      go func(i int){ch3[i] <- Pair{curr,p.clock} }(i)
    }
  }
  p.ReqQueue = p.ReqQueue[1:]
}

func (p *Process) RecvRelease(ch3 []chan Pair , ch []chan bool){
  curr :=p.pid
  for{
    release:= <-ch3[curr]
    p.ReqQueue = p.ReqQueue[1:]
    p.mux.Lock()
    p.clock = max(p.clock,release.TimeStamp)+1
    p.mux.Unlock()
    for{
      if(len(p.ReqQueue)!=0) {break}
    }

    if p.ReqQueue[0].pid == curr{
      ch[curr] <- true
    }
  }
}

func (p *Process) AcquireLock(ch1,ch2,ch3 []chan Pair,ch []chan bool,wg *sync.WaitGroup){
  p.SendRequest(ch1)
  p.RecvReply(ch2)
  for{
    status := <-ch[p.pid]
    if status == true{
      break
    }
  }
  p.inCS = true
  p.DoSomething()
  p.inCS = false
  p.SendRelease(ch3)
  wg.Done()
}

func main(){
  N:=5
  
  ch1 := make([]chan Pair,N)
  //Initialising Channels for Request Sending and Receiving
  for i,_ := range ch1{
    ch1[i] = make(chan Pair)
  }


  ch2 := make([]chan Pair,N)
  //Initialising Channels for Reply Sending and Receiving
  for i,_ := range ch2{
    ch2[i] = make(chan Pair)
  }
  
  ch3 := make([]chan Pair,N)
  //Initialising Channels for Release Sending and Receiving
  for i,_ := range ch3{
    ch3[i] = make(chan Pair)
  }

  ch := make([] chan bool,N)
  //Initialising Channels for Signal Sending and Receiving
  for i,_ := range ch{
    ch[i] = make(chan bool)
  }
  
  //Initialising List of Processes
  PList := make([] *Process,N)
  for  i:= 0;i < N ; i++{
    PList[i] = NewProcess(i)
  }
  
  for i:=0;i<N;i++{
    go PList[i].RecvRequest(ch1,ch2)
    go PList[i].RecvRelease(ch3,ch)
  }

  var wg1 sync.WaitGroup
  wg1.Add(1)
  go PList[0].AcquireLock(ch1,ch2,ch3,ch,&wg1)
  ch[0] <-true
  wg1.Add(1)
  go PList[1].AcquireLock(ch1,ch2,ch3,ch,&wg1)
  wg1.Add(1)
  go PList[2].AcquireLock(ch1,ch2,ch3,ch,&wg1)
  wg1.Add(1)
  go PList[3].AcquireLock(ch1,ch2,ch3,ch,&wg1)
  wg1.Wait()
}
