package main
import (
	"fmt"
	"net"
	"bytes"
	"strings"
	"encoding/binary"
	"time"
)

var AddUUID uint64
type ConnectNew func (conn net.Conn) uint64
type ReceiveBuffer func (conn net.Conn, sBuffer ByteArray, UUID uint64)
type ServerInit func ()
type ConnectClose func (conn net.Conn, UUID uint64)
var conn net.Conn

type ByteArray struct {  
    bytes *bytes.Buffer
	conn net.Conn
	writer *bytes.Buffer
}

func New(bytes *bytes.Buffer, conn net.Conn, writer *bytes.Buffer) ByteArray {
	this := ByteArray {bytes, conn, writer}
	return this
}

func (this ByteArray) readByte() int {  
	data := make([]byte, 1)
	this.bytes.Read(data)
    return int(data[0])
}

func (this ByteArray) readShort() int {  
	data := make([]byte, 2)
	this.bytes.Read(data)
    return int(binary.LittleEndian.Uint16(reverseArray(data)))
}

func (this ByteArray) readInt() int {  
	data := make([]byte, 4)
	this.bytes.Read(data)
    return int(binary.LittleEndian.Uint32(reverseArray(data)))
}

func (this ByteArray) readUTF() string {  
	length := this.readShort()
	data := make([]byte, length)
	this.bytes.Read(data)
    return string(data)
}

func (this ByteArray) writeByte(arg int) { 
	this.writer.Write([]byte{byte(arg)})
}

func (this ByteArray) readBoolean() bool { 
	d := this.readByte()
	if d != 0 {
		return true
	}
	return false
}


func (this ByteArray) writeBoolean(arg bool) { 
	if arg {
		this.writeByte(1)
	} else {
		this.writeByte(0)
	}
}

func (this ByteArray) writeShort(arg int) { 
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, uint16(arg))
	this.writer.Write(reverseArray(bs))
}

func (this ByteArray) writeInt(arg int) { 
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(arg))
	this.writer.Write(reverseArray(bs))
}

func (this ByteArray) writeUTF(arg string) { 
	this.writeShort(len(arg))
	this.writer.Write([]byte(arg))
}

func (this ByteArray) writeCC(C int,CC int) { 
	this.writeByte(C)
	this.writeByte(CC)
}

func (this ByteArray) sendPacket() {  
	temp := bytes.NewBuffer(make([]byte,0))
	length := this.writer.Len()
	calc1 := length >> 7
	for calc1 != 0 {
		temp.Write([]byte{byte(length)})
		length = calc1
		calc1 = calc1 >> 7
	}
	temp.Write([]byte{byte(length & 127)})
	temp.Write(this.writer.Bytes())
    this.conn.Write(temp.Bytes())
	this.writer.Reset()
}

func Create_Server(webpath string, init ServerInit, funCN ConnectNew, funRB ReceiveBuffer, funCC ConnectClose) {
	listener, err := net.Listen("tcp", webpath)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}

	init()
	defer listener.Close()

	for {
		conn, err = listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			return
		}

		UUID := funCN(conn)

		go tcpHandler(conn, funRB, funCC, UUID)
	}
}

func tcpHandler(conn net.Conn, funRB ReceiveBuffer, funCC ConnectClose, UUID uint64) {
	cache := make([]byte, 1048576)
	buf := bytes.NewBuffer(make([]byte,0, 1048576))

	var contentLen int
	var code int
	var offset int

	for {
		size, err := conn.Read(cache)
		if err == nil {
			buf.Write(cache[:size])
		}

		if strings.Contains(string(cache[:size]),"<policy-file-request/>") {
			conn.Write([]byte("<cross-domain-policy><allow-access-from domain=\"*\" to-ports=\"*\"/></cross-domain-policy>"))
			conn.Close()
			return
		}

 		for {
			if buf.Len() == 0 {
				break
			}
			if contentLen == 0 {
				code = 128
				data := make([]byte, 1)
				for (code & 128) != 0 {
					_, err = buf.Read(data)
					code = int(data[0]) & 255
					contentLen = contentLen | ((code & 127) << (offset * 7))
					offset++;
				}
				contentLen++;
			}
			if buf.Len() >= int(contentLen) {
				data2 := make([]byte, contentLen)
				_, err = buf.Read(data2)
				bufe2 := bytes.NewBuffer(data2)
				bufe3 := New(bufe2, conn, bytes.NewBuffer(make([]byte, 0)))

				go funRB(conn, bufe3, UUID)
				contentLen = 0
				offset = 0
			}
			break
		} 
	}
}
//var connections []net.Conn
func reverseArray(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
        s[i], s[j] = s[j], s[i]
    }
	return s
}
func main() {
	//Log
	start := time.Now()
	Create_Server("0.0.0.0:11801",
		func() {
			fmt.Println("Server loaded in",time.Since(start))
		},

		func(conn net.Conn) uint64 {
			AddUUID++
			//append(connections,conn)
			return AddUUID
		},

		func(conn net.Conn, packet ByteArray, UUID uint64) {
			packet.readByte()
			C := packet.readByte()
			CC := packet.readByte()
			fmt.Println(C,CC)
			if C == 28 {
				if CC == 1 {
					version := packet.readShort()
					langue := packet.readUTF()
					ckey := packet.readUTF()
					if version != 689 {
						packet.conn.Close()
						return
					}
					if strings.Contains(ckey, "CnDnJ") {
						packet.writeCC(26,3)
						packet.writeInt(0)
						packet.writeUTF(langue)
						packet.writeUTF(langue)
						packet.writeInt(0)
						packet.writeBoolean(false)
					} else {
						packet.conn.Close()
					}
				}
				if CC == 4 {
					eC := packet.readByte()
					eCC := packet.readByte()
					packet.readByte()
					packet.readByte()
					fmt.Println("[-] Error ",eC,eCC,packet.bytes)
				}
				if CC == 17 {
					packet.writeCC(20,4)
					packet.writeShort(0)
				}
			}
			if C == 176 {
				if CC == 1 {
					packet.writeCC(176,5)
					packet.writeUTF("en")
					packet.writeUTF("gb")
					packet.writeBoolean(false);
					packet.writeBoolean(true);
					packet.writeUTF("")
				}
				if CC == 2 {
					packet.writeCC(176,6)
					packet.writeShort(1)
					packet.writeUTF("en")
					packet.writeUTF("English")
					packet.writeUTF("gb")
				}
			}
				
			if packet.writer.Len() >= 2 {
				packet.sendPacket()
			}
		},

		func(conn net.Conn, UUID uint64) {
			//remove(connections,conn)
		},
	)
}