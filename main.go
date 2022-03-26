package main

import (
	"encoding/base64"
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
	encodedCrendetials := encodeCredentialsToB64(args.Username, args.Password)
	url := apiUrl

	if args.IsProduction {
		url = apiProductionUrl
	}

	if args.AllMessages {
		messages := getAllMessages(&client, url, encodedCrendetials)
		for _, msg := range messages {
			fmt.Println(msg)
		}
		return
	}

	if args.MsgId != 0 {
		response := flagMsgById(args.MsgId, &client, url, encodedCrendetials)
		fmt.Println(response)
	}
}

func flagMsgById(msgId int, c *http.Client, url, encodedeCredentials string) string {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/flag_tool/%d", url, msgId), nil)
	req = SetRequestHeader(encodedeCredentials, *req)
	if err != nil {
		log.Fatalf(err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	} else if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("BadRequest - This message id: %d might not exist", msgId)
	}
	return fmt.Sprintf("Flagged entry: %d", msgId)
}

func getAllMessages(c *http.Client, url, encodedeCredentials string) []Message {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/flag_tool/msgs", url), nil)
	req = SetRequestHeader(encodedeCredentials, *req)
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

//helper toString
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

func encodeCredentialsToB64(username string, password string) string {
	data := username + ":" + password
	sEnc := base64.StdEncoding.EncodeToString([]byte(data))
	return sEnc
}

func SetRequestHeader(encodededCredentials string, req http.Request) *http.Request {
	req.Header.Set("Authorization", "Basic "+encodededCredentials)
	return &req
}
