/* O(nlogn) algorithm for Leader Election in Rings*/
/*Process with Smallest Uid is selected As Leader*/

package main
import (
  "fmt"
  "math/rand"
  "sync"
)

var wg sync.WaitGroup

func pow(a,k int) int{
  if k==0{
    return 1
  } else if k==1{
    return a
  }
  return pow(a*a,k-1)
}

type msg struct{
  msgtype string
  j,k,d int
}


type Process struct{
  uid int
  electedUID int
  isLeader bool
  LeftSendChannel chan msg
  LeftRecvChannel chan msg
  RightSendChannel chan msg
  RightRecvChannel chan msg
}

func NewProcess(i int) *Process{
  return &Process{uid:i,
                  isLeader:false,
                }
}

/*Receiving and responding from Left Adjacent Process*/
func (p *Process)ReceiveLeft(ch chan int){
  for{
    msgLeft := <-p.LeftRecvChannel
    if p.uid ==0{
    }
    if msgLeft.msgtype=="probe"{
      if msgLeft.j==p.uid{
        p.isLeader=true
        ch<-0
        break
      }
      if msgLeft.j<p.uid && msgLeft.d<pow(2,msgLeft.k){
        msgLeft.d++
        p.RightSendChannel<-msgLeft
      }else if msgLeft.j<p.uid && msgLeft.d==pow(2,msgLeft.k){
        p.LeftSendChannel<-msg{"reply",msgLeft.j,msgLeft.k,0}
      }
    }else if msgLeft.msgtype=="reply"{
      if msgLeft.j!=p.uid{
        p.RightSendChannel<-msgLeft
      }else{
      //Signal to start next round
        ch<-1
      }
    }
  }
}

/*Receiving and responding from Right Adjacent Process*/
func (p *Process)ReceiveRight(ch chan int){
 for{
   msgRight :=<-p.RightRecvChannel
   if p.uid ==0{
   }
   if msgRight.msgtype == "probe"{
     if msgRight.j==p.uid{
       p.isLeader=true
       ch<-0
       break
     }
     if msgRight.j<p.uid && msgRight.d<pow(2,msgRight.k){
       msgRight.d++
       p.LeftSendChannel<-msgRight
     }else if msgRight.j<p.uid &&msgRight.d==pow(2,msgRight.k){
       p.RightSendChannel<-msg{"reply",msgRight.j,msgRight.k,0}
     }
   } else if msgRight.msgtype == "reply"{
     if msgRight.j!=p.uid{ 
       p.LeftSendChannel<-msgRight
     }else{
       //Signal to start Next Round
       ch<-1
     }
   }
  }
}

func (p *Process)StartRound(ch chan int){
  for i :=0; ;i++{
    fmt.Println("Process ",p.uid,"starting round: ",i)
    p.LeftSendChannel<-msg{"probe",p.uid,i,1}
    p.RightSendChannel<-msg{"probe",p.uid,i,1}
    signal1:=<-ch
    signal2:=<-ch
    if signal1==0 && signal2==0{
      fmt.Println("Leader Elected! Process uid: ",p.uid)
      break
    }
  }
}
func (p *Process) StartElection(wg *sync.WaitGroup){
  ch1 := make(chan int)
  go p.ReceiveLeft(ch1)
  go p.ReceiveRight(ch1)
  p.StartRound(ch1)
  defer wg.Done()
}


func main(){
  N:=5
  uidList:=rand.Perm(N)
  PList := make([]*Process,N)
  for i:=0;i<N;i++{
    PList[i] = NewProcess(uidList[i])
    PList[i].LeftSendChannel = make(chan msg)
    PList[i].RightSendChannel= make(chan msg)
  }
  
  for i:=0;i<N;i++{
    PList[i].RightRecvChannel=PList[(i+1)%N].LeftSendChannel
    PList[(i+1)%N].LeftRecvChannel=PList[i].RightSendChannel
  }

  wg.Add(1)
  for i:=0;i<N;i++{
    go PList[i].StartElection(&wg)
  }
  wg.Wait()
}
