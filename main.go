package fragbot

import (
    "strings"
	"time"
	"errors"
	"sort"
	"fmt"
	mc "github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/basic"
	"github.com/Tnze/go-mc/bot/msg"
	aut "github.com/maxsupermanhd/go-mc-ms-auth"
	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/bot/playerlist"
)

func messageJoin(message chat.Message) string {
    toreturn := message.Text
    if message.Extra == nil {
        return toreturn
    }
    for _, extra := range message.Extra {
        toreturn += extra.Text
    }
    return toreturn
}

type FragBot struct {
	email		string
	username  	string
	password	string
	wh        	string
	pL			*playerlist.PlayerList
	bl 		  	[]string
	config    	map[string]interface{}
	Queue     	[]string
	Log 		func(...interface{})
	Current   	string
	client    	*mc.Client
	player	  	*basic.Player
	chat  		*msg.Manager
	messages  	map[string]string
	Handlers  	map[string]func()
}

func New(config map[string]interface{}) *FragBot {
	fb := &FragBot{
		username: 		  config["username"].(string),
		email: 			  config["email"].(string),
		password:		  config["password"].(string),
		wh:       		  config["webhook"].(string),
		bl:		  		  config["blacklist"].([]string),
		Log:			  log,
		config:   		  config,
		Queue:    		  make([]string, 0),
		messages: map[string]string{
			"join":     "Successfully connected to the server as %s",
			"end":      "%s was kicked from the server",
			"invite":   "%i invited %s to the party. His position in queue is %p",
			"joined":   "%s joined %i's party",
			"disband":  "%i didn't join the dungeons in time. L",
			"dungeons": "%i joined the dungeons. Leaving the party...",
			"disbanded": "%i disbanded the party.",
			"limbo":     "%s got sent to the Limbo. It's fully safe now",
		},
		Handlers: map[string]func(){},
	}

	fb.client = mc.NewClient()

	
	var err interface{}
	fb.client.Auth, err = aut.GetMCcredentials("","88650e7e-efee-4857-b9a9-cf580a00ef43")
	if err != nil {
	 	panic(err)
	}
	fb.player = basic.NewPlayer(fb.client, basic.DefaultSettings, basic.EventsListener{
		GameStart:  func() error {
			fmt.Printf(fb.tm(fb.messages["join"])+"\n")
			if fn, ok := fb.Handlers["join"]; ok {
				fn();
			}
			return nil
		},
		Disconnect: func(reason chat.Message) error {
			fmt.Printf(fb.tm(fb.messages["end"])+"\n", fb.username)
			if fn, ok := fb.Handlers["end"]; ok {
				fn();
			}
			return nil
		},
	})
	fb.pL = playerlist.New(fb.client)
	fb.chat = msg.New(fb.client, fb.player, fb.pL, msg.EventsHandler{
		SystemChat: func(c chat.Message, _ bool) error {
			fb.chatHandler(c)
			return nil
		},
		PlayerChatMessage: func(c chat.Message, _ bool) error {
			fb.chatHandler(c)
			return nil
		},
		DisguisedChat:    func(c chat.Message) error {
			fb.chatHandler(c)
			return nil
		},
	})
	return fb
}


func (fb *FragBot) Start() {
	err := fb.client.JoinServer("mc.hypixel.net:25565")
	if err != nil {
		panic(err)
	}
	for {
		if err = fb.client.HandleGame(); err == nil {
			panic("HandleGame never return nil")
		}

		if err2 := new(mc.PacketHandlerError); errors.As(err, err2) {
			if err := new(DisconnectErr); errors.As(err2, err) {
				fb.Log(fb.wh,"Disconnect, reason: ", err.Reason)
				return
			} else {
				fb.Log(fb.wh,err2)
			}
		} else {
			fb.Log(fb.wh,err)
		}
	}
}

type DisconnectErr struct {
	Reason chat.Message
}

func (d DisconnectErr) Error() string {
	return "disconnect: " + d.Reason.String()
}


func (fb *FragBot) joinNextParty() {
	fb.Current = ""
	fb.Queue = fb.Queue[1:]
	if len(fb.Queue) == 0 {
		return
	}

	username := fb.Queue[0]
	fb.chat.SendMessage("/p accept " + username)
	fb.Current = username
	fmt.Printf(fb.tm(fb.messages["joined"])+"\n", fb.username, username)
	
	go func() {
		time.Sleep(time.Second * 1)
		if _, ok := fb.config["message"]; ok {
			fb.chat.SendMessage("/pc " + fmt.Sprintf("%v",fb.config["message"]))
		}
		time.Sleep(time.Second * 6)
		if fb.Current != username {
			return
		}
		fb.Log(fb.wh,strings.ReplaceAll(fb.tm(fb.messages["disbanded"]),"%i", username))
		fb.chat.SendMessage("/p leave")
		fb.joinNextParty()
	}()
}

func (fb *FragBot) tm(message string) string {
	return strings.ReplaceAll(strings.ReplaceAll(message, "%s", fb.username), "%p", fmt.Sprint(len(fb.Queue)))
}

func indexOf(slice []string, str string) int {
    for i, s := range slice {
        if s == str {
            return i
        }
    }
    return -1
}


func (fb *FragBot) chatHandler(aaa chat.Message) error {
	msg := messageJoin(aaa);
	if msg == "You are AFK. Move around to return from AFK." {
		fb.Log(fb.wh,fb.tm(fb.messages["limbo"]))
		if fn, ok := fb.Handlers["limbo"]; ok {
			fn();
		}
		return nil
	}
	cur := fb.Current
	if cur == "" {
		cur = "-"
	}

	if strings.Contains(msg, cur + " entered The ") {
		fb.Log(fb.wh,strings.ReplaceAll(fb.tm(fb.messages["dungeons"]), "%i", fb.Current))
		fb.chat.SendMessage("/p leave")
		if fn, ok := fb.Handlers["limbo"]; ok {
			fn();
		}
		fb.joinNextParty()
		return nil
	}
	
	if strings.Contains(msg, "The party was disbanded because all invites expired and the party was empty") || strings.Contains(msg, cur + " has disbanded the party!") {
		if fn, ok := fb.Handlers["disbanded"]; ok {
			fn();
		}
		fb.Log(fb.wh,strings.ReplaceAll(fb.tm(fb.messages["disbanded"]), "%i", fb.Queue[0]))
		fb.joinNextParty()
		return nil
	}
	
	if strings.HasPrefix(msg, "The party invite from") {
		username := removeRankPrefix(strings.Split(strings.Split(msg, "invite from ")[1], " has expired")[0])
		index := indexOf(fb.Queue, username)
		if index > -1 {
			fb.Queue = append(fb.Queue[:index], fb.Queue[index+1:]...)
		}
	}
	

	if strings.Contains(msg, " has invited you to join their party!") {
		username := partyInviteGetUsername(msg)
		if contains(fb.bl,username) {
			return errors.New(username+" is in the blacklist")
		}
		if strings.Contains(username, " ") {
			return errors.New("Invalid username")
		}
		if strings.Index((strings.Join(fb.Queue," ")+" "),username+" ")>-1 {
			return errors.New("Username already in the queue")
		}
		fb.Queue = append(fb.Queue,username)
		if fn, ok := fb.Handlers["invite"]; ok {
			fn();
		}
		fb.Log(fb.wh,strings.ReplaceAll(fb.tm(fb.messages["invite"]),"%i", username))
		if len(fb.Queue) != 1 {
			return nil
		}
		fb.chat.SendMessage("/p accept " + username)

		fb.Current = username
		fb.Log(fb.wh,strings.ReplaceAll(fb.tm(fb.messages["joined"]),"%i", username))
		go func() {
			time.Sleep(time.Second * 1)
			if _, ok := fb.config["message"]; ok {
				fb.chat.SendMessage("/pc " + fmt.Sprintf("%v",fb.config["message"]))
			}
			time.Sleep(time.Second * 6)
			if fb.Current != username {
				return
			}
			fb.Log(fb.wh,strings.ReplaceAll(fb.tm(fb.messages["disbanded"]),"%i", username))
			fb.chat.SendMessage("/p leave")
			fb.joinNextParty()
		}()
	}
	return nil;
}

func (fb *FragBot) setMessage(msgType string, message string) {
	fb.messages[msgType] = message
}

func contains(s []string, searchterm string) bool {
    i := sort.SearchStrings(s, searchterm)
    return i < len(s) && s[i] == searchterm
}