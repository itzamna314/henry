henry
=====
A simple, no bs Slackbot in Go

## Quickstart

![screen shot 2016-07-14 at 11 53 51
pm](https://cloud.githubusercontent.com/assets/4248167/16866056/4e0f72a8-4a1e-11e6-950e-584af50c6426.png)


```go
package main

import (
    "fmt"
    "github.com/bndw/henry"
)

func main() {
    // Create a bot
    bot, err := henry.Create("YOUR_SLACK_API_TOKEN")
    if err != nil {
        panic(err)
    }

    // Add a handler
    bot.Handle("weather", weatherHandler)

    // Start listening
    if err = bot.Listen(); err != nil {
        panic(err)
    }
}

func weatherHandler(message *henry.Message) string {
    if len(message.Args) < 1 {
        return "the weather is great!"
    }
  
    day := message.Args[0]

    return fmt.Sprintf("The weather %s is great!", day)
}
```

## License
MIT

## Credits
_Inspired by [mybot](https://www.opsdash.com/blog/slack-bot-in-golang.html)_
