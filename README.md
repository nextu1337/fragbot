# go-fragbot
A very simple and straightforward Hypixel SkyBlock FragBot library made in Go

## Installation
`go get github.com/nextu1337/fragbot`

## Example use

```go
package main

import (
	fragbot "github.com/nextu1337/fragbot"
)

func main() {
	config := map[string]interface{}{"message":"join gg/url for more FragBots!", // Message sent in party chat after someone joins [OPTIONAL]
  	"username":"FRAGBOT_username", // Fragbot's username
  	"email":"example@outlook.com", // Fragbot's e-mail
  	"password":"example_password", // Fragbot's password
  	"webhook":"https://discord.com/api/webhooks/XXXXXXXXXXXXXXXXXX/XXXXXXXX...",  // Webhook URL [OPTIONAL]
 	 "blacklist": []string{}} // Blacklist, can be empty
  
	frag := fragbot.New(config)

	frag.Handlers["join"] = func() {
		fmt.Println("FragBot connected to the server");
	}

	frag.Handlers["end"] = func() {
		fmt.Println("FragBot disconnected from the server");
	}

	frag.Handlers["invite"] = func() {
		fmt.Sprintln("FragBot was invited by %v",frag.Queue[len(frag.Queue) - 1]);
	}
  	frag.SetMessage("limbo","Bot was sent to the limbo")
  
	frag.Start() // Infinite loop starts here
}
```
