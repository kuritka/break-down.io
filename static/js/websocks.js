var socket = new WebSocket("ws://localhost:8080/api/websockets/stream");
const avatar = document.getElementById("avatar-container");

socket.onopen = function () {
    console.log("Connection opened");
    avatar.setAttribute("class", "avatar-container avatar-container-live")
}

socket.onmessage = function (e) {
    console.log(e.data)
    switch(e.data) {
        case "off":
            avatar.setAttribute("class", "avatar-container avatar-container-down");
            break;
        default:
            avatar.setAttribute("class", "avatar-container avatar-container-live");
    }
}

socket.onerror = function (e) {
    console.log("error occured" + e);
}

socket.onclose = function () {
    console.log("Connection closed");
    avatar.setAttribute("class", "avatar-container avatar-container-down")
}

