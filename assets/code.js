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
    if (event.data != "green" && event.data != "red") {
        var uid = event.data;
        document.cookie = "uid=" + uid + ";path=/";
        console.log("cookie créer avec l'uid : ",uid)
    } else {
        color.innerText = (event.data)
        button.style.backgroundColor = event.data;
    }
    };
ws.onclose = function(event) {
    console.log("websocket closed")
}