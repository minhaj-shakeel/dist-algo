/* O(nlogn) algorithm for Leader Election in Rings*/
package main
import (
  "fmt"
)

type msg{
  msgtype string
  j,k,d int
}


type Process struct{
  uid int
  electedUID int
  isLeader bool
  state string
  LeftSendChannel chan msg
  LeftRecvChannel chan msg
  RightSendChannel chan msg
  RightRecvChannel chan msg
}

func main(){
  fmt.Println("Hirchberg-Sinclair")
}
