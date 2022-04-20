// package main

import (
	// "encoding/json"

	"fmt"
	"log"
	"os/exec"
)


func breakFile() {
	// ffmpeg -y -i DearZindagi.mkv -codec copy -map 0 -f segment -segment_time 7 -segment_format mpegts -segment_list DearZindagi_index.m3u8 -segment_list_type m3u8 ./segment_no_%d.ts

	cmd := exec.Command("ffmpeg", "-y" , "-i" , "./video/DearZindagi.mkv", "-codec", "copy", "-map", "0","-f", "segment", "-segment_time", "10", "-segment_format", "mpegts", "-segment_list", "D:/segments/DearZindagi_index.m3u8", "-segment_list_type", "m3u8", "D:/segments/segment_no_%d.ts")

	output, err := cmd.CombinedOutput()

	// err := cmd.Run()
	if err != nil {
		fmt.Printf("%s\n", output)
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", output)
	}
}

func main() {
	breakFile()
}