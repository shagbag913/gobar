package main

import (
    "bytes"
    "time"
    "strconv"
    "os"
    "fmt"
)

func getVolumeGlyph(percentage int, muted bool) string {
    glyphs := []string{"", "", "", ""}

    switch {
    case muted:
        return glyphs[0]
    case percentage == 0:
        return glyphs[1]
    case percentage <= 50:
        return glyphs[2]
    default:
        return glyphs[3]
    }
}


func setVolumeString() {
    volTempPath := os.Getenv("HOME") + "/.cache/volume/percentage"

    for {
        /*
         * Fetch volume from a temp file, so we don't have to poll
         * from ALSA or some wrapper constantly
         */
        file, err := os.Open(volTempPath)
        if err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            break
        }
        volFromFile := make([]byte, 4)
        var num int
        num, err = file.Read(volFromFile)
        if err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            break
        }

        /* If 'M' succeeds the percentage, we're muted */
        muted := false
        if bytes.Contains(volFromFile, []byte("M")) {
            muted = true
        } else {
            /* Append random char to end so later logic is simpler */
            volFromFile = append(volFromFile, 'N')
            num++
        }

        percentage, err := strconv.Atoi(string(volFromFile[:num-1]))
        if err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            break
        }

        glyph := getVolumeGlyph(percentage, muted)
        newVolumeString := glyph + " " + strconv.Itoa(percentage) + "%"

        if newVolumeString != volumeString {
            volumeString = newVolumeString
            printBuffer()
        }
        time.Sleep(time.Second * 1)
    }
}

