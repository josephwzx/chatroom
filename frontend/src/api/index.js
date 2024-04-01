var socket = new WebSocket("ws://localhost:8080/ws");

let connect = cb => {
    console.log("Connectting...");

    socket.onopen = () => {
        console.log("Successfully Connected");
    };

    socket.onmessage = msg => {
        console.log(msg);
        cb(msg);
    };

    socket.onclose = event => {
        console.log("Socket Closed Connection: ", event);
    };

    socket.onerror = error => {
        console.log("Socket Error: ", error);
    };
};

let sendMsg = (msg, username) => {
    const message = {
        username: username,
        message: msg
    };
    console.log("sending msg: ", message);
    socket.send(JSON.stringify(message));
};

export { connect, sendMsg };