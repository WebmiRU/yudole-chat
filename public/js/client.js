const messages = document.querySelector('main')
const messageTemplate = document.querySelector('template.message');
const successTemplate = document.querySelector('template.success');
// const errorTemplate = document.querySelector('template.error');
const infoTemplate = document.querySelector('template.info');
const messageAudio = document.querySelector('audio.message');

let socket = new ReconnectingWebSocket("ws://" + location.hostname + ":5367/chat");
let template = null

function message(msg) {
    switch (msg.type) {
        case "message/channel":
            template = document.createElement('div');
            template.appendChild(messageTemplate.content.cloneNode(true));
            template.innerHTML = template.innerHTML
                .replace(/{{service}}/g, msg.service)
                .replace(/{{user_name}}/g, msg.user.nick ?? msg.user.login)
                .replace(/{{content}}/g, msg.message.html);

            messages.prepend(template);
            break

        case "success/connection/server":
            template = document.createElement('div');
            template.appendChild(successTemplate.content.cloneNode(true));
            template.innerHTML = template.innerHTML
                .replace(/{{service}}/g, msg.service)
                .replace(/{{content}}/g, `Успешное подключение к <u>серверу</u> <b>${msg.service}</b>`);

            messages.prepend(template);
            break

        case "success/join/channel":
            template = document.createElement('div');
            template.appendChild(successTemplate.content.cloneNode(true));
            template.innerHTML = template.innerHTML
                .replace(/{{service}}/g, msg.service)
                .replace(/{{content}}/g, `Успешное подключение к <u>каналу</u> <b>${msg.channel}</b>`);

            messages.prepend(template);
            break

        case "user/join/channel":
            template = document.createElement('div');
            template.appendChild(infoTemplate.content.cloneNode(true));
            template.innerHTML = template.innerHTML
                .replace(/{{service}}/g, msg.service)
                .replace(/{{content}}/g, `Пользователь <b>${msg.user.nick.length ? msg.user.nick : msg.user.login}</b> <u>зашел</u> на канал канал <b>${msg.channel}</b>`);

            messages.prepend(template);
            break

        default:
            break

    }

    messageAudio.play()
}

const msg21 = {"type": "user/join/channel", "service": "twitch", "text": "", "user": {"login": "xxxyyy", "nick": "", "avatar_url": "", "color": ""}, "channel": "ewolf34", "value": ""};
const msg = {service: "twitch", type: "message/channel", user: {nick: "EWolf34", login: "ewolf34"}, message: {html: "Hello world!", text: "Hello world!"}}
const msg2 = {"type": "user/join/channel", "service": "twitch", "text": "", "user": {"login": "yudolebot", "nick": "", "avatar_url": "", "color": ""}, "channel": "ewolf34", "value": ""};
const msg3 = {"type": "success/join/channel", "service": "trovo", "user": {"login": "", "nick": "", "avatar_url": "", "color": ""}, "channel": "EWolf34"};
const msg4 = {"service": "twitch", "type": "message/channel", "user": {"login": "ewolf34", "nick": "EWolf34", "avatar_url": "", "color": ""}, "message": {"text": "432423432", "html": "432423432"}};

// setInterval(() => {
//     // message(msg4)
//     message(msg3)
//     // message(msg21)
// }, 2000);

socket.onmessage = function (event) {
    console.log(JSON.parse(event.data))
    message(JSON.parse(event.data))
};

socket.onopen = (event) => {
    console.log('Соединение с сервером установлено')
}

socket.onclose = (event) => {
    if (event.wasClean) {
        console.log('Соединение с сервером закрыто')
    } else {
        console.log('Соединение с сервером прервано')
    }
};

socket.onerror = (error) => {
    // socket.close();
    console.log(error);
};
