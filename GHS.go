/*Distributed Algorithm to Find Minimum Spanning Tree*/

package main

import (
  "fmt"
  "os"
  "bufio"
  "strconv"
  "strings"
)
type msg struct{
  Type string
  level int
  name string
  state string
}

type Node struct{
  name string
  id int
  EdgeList [] int
  EdgeWeight []int
  status [] string
  level int
  state string
  rec int
  parent int
  SendChannel [] chan msg
  RecvChannel [] chan msg

  /*Change the types According to Low level Implementation*/
  bestNode int
  bestWt int
  testNode int
}


func NewNode(i int) *Node{
  return &Node{id:i,
               EdgeList: make([] int,0),
               EdgeWeight : make([]int,0),
               status : make([]string,0),
               level : 0,
               state : "SLEEP",
               rec : 0,
               SendChannel : make([] chan msg,0),
               RecvChannel : make([] chan msg,0),
              }
}

//Minimum weighted edge in the Graph
func (n *Node) FindMin(){
}

func minIndex(EdgeWeight [] int) int{
  min:=0
  for i,v:=range EdgeWeight{
    if v<EdgeWeight[min]{
      min = i 
    }
  }
  return min
}

func (n *Node)Init(){
  minI := minIndex(n.EdgeWeight)
  n.status[minI] = "branch"
  n.SendChannel[minI]<-msg{"INITIATE",0}
  msgRecv:=<-RecvChannel[minI]


}
func (n *Node) RecvMsg(){
  for q:=0;q<len(n.RecvChannel);q++{
    go func(i int){ msgRecv:=<-RecvChannel[i]
                    if msgRecv.name=="connect" {
                      
                    }    
                  }(q)
  } 

}
func (n *Node) RecvConnect(msgRecv msg,q int){
  L:=msgRecv.level
  if L<n.level{
    /*Combine with Rule LT*/
    n.status[q]="branch"
    n.SendChannel[q]<-msg{"initiate",n.level,n.name,n.state}
  } else if status[q]=="basic"{
    /*wait*/
  } else {
    /*Combine with rule EQ*/
    newname:=strconv.Itoa(n.Id)+strconv.Itoa(n.EdgeList[q])
    n.SendChannel[q]<-msg{"initiate",n.level+1,newname,"find")}
  }
}

func (n *Node) RecvInitiate(msgRecv msg , q int){
  n.level = msgRecv.level
  n.name  = msgRecv.name
  n.state = msgRecv.state
  n.parent = q
  
  n.bestNode = -1
  n.bestWt = 100000
  n.testNode = -1

  for r:=0;r<len(n.SendChannel);r++{
    if n.status[r]=="branch" && r!=q{
      //change to concurrent send
      n.SendChannel[r]<-msgRecv
    }
  }
  if state =="find"{
    n.rec=0
    n.FindMin()
  }
}

func main(){
  file,_ := os.Open("sample-input.txt")
  scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)
  
  scanner.Scan()
  N,_:=strconv.Atoi(scanner.Text())
  NodeList := make([] *Node ,N)
  for i:=0;i<N;i++{
    NodeList[i] = NewNode(i)
  }

  for scanner.Scan(){
    txtline:=strings.Split(scanner.Text()," ")
    
    n1,_:=strconv.Atoi(txtline[0])
    n2,_:=strconv.Atoi(txtline[1])
    w,_:=strconv.Atoi(txtline[2])
    
    ch1 := make(chan msg)
    ch2 := make(chan msg)
    
    NodeList[n1].EdgeList = append(NodeList[n1].EdgeList,n2)
    NodeList[n1].EdgeWeight = append(NodeList[n1].EdgeWeight,w)
    NodeList[n1].status = append(NodeList[n1].status,"UNUSED")
    NodeList[n1].SendChannel = append(NodeList[n1].SendChannel,ch1)
    NodeList[n1].RecvChannel = append(NodeList[n1].RecvChannel,ch2)
    
    NodeList[n2].EdgeList = append(NodeList[n2].EdgeList,n1)
    NodeList[n2].EdgeWeight = append(NodeList[n2].EdgeWeight,w)
    NodeList[n2].status = append(NodeList[n2].status,"UNUSED")
    NodeList[n2].SendChannel = append(NodeList[n2].SendChannel,ch2)
    NodeList[n2].RecvChannel = append(NodeList[n2].RecvChannel,ch1)

  }

  file.Close()
}
