var ws = new WebSocket("ws://82.67.64.82:80/ws");
var color = document.getElementById("jscolor")
var getcolor = document.getElementById("getcolor")
var getcode = document.getElementById("getcode")
var urlcourante = document.location.href;
var currentCode = urlcourante.split("/")[4]
var buttoncouleur = document.getElementById('colorButton');
var bouton1 = document.getElementById("1")
var bouton2 = document.getElementById("2")
var bouton3 = document.getElementById("3")
var bouton4 = document.getElementById("4")
var bouton5 = document.getElementById("5")
var bouton6 = document.getElementById("6")
var bouton7 = document.getElementById("7")
var bouton8 = document.getElementById("8")
var bouton9 = document.getElementById("9")
var statut = document.getElementById("statut")
let buttons = [[bouton1, bouton2, bouton3], [bouton4, bouton5, bouton6], [bouton7, bouton8,bouton9]]
var tour = true
var droitdejouer = true

ws.onopen = function(event) {
    console.log("WebSocket connected.");
};

ws.onerror = function(error) {
    console.error('WebSocket Error:', error);
};
        
function Win(symb) {
    for(var i = 0 ; i < 3 ; i++){
        if (buttons[i][0].innerText == symb && buttons[i][1].innerText == symb && buttons[i][2].innerText == symb) {
            return true;
        }
    }
	for (var col = 0; col < 3;col++) {
		if (buttons[0][col].innerText == symb && buttons[1][col].innerText == symb && buttons[2][col].innerText == symb) {
			return true;
		}
	}
	if (buttons[0][0].innerText == symb && buttons[1][1].innerText == symb && buttons[2][2].innerText == symb) {
		return true;
	}
	if (buttons[0][2].innerText == symb && buttons[1][1].innerText == symb && buttons[2][0].innerText == symb) {
		return true;
	}
}
     

function Egalite(){
    Draw = true
    for(var i = 0 ; i < 3 ; i++){
        for(var j = 0 ; j < 3 ; j++){
            if(buttons[i][j].innerText == ""){
                Draw = false
            }
        }
    }
    if (Draw){
        statut.textContent = "Egalité sale merde"  
    }   
}

function Morpion(nEvent){
    currentbtn = document.getElementById(nEvent.id)
    if (currentbtn.innerText == ""){
        if (tour && droitdejouer){
            droitdejouer = false
            currentbtn.innerText = "X"
            ws.send(getcode.value+"|"+"false"+"|"+ nEvent.id);
            if (Win("X")) {
                statut.textContent = "Le joueur X a gagné"
            } else {
                statut.textContent = "tour du j1"
            }
            
        } else if (droitdejouer){
            droitdejouer = false
            currentbtn.innerText = "0"
            ws.send(getcode.value+"|"+"true"+"|" + nEvent.id);
            if (Win("O")) {
                statut.textContent = "Le joueur O a gagné"
            } else {
                statut.textContent = "tour du j2"
            }
        }
    }
    Egalite()
       
}




function changeColor(nEvent) {
    var btntemp = document.getElementById(nEvent.id)
    var newColor = btntemp.style.backgroundColor === 'green' ? 'red' : 'green';
    console.log("tient :",newColor)
    ws.send(newColor + "|" + getcode.value+"|"+tour);
    console.log("hidden input",getcolor.value)
    console.log("pute",getcode.value)
    btntemp.style.backgroundColor = newColor;
}
ws.onmessage = function(event) {
    var button = document.getElementById('colorButton');
    var coderecu = event.data.split('|')[0].substring(1);
    var tourrecu = event.data.split('|')[1]
    var casecoche = event.data.split('|')[2]
    if (coderecu == currentCode) {
        console.log("recup :",event.data)
        if (tourrecu == "false") {
            tour = false
            document.getElementById(casecoche).innerText = "X"
        } else {
            tour = true
            document.getElementById(casecoche).innerText = "O"
        }
        droitdejouer = true
    }  
};
ws.onclose = function(event) {
    console.log("websocket closed")
}