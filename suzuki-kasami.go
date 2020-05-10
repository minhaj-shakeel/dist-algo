package main

import  (
  "fmt"
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

func enqueue(queue[] int, element int) []int {
  queue = append(queue, element);
  return queue
}

func dequeue(queue[] int) ([]int) {
  return queue[1:];
}

func isEmpty(queue[] int) bool {
    return len(queue)==0
}

func isPresent(queue[] int, i int) bool{
  found := false
  for _,v := range queue{
    if (i == v){
      found = true
      break
    }
  }
  return found
}

type Token struct{
  Q [] int
  Last [] int
}

func NewToken(N int) *Token{
  return &Token{Q : make([]int,0),
                Last : make([]int,N),
        }
}

type Process struct{
  index int
  HasToken bool
  inCS bool
  Request [] int
  token *Token
}

type Req struct {
  pid int
  seq int
}

func NewProcess(N int, i int) *Process{
  return &Process{index: i,
                  HasToken : false,
                  inCS : false,
                  Request : make([]int,N),
                }
}



func (p Process) ReleaseToken(ch [] chan *Token){
  curr:=p.index
  p.token.Last[curr]=p.Request[curr]
  for{
    for i,v := range p.Request{
      if (i!=curr) && (isPresent(p.token.Q,i)==false) && (v==p.token.Last[i]+1){
        p.token.Q = append(p.token.Q,i)
      }
    }
    if isEmpty(p.token.Q)==false{
      break
    }
  }

  next:=p.token.Q[0]
  p.token.Q = dequeue(p.token.Q)
  p.HasToken=false
  ch[next] <- p.token
}



func (p Process) SendRequest(ch [] chan  Req){
  curr := p.index
  p.Request[curr]+=1
  req:=Req{curr,p.Request[curr]}
  for i := 0 ; i < len(ch) ; i++ {
     go func(i int){ch[i] <- req }(i)
  }
}


func (p Process) ReceiveRequest(ch1 [] chan Req,ch2 [] chan *Token){
  curr :=p.index
  for {
    req := <-ch1[curr]
    p.Request[req.pid]=max(req.seq,p.Request[req.pid])
    if (p.HasToken == true) && (p.inCS == false){
      p.ReleaseToken(ch2)
    }
  }
}

func (p Process) DoSomething(){
  fmt.Printf("Process %d Entered\n",p.index)
  time.Sleep(5*time.Second)
  fmt.Printf("Process %d Exited\n",p.index)
}


func(p Process) EnterCS(ch1 [] chan Req,ch2 [] chan *Token,wg *sync.WaitGroup){
  defer wg.Done()
  curr := p.index
  p.SendRequest(ch1)
  p.token =<- ch2[curr]
  p.HasToken = true
  p.inCS = true
  p.DoSomething()
  p.inCS = false
  go p.ReleaseToken(ch2)

}


func main(){
  N:=5 //Number of Processes

  //Initialising Channels for Request Receiving and Sending
  ch1 := make([] chan Req,N) 
  for i,_:= range ch1{
    ch1[i] = make(chan Req)
  }
  //Initialising Channels for token Receiving and Sending
  ch2 := make([] chan *Token,N)
  for i,_ := range ch2{
    ch2[i] = make(chan *Token)
  }
  
  //Initialising List of Processes
  PList := make([] *Process ,5)
  for i:=0;i<5;i++{
    PList[i] = NewProcess(N,i)
  }
  
  //Giving Token to Process 0 initially
  token := NewToken(N)
  PList[0].HasToken=true
  PList[0].token = token
  
  //Starting to Receiving Requests from other Processes
  for i:=0;i<5;i++{
   go PList[i].ReceiveRequest(ch1,ch2)
  }
  
  var wg sync.WaitGroup
  for{
    var pid int
    fmt.Println("Enter Pid")
    fmt.Scan(&pid)
    if pid==-1{
      break
    }
    wg.Add(1)
    go PList[pid].EnterCS(ch1,ch2,&wg)
  }
  wg.Wait()
}
