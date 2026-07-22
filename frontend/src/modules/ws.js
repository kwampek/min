import { useNavigate } from "react-router-dom";
import store from "../store/index.js";
import { setChats, addMessage, normalizeChats, setSearchResults, addNewChat } from '../store/slices/chatSlice';


let socket = null;

//const WS_URL = "ws://158.160.185.12/ws";
const WS_URL = "ws://localhost:8081/ws";


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

export function initWebSocket({ onUnauthorized }) {
  if (socket) return socket;
  socket = connectWS();
  if (!socket) {
    console.log("Failed to connect")
    return;
  }

  socket.onopen = () => {
    console.log("WS CONNECTED");
  };
  
  socket.onclose = (event) => {
      console.log("WS CLOSED");

      socket = null;

      if (event.code === 4001) {
          onUnauthorized?.();
      }
  };

  socket.onerror = (e) => {
    // TODO
    // theoretically there won't be errors errors during ws send 
    // server catches all errors and revokes connection

    console.log("WS ERROR", e);

    socket = null;

    if (event.code === 4001) {
        onUnauthorized?.();
    }
  };


  // better to greedy add objects and confirm them with ws msges
  
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
        // TODO may be need to return here
    }
  };

  return socket;
}

export function getWebSocket() {
  return socket;
}
