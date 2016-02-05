package dialogboxes

import (
    "gopkg.in/qml.v1"
    "log"
    "strings"
    "os/user"
)

func SelectFolderDialog() string{
    engine := qml.NewEngine()
    component, err := engine.LoadFile("qrc:///qml/filedialog.qml")
    if err != nil {
        log.Panicln(err)
    }
    window := component.CreateWindow(nil)
    curuser,err:=user.Current()
    if err!=nil {
        log.Panicln(err) //system error
    }
    fd:=window.Root().ObjectByName("filedialog")
    fd.Set("folder",curuser.HomeDir)
    window.Root().Call("setupForFolder")
    path:=make(chan string)
    window.Root().On("gotpath",func () {
        if window.Root().Bool("cancelled") {
            path <- ""
        } else {
            path <- fd.String("fileUrl")
        }
    })
    window.Show()
    return strings.Replace(<-path,"file://","",1)
}
