package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	apiURL                = "https://api.rhododevdron.swuwu.dk"
	tweetArrayInitialSize = 5
	helpString            = `ITU-Minitwit Tweet Flagging Tool


Usage:
	flag_tool <tweet_id>
	flag_tool -i
	flag_tool -h

Options:
-h            Show this screen.
-i            Dump all tweets and authors to STDOUT.`
)

type Args struct {
	TweetId   int
	Help      bool
	AllTweets bool
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

func getAllTweets(c *http.Client) []string {
	tweets := make([]string, tweetArrayInitialSize)

	return tweets
}

func retrieveArgs() Args {
	args := Args{}
	help := flag.Bool("h", false, "print help")
	tweets := flag.Bool("i", false, "print all tweets")

	flag.Parse()

	if *help {
		args.Help = true
	} else if *tweets {
		args.AllTweets = true
	} else {
		PotentialTweetId := flag.Arg(0)
		PotentialTweetId = strings.Trim(PotentialTweetId, " ")
		tweetId, _ := strconv.Atoi(PotentialTweetId)
		args.TweetId = tweetId
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

	if args.AllTweets {
		tweets := getAllTweets(&client)
		for _, tweet := range tweets {
			fmt.Println(tweet)
		}
		return
	}
	if args.TweetId != 0 {
		response := flagMsgById(args.TweetId, &client)
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
