package controllers

import (
	"fmt"
	"time"
	"util"

	"github.com/gorilla/websocket"
)

const (
	kMaxMessageSize = 512
	kPongWait       = 10 * time.Second
	kWriteWait      = 10 * time.Second
	kPingPeriod     = kPongWait / 2
)

type Client struct {
	room *Room // joined room
	conn *websocket.Conn
}

func NewClient(conn *websocket.Conn) *Client {
	newClient := &Client{conn: conn}
	go newClient.keepAlive()
	return newClient
}

func (c *Client) ReadMessage() {
	c.conn.SetReadLimit(kMaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(kPongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(kPongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				util.Log.Error(err.Error())
			}
			break
		}
		util.Log.Debug(string(message))
		c.room.Broadcast(string(message))
	}

	// 异常退出
	c.Exit()
}

func (c *Client) keepAlive() {
	ticker := time.NewTicker(kPingPeriod)
	defer ticker.Stop()
	defer c.Exit()
	for {
		select {
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(kWriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *Client) SendMessage(msgs []string) {
	if c.room == nil || c.conn == nil {
		return
	}

	c.conn.SetWriteDeadline(time.Now().Add(kWriteWait))
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		util.Log.Error(err.Error())
		return
	}

	for _, msg := range msgs {
		fmt.Fprintln(w, msg)
	}
	if err := w.Close(); err != nil {
		util.Log.Error(err.Error())
		return
	}
}

func (c *Client) Join(room *Room) {
	room.AddClient(c)
	c.room = room
}

func (c *Client) Exit() {
	if c.room == nil {
		return
	}

	c.room.DelClient(c)
	c.conn.Close()
}
