package modules
import (
	"bytes"
	"math/rand"
	"time"
	"fmt"
)
type Rooms struct {
	Server *Server
	Name string
	Clients map[string]*Client
	Map Maps
	playerList *ByteArray
	RoundId int
	MiceEnter int
	StartTime time.Time
}

func Room(server *Server, name string) *Rooms {
	this := &Rooms {server,name,make(map[string]*Client),*server.Maps[1],ByteArrayNew(bytes.NewBuffer(make([]byte,0))),1,0,time.Now()}
	return this
}

func (this *Rooms) newRound() {
	this.SendPacket(ByteArrayNew(bytes.NewBuffer(make([]byte,0))).WriteCC(5,10).WriteBoolean(true))
	this.RoundId = (this.RoundId + 1) % 100
	this.StartTime = time.Now()
	this.MiceEnter = 0
	rand.Seed(time.Now().Unix())
	this.Map = *this.Server.Maps[rand.Intn(len(this.Server.Maps))]
	for _, client := range this.Clients {
		client.Alive = true
		client.HaveCheese = false
		client.SendMap()
		client.SendPlayerList()
		client.SendPacket(ByteArrayNew(bytes.NewBuffer(make([]byte,0))).WriteCC(8,11).WriteInt(0).WriteInt(0).WriteInt(0).WriteInt(0).WriteShort(0))
		client.SendPacket(ByteArrayNew(bytes.NewBuffer(make([]byte,0))).WriteCC(1,1).WriteUTF(fmt.Sprintf("\x08\x15\x011\x01")))
		client.SendPacket(ByteArrayNew(bytes.NewBuffer(make([]byte,0))).WriteCC(5,22).WriteShort(63))
	}
	go func() {
		time.Sleep(3 * time.Second)
		this.SendPacket(ByteArrayNew(bytes.NewBuffer(make([]byte,0))).WriteCC(5,10).WriteBoolean(false))
	}()

}

func (this *Rooms) CheckChangeMap() {
	for _, client := range this.Clients {
		if (client.Alive) {
			return
		}
	}
	this.newRound()
}

func (this *Rooms) SendPacket(packet *ByteArray) {
	for _, client := range this.Clients {
		client.SendPacket(packet)
	}

}

func (this *Rooms) SendPacketOthers(packet *ByteArray, name string) {
	for _, client := range this.Clients {
		if client.PlayerName != name {
			client.SendPacket(packet)
		}
	}

}

func (this *Rooms) AddClient(client *Client) {
	this.Clients[client.PlayerName] = client
	client.Room = this
	
	if len(this.Clients) == 1 {
		client.Alive = true
		this.newRound()
		return
	}
	
	p := ByteArrayNew(bytes.NewBuffer(make([]byte,0))).WriteCC(144,2)
	p.Writer.Write(client.PlayerData(true).Writer.Bytes())
	this.SendPacketOthers(p,client.PlayerName)
	client.SendMap()
	client.SendPlayerList()
	client.SendPacket(ByteArrayNew(bytes.NewBuffer(make([]byte,0))).WriteCC(8,11).WriteInt(0).WriteInt(0).WriteInt(0).WriteInt(0).WriteShort(0))
}