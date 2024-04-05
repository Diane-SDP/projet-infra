var ws = new WebSocket("ws://82.67.64.82:80/ws");
var color = document.getElementById("jscolor")
var getcolor = document.getElementById("getcolor").value
var getpseudo = document.getElementById("getpseudo").value
var getcode = document.getElementById("getcode").value
ws.onopen = function(event) {
    console.log("WebSocket connected.");
};

ws.onerror = function(error) {
    console.error('WebSocket Error:', error);
};

function changeColor() {
    var button = document.getElementById('colorButton');
    var newColor = button.style.backgroundColor === 'green' ? 'red' : 'green';
    console.log("tient :",newColor)
    var message = newColor+"|"+ getpseudo + "|" + getcode
    ws.send(message)
    console.log("message send : ",message)
    button.style.backgroundColor = newColor;
}
ws.onmessage = function(event) {
    var button = document.getElementById('colorButton');
    console.log("recup :",event.data)
    if (event.data.split("|")[0] != "green" && event.data.split("|")[0] != "red") {
        var cookie = getCookie("uid")
        if (cookie == null) {
            var uid = event.data;
            document.cookie = "uid=" + uid + ";path=/";
            console.log("cookie créer avec l'uid : ",uid)
        }

    } else {
        console.log("ça c'est une couleur : ")
        color.innerText = (event.data.split("|")[0])
        button.style.backgroundColor = event.data.split("|")[0];
    }
    };
ws.onclose = function(event) {
    console.log("websocket closed")
}

function getCookie(name) {
    var dc = document.cookie;
    var prefix = name + "=";
    var begin = dc.indexOf("; " + prefix);
    if (begin == -1) {
        begin = dc.indexOf(prefix);
        if (begin != 0) return null;
    }
    else
    {
        begin += 2;
        var end = document.cookie.indexOf(";", begin);
        if (end == -1) {
        end = dc.length;
        }
    }
    // because unescape has been deprecated, replaced with decodeURI
    //return unescape(dc.substring(begin + prefix.length, end));
    return decodeURI(dc.substring(begin + prefix.length, end));
} 