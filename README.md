# fragbot
A very simple and straightforward Hypixel SkyBlock FragBot library made in Go<br>
As of now it only works with Microsoft accounts

## What is a fragbot? 
The video linked below explains everything you need to know perfectly<br>
https://youtu.be/EWBtOhZOZxQ

## TODO
- Better error handling and more failsafes
- Mojang account support (`auth`-type field in config) 
- Add comments

## Config
- `username` string, username of the fragbot
- `email`    string, fragbot's email address
- `password` string, password of the fragbot
- `webhook`  string, optional, webhook url
- `blacklist` []string, can be empty, list of blacklisted usernames
- `message`  string, optional, message sent in pc after fragbot joins the party

## Events (Handlers)
- `join` happens when bot joins hypixel successfully
- `end`  happens when bot gets disconnected 
- `invite` happens when bot gets added to the party
- `limbo`  happens when bot gets sent to Limbo because AFK
- `dungeons` happens when player joins the dungeon
## Installation
`go get github.com/nextu1337/fragbot`

## Example use

```go
package main

import (
	"github.com/nextu1337/fragbot"
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
