package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	apiURL     = "https://api.rhododevdron.swuwu.dk"
	helpString = `ITU-Minitwit Tweet Flagging Tool


Usage:
	flag_tool <tweet_id>
	flag_tool -i
	flag_tool -h

Options:
-h            Show this screen.
-i            Dump all tweets and authors to STDOUT.`
)

type Args struct {
	MsgId       int
	Help        bool
	AllMessages bool
}

func flagMsgById(msgId int, c *http.Client) string {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/flag_tool/%d", apiURL, msgId), nil)
	if err != nil {
		log.Fatalf(err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	} else if resp.StatusCode != http.StatusBadRequest {
		return fmt.Sprintf("BadRequest - This message id: %d might not exist", msgId)
	}
	return fmt.Sprintf("Flagged entry: %d", msgId)
}

func getAllMessages(c *http.Client) []Message {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/flag_tool/msgs", apiURL), nil)
	if err != nil {
		log.Fatalf(err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		log.Fatalf(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf(err.Error())
	}
	resp.Body.Close()

	var messages []Message
	json.Unmarshal(body, &messages)

	return messages
}

func retrieveArgs() Args {
	args := Args{}
	help := flag.Bool("h", false, "print help")
	messages := flag.Bool("i", false, "print all tweets")

	flag.Parse()

	if *help {
		args.Help = true
	} else if *messages {
		args.AllMessages = true
	} else {
		PotentialMsgId := flag.Arg(0)
		PotentialMsgId = strings.Trim(PotentialMsgId, " ")
		msgId, _ := strconv.Atoi(PotentialMsgId)
		args.MsgId = msgId
	}
	return args
}

func main() {
	args := retrieveArgs()

	if args.Help {
		fmt.Println(helpString)
		return
	}

	client := http.Client{}

	if args.AllMessages {
		messages := getAllMessages(&client)
		for _, msg := range messages {
			fmt.Println(msg)
		}
		return
	}
	if args.MsgId != 0 {
		response := flagMsgById(args.MsgId, &client)
		fmt.Println(response)
	}
}

//helper
func (m Message) String() string {
	return fmt.Sprintf("Author: %s PubDate: %s Text: %s Flagged: %t", m.Author.Username, time.Unix(m.PubDate, 0).String(), m.Text, m.Flagged)
}

// straight copy pasta from Minitwit.models
type User struct {
	UserId       int64
	Username     string
	Email        string
	PasswordHash string
}

type Message struct {
	Author  User
	PubDate int64  // The publish timestamp as UNIX
	Text    string // The message itself
	Flagged bool
}
