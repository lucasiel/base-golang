package modules

import (
	"io/ioutil"
	"fmt"
	"bytes"
)

type Maps struct {
	Code int
	Author string
	Perma int
	Xml []byte
}

type Server struct {
	Clients map[string]*Client
	Maps map[int]*Maps
	Rooms map[string]*Rooms
	ReverseMap bool
}

func ServerNew() *Server {
	this := &Server {make(map[string]*Client),make(map[int]*Maps),make(map[string]*Rooms),false}
	return this
}

func (this *Server) LoadMapDatabase() {
	db2, err := ioutil.ReadFile("maps.dat")
	if err != nil {
		fmt.Println("[-] No file maps.dat")
		return
	}
	db := ByteArrayNew(bytes.NewBuffer(db2))

	mapCount := db.ReadInt()
	fmt.Println("[+] There are",mapCount,"maps in the database.")
	for i,_ := range make([]byte, mapCount) {
	    mapCode := db.ReadInt()
        mapAuthor := db.ReadUTF()
        mapPerma := db.ReadByte()
		length := db.ReadInt()
		xml := make([]byte, length)
		db.Bytes.Read(xml)
		this.Maps[int(i)] = &Maps{mapCode, mapAuthor, mapPerma, xml}
	}
}

func (this *Server) AddClientToRoom(client *Client, roomName string) {
	if val, ok := this.Rooms[roomName]; ok {
		if ok {
			val = val
			this.Rooms[roomName].AddClient(client)
			return
		}
	}
	room := Room(this,roomName)
	this.Rooms[roomName] = room
	room.AddClient(client)
}
