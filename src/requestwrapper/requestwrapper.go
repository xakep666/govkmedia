package requestwrapper

import (
    "fmt"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "errors"
)

type RequestAccesser struct {
    Token string
    UserId int
}

const REQUEST_STRING = "https://api.vk.com/method/%s?%saccess_token=%s&v=5.44"

func (a *RequestAccesser) MakeRequest(method string,params map[string]string) (resp map[string]interface{},e error) {
    var parm_str string
    for k,v:=range params {
        parm_str+=fmt.Sprintf("%s=%s&",k,v)
    }
    req_url:=fmt.Sprintf(REQUEST_STRING,method,parm_str,a.Token)
    httpresp,err:=http.Get(req_url)
    if err!=nil {
        return nil,err
    }
    defer httpresp.Body.Close()
    contents,err:=ioutil.ReadAll(httpresp.Body)
    if err!=nil {
        return nil,err
    }
    err=json.Unmarshal(contents,&resp)
    if err!=nil {
        return nil,err
    }
    resperr:=resp["error"]
    if resperr!=nil {
        errtext:=resperr.(map[string]interface{})["error_msg"].(string)
        return nil,errors.New(errtext)
    }
    return resp,nil
}
