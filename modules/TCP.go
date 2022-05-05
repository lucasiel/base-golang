package modules
import ( 
	"net"
	"bytes"
	"strings"
	"fmt"
)

type ConnectNew func (conn net.Conn) *Client
type ReceiveBuffer func (conn net.Conn, sBuffer ByteArray, this *Client)
type ServerInit func ()
type ConnectClose func (conn net.Conn, this *Client)
var conn net.Conn

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

		client := funCN(conn)

		go tcpHandler(conn, funRB, funCC, client)
	}
}

func tcpHandler(conn net.Conn, funRB ReceiveBuffer, funCC ConnectClose, client *Client) {
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
				bufe3 := *ByteArrayNew(bufe2)

				funRB(conn, bufe3, client)
				contentLen = 0
				offset = 0
			}
			break
		} 
	}
}