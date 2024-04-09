package rest

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
)

func parseRsp(r *resty.Response, t interface{}) error {
	if !r.IsSuccess() {
		return errors.New(r.String())
	}
	return json.Unmarshal(r.Body(), t)
}
