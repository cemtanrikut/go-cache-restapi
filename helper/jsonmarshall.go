package helper

import "encoding/json"

func JsonMarshall(resp *Resp) []byte {
	response, err := json.Marshal(&resp)
	if err != nil {
		panic(err)
	}
	return response
}
