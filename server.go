package main
import (
	"fmt"
	"net"
	"strings"
	"time"
	. "server/modules"
)

var AddUUID int

func main() {
	//Log
	start := time.Now()
	server := ServerNew()
	AddUUID = 2
	Create_Server("0.0.0.0:11801",
		func() {
			server.LoadMapDatabase()
			fmt.Println("[+] Server loaded in",time.Since(start))
		},

		func(conn net.Conn) *Client {
			AddUUID = AddUUID + 1
			pid := ClientNew(conn,"",int(AddUUID),server)
			return pid
		},

		func(conn net.Conn, packet ByteArray, this *Client) {
			packet.ReadByte()
			C := packet.ReadByte()
			CC := packet.ReadByte()
			if C == 4 {
				if CC == 4 {
					packet.WriteCC(4,4)
					packet.WriteInt(this.PlayerID)
					packet.Writer.Write(packet.Bytes.Bytes())
					this.Room.SendPacketOthers(&packet,this.PlayerName)
					return
				}
				if CC == 5 {
					this.Alive = false
					this.Score++
					packet.WriteCC(1,1)
					packet.WriteUTF(fmt.Sprintf("\x08\x05\x01%d\x01%d",this.PlayerID,this.Score))
					this.Room.SendPacket(&packet)
					this.Room.CheckChangeMap()
					return
				}
				if CC == 9 {
					packet.WriteCC(4,9)
					packet.WriteInt(this.PlayerID)
					packet.WriteByte(packet.ReadByte())
					packet.WriteByte(0)
					this.Room.SendPacket(&packet)
					return
				}
			}
			if C == 5 {
				if CC == 19 {
					if this.HaveCheese { return }
					this.HaveCheese = true
					packet.WriteCC(144,6)
					packet.WriteInt(this.PlayerID)
					packet.WriteBoolean(true)
					this.Room.SendPacket(&packet)
					return
				}
				if CC == 18 {
					this.Room.MiceEnter += 1
					
					packet.WriteCC(8,6)
					packet.WriteByte(3).WriteInt(this.PlayerID).WriteShort(0).WriteByte(this.Room.MiceEnter).WriteShort(int(time.Since(this.Room.StartTime)) / 1000)
					if this.Room.MiceEnter < 4 {
						this.Score += 5 - this.Room.MiceEnter
					} else {
						this.Score++
					}
					this.Room.SendPacket(&packet)
					this.Alive = false
					this.Room.CheckChangeMap()
					return
				}
			}
			if C == 6 {
				if CC == 6 {
					message := packet.ReadUTF()
					packet.WriteCC(6,6)
					packet.WriteUTF(this.PlayerName)
					packet.WriteUTF(message)
					packet.WriteBoolean(true)
					this.Room.SendPacket(&packet)
					return
				}
			}
			if C == 8 {
				if CC == 5 {
					packet.WriteCC(8,5)
					packet.WriteInt(this.PlayerID)
					packet.WriteByte(packet.ReadByte())
					this.Room.SendPacket(&packet)
					return
				}
			}
			if C == 26 {
				if CC == 8 {
					username := packet.ReadUTF()
					packet.ReadUTF() //password not needed for now
					packet.ReadUTF() // ??
					this.PlayerName = username
					server.Clients[username] = this
					packet.WriteCC(26,2)
					packet.WriteInt(this.PlayerID)
					packet.WriteUTF(username)
					packet.WriteInt(600000)
					packet.WriteByte(6)
					packet.WriteInt(this.PlayerID)
					packet.WriteBoolean(true);
					packet.WriteByte(13)
					packet.WriteByte(-1)
					packet.WriteByte(13)
					packet.WriteByte(13)
					packet.WriteByte(5)
					packet.WriteByte(-1)
					packet.WriteByte(13)
					packet.WriteByte(15)
					packet.WriteByte(11)
					packet.WriteByte(5)
					packet.WriteByte(5)
					packet.WriteByte(5)
					packet.WriteByte(10)
					packet.WriteByte(10)
					packet.WriteBoolean(false);
					packet.WriteShort(255);
					packet.WriteShort(0);
					this.SendPacket(&packet)
					this.JoinRoom(packet.ReadUTF())
					return
				}
				if CC == 25 {
					packet.WriteCC(26,25)
				}
			}
			if C == 28 {
				if CC == 1 {
					version := packet.ReadShort()
					langue := packet.ReadUTF()
					ckey := packet.ReadUTF()
					if version != 616 {
						this.Transport.Close()
						return
					}
					if strings.Contains(ckey, "yAdByj") {
						packet.WriteCC(26,3)
						packet.WriteInt(0)
						packet.WriteUTF(langue)
						packet.WriteUTF(langue)
						packet.WriteInt(0)
						packet.WriteBoolean(false)
					} else {
						this.Transport.Close()
					}
				}
				if CC == 4 {
					eC := packet.ReadByte()
					eCC := packet.ReadByte()
					packet.ReadByte()
					packet.ReadByte()
					fmt.Println("[-] Error ",eC,eCC,packet.Bytes)
				}
				if CC == 17 {
					packet.WriteCC(20,4)
					packet.WriteShort(0)
				}
			}
			if C == 176 {
				if CC == 1 {
					packet.WriteCC(176,5)
					packet.WriteUTF("en")
					packet.WriteUTF("gb")
					packet.WriteBoolean(false);
					packet.WriteBoolean(true);
					packet.WriteUTF("")
				}
				if CC == 2 {
					packet.WriteCC(176,6)
					packet.WriteShort(2)
					packet.WriteUTF("en")
					packet.WriteUTF("English")
					packet.WriteUTF("gb")
					packet.WriteUTF("ro")
					packet.WriteUTF("Romania")
					packet.WriteUTF("ro")
				}
			}
				
			if packet.Writer.Len() >= 2 {
				this.SendPacket(&packet)
			}
		},

		func(conn net.Conn, this *Client) {
			if val, ok := server.Clients[this.PlayerName]; ok {
				if ok {
					delete(server.Clients,val.PlayerName)
				}
			}
		},
	)
}