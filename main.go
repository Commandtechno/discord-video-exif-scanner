package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Channel struct {
	ID    string `json:"id"`
	Guild Guild  `json:"guild"`
}

type Guild struct {
	ID string `json:"id"`
}

func getMessageUrl(guildID string, channelID, messageID string) string {
	if guildID == "" {
		guildID = "@me"
	}

	return "https://discord.com/channels/" + guildID + "/" + channelID + "/" + messageID
}

func main() {
	packageDir := os.Args[1]
	channels, err := os.ReadDir(filepath.Join(packageDir, "messages"))
	if err != nil {
		fmt.Println("Could not read messages directory", err)
	}

	pending := 0
	wg := &sync.WaitGroup{}

	for _, channelDir := range channels {
		if !channelDir.IsDir() {
			continue
		}

		channelFile, err := os.Open(filepath.Join(packageDir, "messages", channelDir.Name(), "channel.json"))
		if err != nil {
			fmt.Println("Could not open channel file", err)
			continue
		}

		var channel Channel
		json.NewDecoder(channelFile).Decode(&channel)

		messagesFile, err := os.Open(filepath.Join(packageDir, "messages", channelDir.Name(), "messages.csv"))
		if err != nil {
			fmt.Println("Could not open messages file", err)
			continue
		}

		messagesReader := csv.NewReader(messagesFile)

		for {
			record, err := messagesReader.Read()
			if err != nil {
				break
			}

			id := record[0]
			attachments := record[3]

			if attachments != "" && attachments != "Attachments" {
				for _, attachment := range strings.Split(attachments, " ") {
					t := mime.TypeByExtension(filepath.Ext(attachment))
					if !strings.HasPrefix(t, "image/") && !strings.HasPrefix(t, "video/") && !strings.HasPrefix(t, "audio/") {
						continue
					}

					for pending > 100 {
						time.Sleep(time.Second)
					}

					pending++
					wg.Add(1)

					go func(attachment string) {
						defer func() {
							pending--
							wg.Done()
						}()

						cmd := exec.Command("ffmpeg", "-i", attachment, "-f", "ffmetadata", "-")
						metadata, err := cmd.Output()
						if err != nil {
							fmt.Println("Could not get metadata for", attachment, err)
							return
						}

						if strings.Contains(string(metadata), "location") {
							fmt.Println("⚠ Found message with location metadata")
							fmt.Println("> Attachment URL:", attachment)
							fmt.Println("> Message URL:", getMessageUrl(channel.Guild.ID, channel.ID, id))
							fmt.Println()
						}
					}(attachment)
				}
			}
		}
	}
}
