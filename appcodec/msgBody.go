package appcodec

import (
	"encoding/json"
	"gateway/gw"
)

type Body struct {
	BType     string `json:"bt"`
	SessionId string `json:"sid"`

	SvrName    string `json:"sn"`
	Resource   string `json:"rs"`
	Action     string `json:"at"`
	Content    string `json:"ct"`
	Attachment string `json:"atta"`
}

func (b *Body) ToString() string {
	bs, _ := json.Marshal(b)
	return gw.Byte2str(bs)
}

func (b *Body) BytesTo(bs []byte) {
	gw.Logger.Debugf("srcText:%s", gw.Byte2str(bs))
	err := json.Unmarshal(bs, b)
	if err != nil {
		gw.Logger.Errorf("Body parse error:%s\r\n%s", err.Error(), bs)
	}
}

func (b *Body) ToBytes() []byte {
	bs, err := json.Marshal(b)
	if err != nil {
		gw.Logger.Errorf("Body parse error:%s", err.Error())
		return nil
	}
	return bs
}
