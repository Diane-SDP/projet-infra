var ws = new WebSocket("ws://82.67.64.82:80/ws");
var color = document.getElementById("jscolor")
var getcolor = document.getElementById("getcolor")
var code = document.location.href.split('/')[4]



ws.onopen = function(event) {
    console.log("WebSocket connected.");
    ws.send(code);
};

ws.onerror = function(error) {
    console.error('WebSocket Error:', error);
};

function changeColor() {
    var button = document.getElementById('colorButton');
    var newColor = button.style.backgroundColor === 'green' ? 'red' : 'green';
    ws.send(newColor);
    button.style.backgroundColor = newColor;
}
ws.onmessage = function(event) {
    var button = document.getElementById('colorButton');
    console.log("recup :",event.data)
    color.innerText = (event.data)
    button.style.backgroundColor = event.data;
};
ws.onclose = function(event) {
    console.log("websocket closed")
}