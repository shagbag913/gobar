package main

import (
    "os"
    "io/ioutil"
    "fmt"
    "time"
)

func setNetStatus() {
    var newNetStatus string
    for {
        baseDir := "/sys/class/net/"
        netDirs, err := ioutil.ReadDir(baseDir)
        if err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            break
        }

        newNetStatus = ""

        for _, netDir := range netDirs {
            file, err := os.Open(baseDir + netDir.Name() + "/operstate")
            if err != nil {
                fmt.Fprintln(os.Stderr, err.Error())
                continue
            }

            state := make([]byte, 4)

            count, err := file.Read(state)
            if err != nil {
                fmt.Println(err.Error())
                file.Close()
                continue
            }

            if count == 3 {
                if _, err := os.Stat(baseDir + netDir.Name() + "/wireless"); os.IsNotExist(err) {
                    newNetStatus += "   "
                } else {
                    newNetStatus += "   "
                }
            }

            file.Close()
        }

        if newNetStatus != "" {
            newNetStatus = newNetStatus[:len(newNetStatus)-3]
        }

        if newNetStatus != netStatus {
            netStatus = newNetStatus
            printBuffer()
        }

        time.Sleep(time.Second * 5)
    }
}
