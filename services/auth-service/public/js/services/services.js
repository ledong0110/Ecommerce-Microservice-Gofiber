
const CHAT_SERVER_ENDPOINT = window.location.host;
let webSocketConnection = null;
var selectedUser = ""

// Chat list
var chatList = []

const chatListSubscription = (socketPayload) => {
    let newChatList = chatList;

    if (socketPayload.type === 'new-user-joined') {
      const incomingChatList = socketPayload.chatlist;
      if (incomingChatList) {
        newChatList = newChatList.filter(
          (obj) => obj.userID !== incomingChatList.userID
        );
      }

      /* Adding new online user into chat list array */
      newChatList = [...newChatList, ...[incomingChatList]];
    } else if (socketPayload.type === 'user-disconnected') {
      const outGoingUser = socketPayload.chatlist;
      const loggedOutUserIndex = newChatList.findIndex(
        (obj) => obj.userID === outGoingUser.userID
      );
      if (loggedOutUserIndex >= 0) {
        newChatList.splice(loggedOutUserIndex, 1);
      }
    } else {
      newChatList = socketPayload.chatlist;
    }

    // slice is used to create aa new instance of an array
    chatList = newChatList.slice();

    var onlineUsers = document.querySelectorAll('.parent-user')
    for (let i = 0; i < onlineUsers.length; i++)
    {
        var user = onlineUsers[i].querySelector('.child')
        var idOnlineUser = onlineUsers[i].getAttribute('data-id')
        if (chatList.filter(e => e.ID === idOnlineUser).length > 0){
            user.classList.remove("bg-danger")
            user.classList.add("bg-success")
        }
        else {
            user.classList.remove("bg-success")
            user.classList.add("bg-danger")
        }
    }
};



// Conversation

var conversations = []


const scrollMessageContainer = (messageContainer) => {
    if (messageContainer !== null) {
      try {
        setTimeout(() => {
          messageContainer.scrollTop = messageContainer.scrollHeight;
        }, 100);
      } catch (error) {
        console.warn(error);
      }
    }
  }

const newMessageSubscription = (messagePayload) => {
    if (
        selectedUser !== null &&
        selectedUser === messagePayload.fromUserID
    ) {

        conversations = [...conversations, messagePayload];
        RenderUIConversation()
    }
};

const loadConversation = async (toUserID) => {
    selectedUser = toUserID
    let res = await getConversationBetweenUsers(toUserID)
    if (res == null) {
        conversations = []
    }
    else {
        conversations = res.slice()
    }
    RenderUIConversation()
}

async function getConversationBetweenUsers(toUserID) {
    var response
    var dataStr = {data: "hello"}
    await $.ajax({
        type: "POST",
        url: `/chat-app/GetConversation/${toUserID}`,
        data: dataStr,
        success: (res)=>{
            response = JSON.parse(res)
        }

    })
    return response
}

function RenderUIConversation() {
    var ConversationTag = document.getElementById("conversation")
    var html = ""
    for(let i = 0; i < conversations.length; i++) {
        if (conversations[i].fromUserID == userID.id){
            html += `<div class="d-flex flex-row justify-content-end">
            <div>
              <p class="small p-2 me-3 mb-1 text-white rounded-3 bg-primary">${conversations[i].message}</p>
              <p class="small me-3 mb-3 rounded-3 text-muted">12:00 PM | Aug 13</p>
            </div>
            <img src="https://mdbcdn.b-cdn.net/img/Photos/new-templates/bootstrap-chat/ava1-bg.webp"
              alt="avatar 1" style="width: 45px; height: 100%;">
          </div>`
        }
        else {
            html += `
            <div class="d-flex flex-row justify-content-start">
            <img src="https://mdbcdn.b-cdn.net/img/Photos/new-templates/bootstrap-chat/ava6-bg.webp"
                alt="avatar 1" style="width: 45px; height: 100%;">
            <div>
                <p class="small p-2 ms-3 mb-1 rounded-3" style="background-color: #f5f6f7;">${conversations[i].message}</p>
                <p class="small ms-3 mb-3 rounded-3 text-muted float-end">12:00 PM | Aug 13</p>
            </div>
            </div>`
        }
    }
    if (html == "") html = "<h1>you have not chat, say hi !</h1>"
    ConversationTag.innerHTML = html
    scrollMessageContainer(ConversationTag)
}

function sendMessage() {
        const message = document.getElementById("messageText").value
        document.getElementById("messageText").value = ""
        if (message === '' || message === undefined || message === null) {
        alert(`Message can't be empty.`);
        } else if (userID.id === '') {
        alert("There are some errors");
        } else if (selectedUser === undefined) {
        alert(`Select a user to chat.`);
        } else {
        message.value = '';

        const messagePayload = {
            fromUserID: userID.id,
            message: message.trim(),
            toUserID: selectedUser,
        };
        sendWebSocketMessage(messagePayload);
        conversations = [...conversations, messagePayload];
        RenderUIConversation()
    }
}


// Socket service
function connectToWebSocket(userID) {
    if (userID === "" && userID === null && userID === undefined) {
        return {
            message: "You need User ID to connect to the Chat server",
            webSocketConnection: null
        }
    } else if (!window["WebSocket"]) {
        return {
            message: "Your Browser doesn't support Web Sockets",
            webSocketConnection: null
        }
    }
    if (window["WebSocket"]) {
        if (window.location.protocol === 'https:'){

            webSocketConnection = new WebSocket("wss://" + CHAT_SERVER_ENDPOINT + "/chat-app/ws/register");
        }
        else {
            webSocketConnection = new WebSocket("ws://" + CHAT_SERVER_ENDPOINT + "/chat-app/ws/register");
        }
        return {
            message: "You are connected to Chat Server",
            webSocketConnection
        }
    }

}

function sendWebSocketMessage(messagePayload) {
    if (webSocketConnection === null) {
      return;
    }
    webSocketConnection.send(
      JSON.stringify({
        eventName: 'message',
        eventPayload: messagePayload
      })
    );
}

function emitLogoutEvent(userID) {
    if (webSocketConnection === null) {
        return;
    }
    webSocketConnection.close();
}

function listenToWebSocketEvents() {

    if (webSocketConnection === null) {
        return;
    }

    webSocketConnection.onclose = (event) => {
        // eventEmitter.emit('disconnect', event);
    };

    webSocketConnection.onmessage = (event) => {
        try {
            const socketPayload = JSON.parse(event.data);
            switch (socketPayload.eventName) {
                case 'chatlist-response':
                    if (!socketPayload.eventPayload) {
                        return
                    }
                    chatListSubscription(socketPayload.eventPayload)
                    break;

                case 'disconnect':
                    if (!socketPayload.eventPayload) {
                        return
                    }
                    chatListSubscription(socketPayload.eventPayload)
                    break;
                    // eventEmitter.emit(
                    //   'chatlist-response',
                    //   socketPayload.eventPayload
                    // );


                case 'message-response':

                    if (!socketPayload.eventPayload) {
                        return
                    }    
                    newMessageSubscription(socketPayload.eventPayload)
                    break;

                default:
                    break;
            }
        } catch (error) {
            console.log(error)
            console.warn('Something went wrong while decoding the Message Payload')
        }
    };
}


function emitLogoutEvent() {
    if (webSocketConnection === null) {
        return;
    }
    webSocketConnection.close();
}
