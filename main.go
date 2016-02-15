package main

import (
        "gopkg.in/qml.v1"
        "strconv"
        "net/url"
        "log"
        "strings"
        "govkmedia/requestwrapper"
)

const APP_ID int = 5255931
const APP_SCOPE string = "audio,video,photos,offline"

//go:generate genqrc qml

func main() {
    err:=qml.Run(run)
    if err!=nil {
        log.Panicln(err.Error())
    }
}

func run() error {
    registerQMLTypes()
	engine := qml.NewEngine()
        component, err := engine.LoadFile("qrc:///qml/main.qml")
	if err != nil {
                return err
	}

        context := engine.Context()
        window := component.CreateWindow(nil)
        context.SetVar("appEngine",&AppEngine{QMLEngine: engine,MainWindow: window})
	window.Show()
	window.Wait()
        return nil
}

type AppEngine struct {
    requestwrapper.RequestAccesser
    QMLEngine *qml.Engine
    MainWindow *qml.Window
}

func (ae *AppEngine) ShowOauth() {
    component, err := ae.QMLEngine.LoadFile("qrc:///qml/oauth.qml")
    if err != nil {
        log.Panicln(err.Error())
    }
    window:=component.CreateWindow(nil)
    url:="https://oauth.vk.com/authorize?client_id="+strconv.Itoa(APP_ID)+
         "&redirect_uri=https://oauth.vk.com/blank.html"+
         "&display=mobile"+
         "&response_type=token"+
         "&scope="+APP_SCOPE
    webview:=window.Root().ObjectByName("oauthww")
    webview.Set("url",url)
    window.Show()
}

func (ae *AppEngine) CheckAuth(inurl string) {
    inurl=strings.Replace(inurl,"#","?",1) //hack to get parser working
    authUrl,err:=url.Parse(inurl)
    if err!=nil {
        log.Println(err.Error())
        return
    }
    authQuery:=authUrl.Query()
    if authQuery.Get("error")!="" {
        ae.toggleWindow(false)
        mwroot:=ae.MainWindow.Root()
        txt:=mwroot.ObjectByName("noauthtxt")
        txt.Set("text","Ошибка авторизаци: "+authQuery.Get("error_description"))
        return
    }
    //auth successful, rerender window
    ae.Token=authQuery.Get("access_token")
    ae.UserId,_=strconv.Atoi(authQuery.Get("user_id"))
    go ae.loadAvatar()
    go ae.loadName()
    go ae.loadAudios(ae.UserId)
    ae.toggleWindow(true)
}

func (ae *AppEngine) toggleWindow(authDone bool) {
    root:=ae.MainWindow.Root()
    noauthobjects:=[...]string{"noauthtxt","authbtn"}
    for _,v:=range noauthobjects {
        object:=root.ObjectByName(v)
        object.Set("visible",!authDone)
    }
    authobjects:=[...]string{"tabs","avatar","name","open","dlallbtn","dlselbtn"}
    for _,v:=range authobjects {
        object:=root.ObjectByName(v)
        object.Set("visible",authDone)
    }
}

func registerQMLTypes () {
    types:=[]qml.TypeSpec{
        {
            Init: func(m *MusicItem, obj qml.Object) {},
        },
    }
    qml.RegisterTypes("GoExtensions",1,0,types)
}
