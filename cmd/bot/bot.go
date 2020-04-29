package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Tnze/go-mc/cmd/bot/authUser"
	"github.com/jinzhu/configor"
	"github.com/mattn/go-colorable"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/chat"
	_ "github.com/Tnze/go-mc/data/lang/en-us"
)

const timeout = 45

var Config = struct {
	APPName string `default:"app name"`

	Info struct {
		Server string  `required:"true"`
		Yaw    float32 `required:"true"`
		Pitch  float32 `required:"true"`
	}

	Account struct {
		Username string `required:"true"`
		Password string `required:"true"`
	}
}{}

var (
	c     *bot.Client
	watch chan time.Time
)

func main() {
	log.SetOutput(colorable.NewColorableStdout())
	c = bot.NewClient()

	err := configor.Load(&Config, "config.json")

	if err != nil {
		fmt.Println(err)
		return
	}

	userAuth := authUser.GetToken(Config.Account.Username, Config.Account.Password)

	if userAuth != nil {
		// c.Auth.Name = userAuth.Name
		c.Auth.UUID = userAuth.ID
		c.Name = userAuth.Name
		c.AsTk = userAuth.Tokens.AccessToken
		// c.Auth.AsTk = userAuth.Tokens.AccessToken
	}

	//Login
	err = c.JoinServer(Config.Info.Server, 25565)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Login success")

	//Register event handlers
	c.Events.GameStart = onGameStart
	c.Events.ChatMsg = onChatMsg
	c.Events.Disconnect = onDisconnect
	c.Events.Die = onDeath
	c.Events.

		//JoinGame
		err = c.HandleGame()
	if err != nil {
		log.Fatal(err)
	}
}

func onDeath() error {
	log.Println("!!!!!!!!!!! Died and Respawned")
	c.Respawn()
	return nil
}

func onGameStart() error {
	log.Println("Game start")

	watch = make(chan time.Time)
	go watchDog()

	return c.UseItem(0)
}

// func onSound(name string, category int, x, y, z float64, volume, pitch float32) error {
// 	if name == "entity.fishing_bobber.splash" {
// 		if err := c.UseItem(0); err != nil { //retrieve
// 			return err
// 		}
// 		log.Println("gra~")
// 		time.Sleep(time.Millisecond * 300)
// 		if err := c.UseItem(0); err != nil { //throw
// 			return err
// 		}
// 		watch <- time.Now()
// 	}
// 	return nil
// }

func onChatMsg(c chat.Message, pos byte) error {
	log.Println("Chat:", c)
	return nil
}

func onDisconnect(c chat.Message) error {
	log.Println("Disconnect:", c)

	return nil
}

func watchDog() {
	to := time.NewTimer(time.Second * timeout)
	for {
		select {
		case <-to.C:
			log.Println("rethrow")
			if err := c.UseItem(0); err != nil {
				panic(err)
			}
		}
		to.Reset(time.Second * timeout)
	}
}
