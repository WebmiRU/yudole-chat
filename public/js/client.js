let messages = document.querySelector('main')
let socket = new ReconnectingWebSocket("ws://" + location.hostname + ":5367/chat");

// function connect() {
//     try {
//         let socket = new WebSocket("ws://" + location.hostname + ":5367/chat");
//
//         socket.onopen = function (e) {
//             console.log("Успешное подключение к серверу чата")
//         };
//     } catch {
//         console.log("SOCKET ERROR");
//     }
// }

// setInterval(function () {
//     try {
//         socket.send('0');
//     } catch {
//         console.log('SERROR');
//         socket = new WebSocket("ws://" + location.hostname + ":5367/chat");
//     }
// }, 3000);

// // Debug function
// setInterval(function () {
//     let m = {
//         "service": "twitch",
//         "type": "message/channel",
//         "user": {
//             "login": "ewolf34",
//             "nick": "EWolf34",
//             "avatar_url": "",
//             "color": ""
//         },
//         "message": {
//             "text": "Проверка "+Math.random()+" HahaSweat HahaSweat HahaSweat HahaSweat HahaSweat  123",
//             "html": "Проверка" +Math.random()+ "<img class=\"smile smile-twitch\" src=\"https://static-cdn.jtvnw.net/emoticons/v2/301108037/default/dark/1.0\" alt=\"HahaSweat\"/> <img class=\"smile smile-twitch\" src=\"https://static-cdn.jtvnw.net/emoticons/v2/301108037/default/dark/1.0\" alt=\"HahaSweat\"/> <img class=\"smile smile-twitch\" src=\"https://static-cdn.jtvnw.net/emoticons/v2/301108037/default/dark/1.0\" alt=\"HahaSweat\"/> <img class=\"smile smile-twitch\" src=\"https://static-cdn.jtvnw.net/emoticons/v2/301108037/default/dark/1.0\" alt=\"HahaSweat\"/> <img class=\"smile smile-twitch\" src=\"https://static-cdn.jtvnw.net/emoticons/v2/301108037/default/dark/1.0\" alt=\"HahaSweat\"/>  123"
//         }
//     };
//     let message = `<div class="message channel ${m.service}"><p class="user">${m.user.nick ?? m.user.login}</p><p class="content">${m.message.html ?? m.message.text}</p></div>`
//     messages.insertAdjacentHTML('afterend', message)
//
// }, 3000)

socket.onmessage = function (event) {
    let m = JSON.parse(event.data)
    switch (m.type) {
        case "message/channel":
            let message = `<div class="message channel ${m.service}"><p class="user">${m.user.nick ?? m.user.login}</p><p class="content">${m.message.html ?? m.message.text}</p></div>`
            console.log("MESSAGE APPEND", message)
            messages.insertAdjacentHTML('afterend', message)
            break

        case "channel/join/success":
            let message2 = `<div class="message system success ${m.service}"><p class="user">${m.service}</p><p class="content">${m.text}</p></div>`
            console.log("MESSAGE APPEND", message2)
            messages.insertAdjacentHTML('afterend', message2)
            break

        default:
            break

    }

    console.log(m)
};

socket.onclose = (event) => {
    if (event.wasClean) {
        console.log('Соединение с сервером закрыто')
    } else {
        console.log('Соединение с сервером прервано')
    }

    // console.log("RECONNECTING");
    // console.log("RECONNECTING FUNC");
    // socket = new WebSocket("ws://localhost:5367/chat");
    // console.log("RECONNECTING DONE");
    // // setTimeout(() => {
    // //     console.log("RECONNECTING FUNC");
    // //     socket = new WebSocket("ws://localhost:5367/chat");
    // //     console.log("RECONNECTING DONE");
    // // }, 1000);
    // setInterval(() => {
    //     console.log("RECONNECTING FUNC");
    //     socket = new WebSocket("ws://localhost:5367/chat");
    //     console.log("RECONNECTING DONE");
    // }, 5000);
};

socket.onerror = (error) => {
    // socket.close();
    console.log(error);
};
