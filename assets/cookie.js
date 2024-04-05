var socket = new WebSocket("ws://localhost:8080/ws");
        socket.onopen = function (event) {
            // Receive the unique identifier from the server
            socket.onmessage = function (event) {
                var uid = event.data;
                // Set a cookie with the received unique identifier
                document.cookie = "uid=" + uid + ";path=/";
            };
        };