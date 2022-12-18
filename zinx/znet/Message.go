package znet

type Message struct {
	// 消息ID
	Id uint32
	// 消息长度
	DateLen uint32
	// 消息内容
	Data []byte
}

func NewMessage(id uint32, data []byte) *Message{
	return &Message{
		Id: id,
		DateLen: uint32(len(data)),
		Data: data,
	}
}

// getter
func (msg *Message) GetMegId() uint32 {
	return msg.Id
}
func (msg *Message) GetMegDataLen() uint32 {
	return msg.DateLen
}
func (msg *Message) GetMegData() []byte {
	return msg.Data
}

// setter
func (msg *Message) SetMegId(id uint32) {
	msg.Id = id
}
func (msg *Message) SetMegDataLen(data_len uint32) {
	msg.DateLen = data_len
}
func (msg *Message) SetMegData(data []byte) {
	msg.Data = data
}
