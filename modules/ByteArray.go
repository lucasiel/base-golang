package modules
import (
	"bytes"
	"encoding/binary"
)

type ByteArray struct {  
    Bytes *bytes.Buffer
	Writer *bytes.Buffer
}
func ByteArrayNew(byter *bytes.Buffer) *ByteArray {
	this := &ByteArray {byter, bytes.NewBuffer(make([]byte,0))}
	return this
}

func (this *ByteArray) ReadByte() int {  
	data := make([]byte, 1)
	this.Bytes.Read(data)
    return int(data[0])
}

func (this *ByteArray) ReadShort() int {  
	data := make([]byte, 2)
	this.Bytes.Read(data)
    return int(binary.BigEndian.Uint16(data))
}

func (this *ByteArray) ReadInt() int {  
	data := make([]byte, 4)
	this.Bytes.Read(data)
    return int(binary.BigEndian.Uint32(data))
}

func (this *ByteArray) ReadUTF() string {  
	length := this.ReadShort()
	data := make([]byte, length)
	this.Bytes.Read(data)
    return string(data)
}

func (this *ByteArray) WriteByte(arg int) *ByteArray { 
	this.Writer.Write([]byte{byte(arg)})
	return this
}

func (this *ByteArray) ReadBoolean() bool { 
	d := this.ReadByte()
	if d != 0 {
		return true
	}
	return false
}


func (this *ByteArray) WriteBoolean(arg bool) *ByteArray { 
	if arg {
		this.WriteByte(1)
	} else {
		this.WriteByte(0)
	}
	return this
}

func (this *ByteArray) WriteShort(arg int) *ByteArray { 
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, uint16(arg))
	this.Writer.Write(bs)
	return this
}

func (this *ByteArray) WriteInt(arg int) *ByteArray { 
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(arg))
	this.Writer.Write(bs)
	return this
}

func (this *ByteArray) WriteUTF(arg string) *ByteArray { 
	this.WriteShort(len(arg))
	this.Writer.Write([]byte(arg))
	return this
}

func (this *ByteArray) WriteCC(C int,CC int) *ByteArray { 
	this.WriteByte(C)
	this.WriteByte(CC)
	return this
}