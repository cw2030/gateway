package appcodec

import (
	"encoding/json"
	"fmt"
)

type Body struct {
	BType     string `json:"bt"`
	SessionId string `json:"sid"`
	SvrType   string `json:"st"`

	SvrName    string `json:"sn"`
	Resource   string `json:"rs"`
	Action     string `json:"act"`
	Content    string `json:"content"`
	Attachment string `json:"atta"`
}

func (b *Body) ToString() string {
	bs, _ := json.Marshal(b)
	return string(bs)
}

func (b *Body) bytesTo(bs []byte) {
	json.Unmarshal(bs, b)
	fmt.Println(b)
}

func (b *Body) ToBytes() []byte {
	bs, err := json.Marshal(b)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return bs
}
