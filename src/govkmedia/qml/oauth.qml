import QtQuick 2.0
import QtQuick.Controls 1.3
import QtWebKit 3.0
import QtQuick.Window 2.0
Window{
    id: oauthwindow
    width: 500
    height: 450
    title: qsTr("Авторизация")
    WebView {
        id: oauthww
        objectName: qsTr("oauthww")
        width: parent.width
        height: parent.height
        onLoadingChanged: {
            var stringurl = loadRequest.url.toString()
            if (stringurl.substring(0,31)==="https://oauth.vk.com/blank.html") {
                oauthwindow.close()
                appEngine.checkAuth(loadRequest.url)
            }
        }
    }
}
