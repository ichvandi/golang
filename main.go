package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	var err error
	var client = &http.Client{}

	request, err := http.NewRequest("GET", "https://cast4.my-control-panel.com/proxy/richard7/stream", nil)
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("User-Agent", "VLC/3.0.6 LibVLC/3.0.6")
	request.Header.Set("Icy-MetaData", "1")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	icyHeader := response.Header.Get("icy-metaint")
	if icyHeader == "" {
		log.Fatal("IceCast server doesn't support metadata")
		defer response.Body.Close()
	}

	icyMetaInt, err := strconv.Atoi(icyHeader)
	if err != nil {
		log.Fatal(err)
	}

	var x []byte

	reader := bufio.NewReader(response.Body)
	for {
		// fmt.Println("A")
		skipped, err := io.CopyN(io.Discard, reader, int64(icyMetaInt))
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println("B")
		for skipped != int64(icyMetaInt) {
			newSkipped, err := io.CopyN(io.Discard, reader, int64(icyMetaInt-int(skipped)))
			if err != nil {
				log.Fatal(err)
			}
			skipped += newSkipped
		}

		// fmt.Println("C")
		symbolLength, err := reader.ReadByte()
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println("D")
		metadataLength := symbolLength * 16
		if metadataLength > 0 {
			// fmt.Println("E")
			for i := 0; i < int(metadataLength); i++ {
				metadataSymbol, err := reader.ReadByte()
				if err != nil {
					log.Fatal(err)
				}

				if metadataSymbol > 0 {
					x = append(x, metadataSymbol)
				}
			}

			streamTitle := string(x[:])
			x = []byte{}
			fmt.Println("Before: " + streamTitle)
			fmt.Println("After: " + parseMetadata(streamTitle))
		}
	}
}

func parseMetadata(result string) string {
	return strings.TrimLeft(strings.TrimRight(result, "';"), "StreamTitle='")
}
