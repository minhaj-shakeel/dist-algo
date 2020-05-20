/*Distributed Algorithm to Find Minimum Spanning Tree*/

package main

import (
  "fmt"
  "os"
  "bufio"
  "strconv"
  "strings"
)

type Node struct{
  id int
  EdgeList [] int
  EdgeWeight []int
  status [] string
  level int
  state string
  rec int
}


func NewNode(i int) *Node{
  return &Node{id:i,
               EdgeList: make([] int,0),
               EdgeWeight : make([]int,0),
               status : make([]string,0),
               level : 0,
               state : "",
               rec : 0,
              }
}

//Minimum weighted edge in the Graph
func FindMin(){
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
    
    NodeList[n1].EdgeList = append(NodeList[n1].EdgeList,n2)
    NodeList[n1].EdgeWeight = append(NodeList[n1].EdgeWeight,w)
    
    NodeList[n2].EdgeList = append(NodeList[n2].EdgeList,n1)
    NodeList[n2].EdgeWeight = append(NodeList[n2].EdgeWeight,w)

  }
  fmt.Println("Nodes initialised")
  file.Close()
}
