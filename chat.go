package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	verbose := flag.Bool("verbose", false, "enable verbose logging")
	flag.Parse()

	if *verbose {
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Verbose logging enabled")
	} else {
		log.SetOutput(os.Stdout)
		log.SetFlags(0)
		log.SetPrefix("")
	}
	client := anthropic.NewClient()

	scanner := bufio.NewScanner(os.Stdin)
	// get user message
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}

	// create the agent and run it
	agent := NewAgent(&client, getUserMessage, *verbose)
	err := agent.Run(context.TODO())
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}
}

func NewAgent(client *anthropic.Client, getUserMessage func() (string, bool), verbose bool) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
		verbose:        verbose,
	}
}

type Agent struct {
	client         *anthropic.Client
	getUserMessage func() (string, bool)
	verbose        bool
}

func (a *Agent) Run(ctx context.Context) error {
	conversation := []anthropic.MessageParam{}

	if a.verbose {
		log.Println("Starting conversation")
	}
	fmt.Println("Chat with acron!")

	for {
		fmt.Print("\u001b[94mYou\u001b[0m: ")
		userInput, ok := a.getUserMessage()
		if !ok {
			if a.verbose {
				log.Println("User input ended, breaking from chat loop")
			}
			break
		}

		// Skip empty messages
		if userInput == "" {
			if a.verbose {
				log.Println("Skipping empty message")
			}
			continue
		}

		if a.verbose {
			log.Printf("User input received: %q", userInput)
		}

		userMessage := anthropic.NewUserMessage(anthropic.NewTextBlock(userInput))
		conversation = append(conversation, userMessage)

		if a.verbose {
			log.Printf("Sending message to Claude, conversation length: %d", len(conversation))
		}

		message, err := a.sendMessage(ctx, conversation)
		if err != nil {
			if a.verbose {
				log.Printf("Error sending message: %v", err)
			}
			break
		}

		conversation = append(conversation, message.ToParam())

		if a.verbose {
			log.Printf("Received response from Claude with %d content blocks", len(message.Content))
		}

		fmt.Printf("\u001b[92mAcron\u001b[0m: %s\n", message.Content[0].Text)
	}
	if a.verbose {
		log.Println("Chat session ended")
	}
	return nil
}

func (a *Agent) sendMessage(ctx context.Context, conversation []anthropic.MessageParam) (*anthropic.Message, error) {
	if a.verbose {
		log.Println("Sending the message with %s modle", anthropic.ModelClaude3_7SonnetLatest)
	}

	message, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_7SonnetLatest,
		MaxTokens: int64(1024),
		Messages:  conversation,
	})

	if a.verbose {
		if err != nil {
			log.Printf("API call failed: %v", err)
		} else {
			log.Printf("API call successful, response received")
		}
	}

	return message, err
}
