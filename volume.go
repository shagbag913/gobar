package main

import (
    "bytes"
    "time"
    "strconv"
    "os"
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
    file, err := os.Open(volTempPath)
    if logFatal(err) {
        return
    }
    defer file.Close()

    for {
        /*
         * Fetch volume from a temp file, so we don't have to poll
         * from ALSA or some wrapper constantly
         */
        _, err := file.Seek(0, 0)
        if logFatal(err) {
            break
        }
        volFromFile := make([]byte, 4)
        var num int
        num, err = file.Read(volFromFile)
        if logFatal(err) {
            break
        }

        /* If 'M' succeeds the percentage, we're muted */
        muted := false
        if bytes.Contains(volFromFile[:num], []byte("M")) {
            muted = true
            num--
        }

        percentage, err := strconv.Atoi(string(volFromFile[:num]))
        if logFatal(err) {
            break
        }

        glyph := getVolumeGlyph(percentage, muted)
        newVolumeString := glyph + " " + strconv.Itoa(percentage) + "%"

        if newVolumeString != volumeString {
            volumeString = newVolumeString
            printBuffer()
            time.Sleep(time.Millisecond * 100)
        } else {
            time.Sleep(time.Millisecond * 300)
        }
    }
}

