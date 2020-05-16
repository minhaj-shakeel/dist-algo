/*Chang Roberts Algorithm For Leader Election in Destributed Network*/
package main
import (
    "fmt"
    "math/rand"
    "sync"
)
var wg sync.WaitGroup

type Process struct{
  uid int
  electedUID int
  isLeader bool
  state string
  SendChannel chan msg
  RecvChannel chan msg
}

type msg struct {
  uid int
  status string
}


func NewProcess(i int) *Process{
  return &Process{uid:i,
                  isLeader:false,
                  state : "non-participant",
                }
}

func (p *Process)Initiate(){
  p.SendChannel<-msg{p.uid,"Find"}
  p.state = "participant"
}

func (p *Process)RecvMsg(wg *sync.WaitGroup){
  for{
    message:=<-p.RecvChannel
    if message.status == "LEADER"{
      if p.uid ==message.uid{
        break
      }else{
        p.electedUID = message.uid
        p.state = "non-participant"
        p.SendChannel<-message
        break
      }
    }

    if message.uid>p.uid{
      p.SendChannel<-message
    }else if message.uid<p.uid{
      if p.state == "non-participant"{
        message.uid = p.uid
        p.SendChannel<-message
        p.state = "participant"
      }
    }else{
      fmt.Println("Leader Elected",p.uid)
        p.state = "non-participant"
        p.SendChannel<-msg{p.uid,"LEADER"}
    }
  }
  wg.Done()
}



func main(){
  N:=50
  uidList:=rand.Perm(N)
  PList := make([]*Process,N)
  for i:=0;i<N;i++{
    PList[i] = NewProcess(uidList[i])
    PList[i].SendChannel = make(chan msg)
  }
  
  //Setting Up receiving channel
  PList[0].RecvChannel = PList[N-1].SendChannel
  for i:=1;i<N;i++{
    PList[i].RecvChannel=PList[i-1].SendChannel
  }
  
  for i:=0;i<N;i++{
    wg.Add(1)
    go PList[i].RecvMsg(&wg)
  }

  /*Any Process can initiate the Process*/
  PList[6].Initiate()
  fmt.Println("Election Completed")
  wg.Wait()


}

