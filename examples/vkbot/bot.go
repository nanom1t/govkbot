package main

import (
	"fmt"
	"github.com/nikepan/govkbot"
	"log"
	"net/url"
	"strconv"
)

type configuration struct {
	VKToken      string
	ServiceURL   string
	ServiceToken string
	AdminID      int
}

func getMeMessage(uid int) (reply string) {
	me := govkbot.API.Me()
	return fmt.Sprintf("You: %+v %+v", me.FirstName, me.LastName)
}

func anyHandler(m *govkbot.Message) (reply string) {
	notifyAdmin(fmt.Sprintf("Command %+v by user vk.com/id%+v in chat %+v", m.Body, m.UserID, m.Title))
	return reply
}

func meHandler(m *govkbot.Message) (reply string) {
	return getMeMessage(m.UserID)
}

func helpHandler(m *govkbot.Message) (reply string) {
	return availableCommands
}

func inviteHandler(m *govkbot.Message) (reply string) {
	log.Printf("invite: %+v %+v %+v\n", m.ActionMID, govkbot.API.Uid, m.ActionMID == govkbot.API.Uid)
	if m.ActionMID == govkbot.API.Uid {
		go m.MarkAsRead()
		notifyAdmin(fmt.Sprintf("I'm invited to chat %+v )", m.Title))
		reply = replyGreet()
	} else {
		log.Printf("greet: %+v %+v\n", m.ActionMID, m)
		reply = greetUser(m.ActionMID)
	}
	return reply
}

func kickHandler(m *govkbot.Message) (reply string) {
	if m.ActionMID == govkbot.API.Uid {
		go m.MarkAsRead()
		notifyAdmin(fmt.Sprintf("I'm kicked from chat %+v (", m.Title))
	}
	return reply
}

func greetUser(uid int) (reply string) {
	u, ok := govkbot.API.User(uid)
	if ok {
		reply = fmt.Sprintf("Hello, %+v %+v", u.FirstName, u.LastName)
	}
	return reply
}

func replyGreet() (reply string) {
	reply = "Hi all. I'am bot\n" + availableCommands
	return reply
}

func addFriendHandler(m *govkbot.Message) (reply string) {
	log.Printf("friend %+v added\n", m.UserID)
	notifyAdmin(fmt.Sprintf("user vk.com/id%+v add me to friends", m.UserID))
	return reply
}

func deleteFriendHandler(m *govkbot.Message) (reply string) {
	log.Printf("friend %+v deleted\n", m.UserID)
	notifyAdmin(fmt.Sprintf("user vk.com/id%+v delete me from friends", m.UserID))
	return reply
}

func notifyAdmin(msg string) {
	if config.AdminID != 0 {
		p := url.Values{}
		p.Add("user_id", strconv.Itoa(config.AdminID))
		p.Add("message", msg)
		_ = govkbot.API.Call("messages.send", p)
	}
}
