package httpclient

import (
	"encoding/json"
)

type M map[string]interface{}

func (m M) MarshalJSON() ([]byte, error) {
	type tmpType M
	var tmp = tmpType(m)
	return json.Marshal(&tmp)
}
