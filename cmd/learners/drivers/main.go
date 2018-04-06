package main

import (
    "bufio"
    "fmt"
    "net"
)

func handleStreamIn(userProgConn net.Conn, inputConn net.Conn) {
    for {
        streamIn, err := bufio.NewReader(inputConn).ReadString('\n')
        if err != nil {
            panic(err)
        }
        fmt.Fprintf(userProgConn, streamIn)
    }
}

func handleStreamOut(userProgConn net.Conn, inputConn net.Conn) {
    for {
        streamOut, err := bufio.NewReader(userProgConn).ReadString('\n')
        if err != nil {
            panic(err)
        }
        fmt.Fprintf(inputConn, streamOut)
    }
}

func main() {

    userProgConn, err := net.Dial("tcp", ":9999")
    if err != nil {
        panic(err)
    }
    
    ln, err := net.Listen("tcp", ":8888")
    if err != nil {
        panic(err)
    }

    for {
        inputConn, err := ln.Accept()
        if err != nil {
            panic(err)
        }
        go handleStreamIn(userProgConn, inputConn)
        go handleStreamOut(userProgConn, inputConn)
    }

}