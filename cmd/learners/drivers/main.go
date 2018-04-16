package main

import (
    "bufio"
    "fmt"
    "net"
    "log"
    "os"
    "os/exec"
    "os/signal"
    "syscall"

    "github.com/kelseyhightower/envconfig"
)

func signalHandler(cmd *exec.Cmd) {
    signal_chan := make(chan os.Signal, 1)
    signal.Notify(signal_chan, syscall.SIGINT)
    go func() {
        for {
            s := <-signal_chan
            switch s {
            case syscall.SIGINT:
                if err := cmd.Process.Kill(); err != nil {
                    panic(err)
                }
                os.Exit(0)
            }
        }
    }()
}

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

    driverConfig := config{}
    if err := envconfig.Process("driver", &driverConfig); err != nil {
        log.Fatalf("Error on processing configuration: %s", err.Error())
        return
    }

    cmd := exec.Command("python", "sgdclassifier.py")
    if err := cmd.Start(); err != nil {
        panic(err)
    }
    go signalHandler(cmd)

    userProgConn, err := net.Dial("tcp", driverConfig.StreamInURL)
    if err != nil {
        panic(err)
    }
    
    ln, err := net.Listen("tcp", driverConfig.StreamOutURL)
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