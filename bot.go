package henry

import (
	"fmt"
	"strings"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

type Bot struct {
	id       string             // UserId from slack once we're connected
	conn     *websocket.Conn    // Websocket conn to slack
	token    string             // Slack API Key
	handlers map[string]Handler // User registered handlers
	counter  uint64             // Tracks the count of messages sent
}

// Create connects to Slack and returns a Bot
func Create(apiToken string) (*Bot, error) {
	conn, id, err := slackConnect(apiToken)
	if err != nil {
		return nil, err
	}

	bot := Bot{
		id:       id,
		conn:     conn,
		token:    apiToken,
		handlers: make(map[string]Handler),
		counter:  0,
	}

	return &bot, nil
}

// Listen starts the bot in the listening loop
func (bot *Bot) Listen() error {
	fmt.Println("Listening for messages...")

	for {
		m, err := bot.receive()
		if err != nil {
			fmt.Println(err.Error())
			bot.reply(m, err.Error())
			continue
		}

		if m == nil {
			continue
		}

		handler, ok := bot.handlers[m.Command]
		if !ok {
			bot.reply(m, fmt.Sprintf("No handler for %s", m.Command))
			continue
		}

		// TODO: Don't let this blow up
		resp := handler.Fn(m)

		bot.reply(m, resp)
	}
}

type Message struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	User    string `json:"user"`
	Command string
	Args    []string
}

type Handler struct {
	Fn func(*Message) string
}

// Handle registers a handler for a command
func (bot *Bot) Handle(cmd string, handler func(*Message) string) {
	bot.handlers[cmd] = Handler{Fn: handler}
}

// receive receives a message from Slack
func (bot *Bot) receive() (*Message, error) {

	var msg Message
	err := websocket.JSON.Receive(bot.conn, &msg)
	if err != nil {
		return nil, err
	}

	// it's not a message or not for us
	if msg.Type != "message" || !strings.HasPrefix(msg.Text, "<@"+bot.id+">") {
		return nil, nil
	}

	parts := strings.Fields(msg.Text)
	// [<@id>, deploy, foo/bar:123, beta]

	if len(parts) <= 1 {
		return nil, fmt.Errorf("Unsupported command%s", "")
	}

	msg.Command = parts[1]
	msg.Args = parts[2:]

	return &msg, nil
}

// reply sends a response to Slack
func (bot *Bot) reply(msg *Message, text string) {
	go func(msg *Message) {
		msg.Id = atomic.AddUint64(&bot.counter, 1)
		msg.Text = text
		_ = websocket.JSON.Send(bot.conn, msg)
	}(msg)
}
