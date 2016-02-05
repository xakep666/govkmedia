package dialogboxes

import (
    "gopkg.in/qml.v1"
    "log"
)

//go:generate genqrc qml icons

func ShowErrorDialog(text string) {
    engine := qml.NewEngine()
    component, err := engine.LoadFile("qrc:///qml/errordialog.qml")
    if err != nil {
        log.Panicln(err)
    }
    window := component.CreateWindow(nil)
    window.Root().ObjectByName("errortext").Set("text",text)
    window.Show()
}
