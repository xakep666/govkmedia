package main
import (
    "log"
)

func (ae *AppEngine) showErrorDialog(text string) {
    component, err := ae.QMLEngine.LoadFile("qrc:///qml/errordialog.qml")
    if err != nil {
        log.Panicln(err.Error())
    }
    window:=component.CreateWindow(nil)
    window.Root().ObjectByName("errortext").Set("text",text)
    window.Show()
}
