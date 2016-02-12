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

var VkAudioGenres map[float64]string = map[float64]string{
    1:"Rock",
    2:"Pop",
    3:"Rap & Hip-Hop",
    4:"Easy Listening",
    5:"Dance & House",
    6:"Instrumental",
    7:"Metal",
    21:"Alternative",
    8:"Dubstep",
    1001:"Jazz & Blues",
    10:"Drum & Bass",
    11:"Trance",
    12:"Chanson",
    13:"Ethnic",
    14:"Acoustic & Vocal",
    15:"Reggae",
    16:"Classical",
    17:"Indie Pop",
    19:"Speech",
    22:"Electropop & Disco",
    18:"Other",
}
