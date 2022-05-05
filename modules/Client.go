package modules

import (
	"net"
	"bytes"
)

type Client struct {  
	Transport net.Conn
	PlayerName string
	PlayerID int
	Server *Server
	Room *Rooms
	Alive bool
	HaveCheese bool
	Score int
}

func ClientNew(Conn net.Conn, name string, pid int, Server *Server) *Client {
	this := &Client {Conn, name, pid, Server, &Rooms{},false,false,0}
	return this
}

func (this *Client) SendMap() {
	p := ByteArrayNew(bytes.NewBuffer(make([]byte,0)))
	p.WriteCC(5,2)
	p.WriteInt(this.Room.Map.Code)
	p.WriteShort(len(this.Room.Clients))
	p.WriteByte(this.Room.RoundId)
	p.WriteInt(len(this.Room.Map.Xml))
	p.Writer.Write(this.Room.Map.Xml)
	p.WriteUTF(this.Room.Map.Author)
	p.WriteByte(this.Room.Map.Perma)
	p.WriteBoolean(this.Server.ReverseMap)
	this.SendPacket(p)
}

func (this *Client) PlayerData(add bool) *ByteArray {
	p := ByteArrayNew(bytes.NewBuffer(make([]byte,0)))
	p.WriteUTF(this.PlayerName)
	p.WriteInt(this.PlayerID)
	p.WriteBoolean(false)
	p.WriteBoolean(!this.Alive)
	p.WriteShort(this.Score)
	p.WriteBoolean(this.HaveCheese)
	p.WriteShort(0)
	p.WriteByte(0)
	p.WriteByte(0)
	p.WriteUTF("")
	p.WriteUTF("1;0,0,0,0,0,0,0,0,0,0")
	p.WriteBoolean(false)
	p.WriteInt(7886906)
	p.WriteInt(9820630)
	p.WriteInt(0)
	p.WriteInt(-1)
	if (add) {
		p.WriteBoolean(false)
		p.WriteBoolean(true)
	}
	return p
}
func (this *Client) SendPlayerList() {
	p := ByteArrayNew(bytes.NewBuffer(make([]byte,0)))
	p.WriteCC(144,1)
	p.WriteShort(len(this.Room.Clients))
	for _, client := range this.Room.Clients {
		p.Writer.Write(client.PlayerData(false).Writer.Bytes())
	}
	this.SendPacket(p)
}

func (this *Client) SendPacket(packet *ByteArray) {  
	temp := bytes.NewBuffer(make([]byte,0))
	length := packet.Writer.Len()
	calc1 := length >> 7
	for calc1 != 0 {
		temp.Write([]byte{byte(((length & 127) | 128))})
		length = calc1
		calc1 = calc1 >> 7
	}
	temp.Write([]byte{byte(length & 127)})
	temp.Write(packet.Writer.Bytes())
    this.Transport.Write(temp.Bytes())
}

func (this *Client) JoinRoom(roomName string) {  
	this.Server.AddClientToRoom(this, roomName)
	this.SendPacket(ByteArrayNew(bytes.NewBuffer(make([]byte,0))).WriteCC(7,1).WriteByte(0))
	this.SendPacket(ByteArrayNew(bytes.NewBuffer(make([]byte,0))).WriteCC(7,30).WriteByte(9))
	this.SendPacket(ByteArrayNew(bytes.NewBuffer(make([]byte,0))).WriteCC(5,21).WriteBoolean(true).WriteUTF(roomName).WriteUTF("gb"))
}