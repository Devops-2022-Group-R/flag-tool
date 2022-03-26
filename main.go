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
	apiUrl           = "http://localhost:8080"
	apiProductionUrl = "https://api.rhododevdron.swuwu.dk"
	helpString       = `ITU-Minitwit Tweet Flagging Tool


Usage:
	flag_tool <tweet_id>
	flag_tool -i
	flag_tool -h
	flag_tool -u <username> -pwd <password>

Options:
-h            Show this screen.
-i            Dump all tweets and authors to STDOUT.
-p 		Target the production url. 
-u 		Username
-pwd		Password`
)

type Args struct {
	MsgId        int
	Help         bool
	AllMessages  bool
	IsProduction bool
	Username     string
	Password     string
}

func retrieveArgs() Args {
	args := Args{}
	help := flag.Bool("h", false, "print help")
	messages := flag.Bool("i", false, "print all tweets")
	isProduction := flag.Bool("p", false, "Target production api url")
	username := flag.String("u", "", "Username")
	password := flag.String("pwd", "", "Password")

	flag.Parse()

	if *isProduction {
		args.IsProduction = true
	}

	args.Username = *username
	args.Password = *password

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
	url := apiUrl
	if args.IsProduction {
		url = apiProductionUrl
	}

	if args.AllMessages {
		messages := getAllMessages(&client, url)
		for _, msg := range messages {
			fmt.Println(msg)
		}
		return
	}
	if args.MsgId != 0 {
		response := flagMsgById(args.MsgId, &client, url)
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
