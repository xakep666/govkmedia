package dialogboxes

import (
    "gopkg.in/qml.v1"
    "log"
    "govkmedia/requestwrapper"
    "strconv"
)

type groupSelectionDialog struct {
    engine *qml.Engine
    window *qml.Window
    ra      requestwrapper.RequestAccesser
    MyUid   int
}

type group struct {
    Gid int
    Name string
    Avatar string
    Subscribers int
}

func SelectGroupDialog(uid int, token string) int {
    gsd:=groupSelectionDialog{MyUid: uid}
    gsd.engine=qml.NewEngine()
    component,err:=gsd.engine.LoadFile("qrc:///qml/groupselectdialog.qml")
    if err!=nil {
        log.Panicln(err)
    }
    gsd.window=component.CreateWindow(nil)
    gsd.engine.Context().SetVar("engine",&gsd)
    gsd.ra=requestwrapper.RequestAccesser{Token: token}
    gsd.registerQMLTypes()
    gsd.LoadUserGroups(uid)
    gid:=make(chan int) 
    gsd.window.Root().On("selectedgidchanged",func(){
        gid <-gsd.window.Root().Int("selectedgid")
    })
    gsd.window.Show()
    return <-gid
}

func (gsd *groupSelectionDialog) registerQMLTypes() {
    types:=[]qml.TypeSpec{
        {
            Init: func(m *group, obj qml.Object) {},
        },
    }
    qml.RegisterTypes("GoExtensions",1,0,types)
}

func (gsd *groupSelectionDialog) LoadUserGroups(uid int) {
    resp,err:=gsd.ra.MakeRequest("groups.get",map[string]string{
        "user_id":strconv.Itoa(uid),
        "extended": "1",
        "fields": "members_count",
    })
    if err!=nil {
        ShowErrorDialog(err.Error())
        return
    }
    mwroot:=gsd.window.Root()
    mwroot.ObjectByName("grouplistview").Set("subscribersvisible",true)
    gsd.loadGroups(resp)
}

func (gsd *groupSelectionDialog) SearchGroups(query string) {
    resp,err:=gsd.ra.MakeRequest("groups.search",map[string]string{
        "q":query,
    })
    if err!=nil {
        ShowErrorDialog(err.Error())
        return
    }
    mwroot:=gsd.window.Root()
    mwroot.ObjectByName("grouplistview").Set("subscribersvisible",false)
    gsd.loadGroups(resp)
}

func (gsd *groupSelectionDialog) loadGroups(resp map[string]interface{}) {
    mwroot:=gsd.window.Root()
    model:=mwroot.ObjectByName("grouplist")
    model.Call("clear")
    resp_groups:=resp["response"].(map[string]interface{})["items"].([]interface{})
    for _,v:=range resp_groups {
        resp_group:=v.(map[string]interface{})
        subsc,_:=resp_group["members_count"].(float64)
        append_group:=group{
            Gid: int(resp_group["id"].(float64)),
            Name: resp_group["name"].(string),
            Avatar: resp_group["photo_100"].(string),
            Subscribers: int(subsc),
        }
        model.Call("appendStruct",&append_group)
    }
}