import store from "../store/index.js";
import { setChats, addMessage, normalizeChats, setSearchResults, addNewChat } from '../store/slices/chatSlice';


let socket = null;

//const WS_URL = "ws://158.160.185.12/ws";
const WS_URL = "ws://localhost:8080/ws";


function connectWS() {
  const raw = localStorage.getItem("userData");
  if (!raw) {
    return null;
  }

  let userData = JSON.parse(raw);

  const token = userData?.token;
  if (!token) {
    localStorage.removeItem("userData");
    return null;
  }

  return new WebSocket(`${WS_URL}?token=${encodeURIComponent(token)}`);
}



export function initWebSocket() {
  if (socket) return socket;
  socket = connectWS();
  if (!socket) {
    console.log("Failed to connect")
    return;
  }

  socket.onopen = () => console.log("WS CONNECTED");
  socket.onclose = () => console.log("WS CLOSED");
  socket.onerror = (e) => console.log("WS ERROR", e);

  
  socket.onmessage = (event) => {
    try {
        const msg = JSON.parse(event.data);
        if (msg.type === "unauthorized") {
            localStorage.removeItem("userData");
        }

        switch (msg.type) {
        case "initial_state": {
            const normalized = normalizeChats(msg.initial_state);
            store.dispatch(setChats(normalized));
            break;
        }

        case "new_message": {
            store.dispatch(
            addMessage({
                chat_id: msg.chat_id,
                message: {
                id: msg.message.message_id,
                chat_id: msg.chat_id,
                text: msg.message.message_text,
                time: msg.message.send_time,
                from_chat_member_id: msg.message.from_chat_member_id,
                mine: false,
                media: msg.message.media_id || null,
                status: msg.message.status
                }
            })
            );
            break;
        }

        case "search": {
            store.dispatch(setSearchResults(msg.payload.users));
            break;
        }

        case "create_chat": {
            store.dispatch(addNewChat(msg.payload))
            break;
        }

        default:
            break;
        }


    } catch (err) {
        console.error("WS parse error:", err);
    }
  };

  return socket;
}

export function getWebSocket() {
  return socket;
}
