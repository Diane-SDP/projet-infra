var ws = new WebSocket("ws://82.67.64.82:80/ws");
ws.onopen = function(event) {
    console.log("WebSocket connected.");
};

ws.onerror = function(error) {
    console.error('WebSocket Error:', error);
};

ws.onmessage = function(event) {
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