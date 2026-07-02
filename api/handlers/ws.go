package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type WSAddMessage struct {
	ChatID       int    `json:"chat_id"`
	ChatMemberId int    `json:"from_chat_member_id"`
	Text         string `json:"text"`
	Media        string `json:"media"`
}

type WSSearchPayload struct {
	Text string `json:"query"`
}

type WSMessageWrapper struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type WSCreatePayload struct {
	UserId int `json:"userId"`
}

type WSCreateFolderPayload struct {
	UserId int    `json:"user_id"`
	Name   string `json:"folder_name"`
}

type WSEditFolderNamePayload struct {
	UserId   int    `json:"user_id"`
	PrevName string `json:"prev_name"`
	NewName  string `json:"new_name"`
}

type WSAddChatToFolderPayload struct {
	UserId     int    `json:"user_id"`
	ChatId     int    `json:"chat_id"`
	Foldername string `json:"folder_name"`
}

type WSRemoveChatFromFolderPayload struct {
	UserId     int    `json:"user_id"`
	ChatId     int    `json:"chat_id"`
	Foldername string `json:"folder_name"`
}

type WSToggleChatInFolderPayload struct {
	UserId     int    `json:"user_id"`
	ChatId     int    `json:"chat_id"`
	Foldername string `json:"folder_name"`
	IsChecked  bool   `json:"is_checked"`
}

type WSEditChatNamePayload struct {
	ChatId  int    `json:"chat_id"`
	NewName string `json:"new_name"`
}

type WSEditProfilePayload struct {
	UserId      int    `json:"user_id"`
	Login       string `json:"login"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Birthday    string `json:"birthday"`
	Sex         string `json:"sex"`
}

func (A *API) wsHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	claims, err := A.AuthService.Jwt.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	exists, err := A.AuthService.Storage.IsTokenExists(tokenString)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "token revoked", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID

	// TODO

	conn, err := A.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS upgrade error:", err)
		return
	}

	A.WsClients.Lock()
	A.WsClients.Conns[userID] = conn
	A.WsClients.Unlock()

	allMessages, err := A.loadAllMessages(userID)
	if err != nil {
		log.Println("load messages error:", err)
		conn.WriteJSON(map[string]interface{}{
			"type":  "error",
			"error": "failed to load messages",
		})
		return
	}

	conn.WriteJSON(map[string]interface{}{
		"type":          "initial_state",
		"initial_state": allMessages,
	})

	go func() {
		defer func() {
			A.WsClients.Lock()
			delete(A.WsClients.Conns, userID)
			A.WsClients.Unlock()
			conn.Close()
			log.Println("WS disconnected:", userID)
		}()

		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				// TODO here must be reconnect
				log.Println("WS read error:", err)
				return
			}

			log.Println("Raw WS message:", string(p))

			var wrapper WSMessageWrapper
			if err := json.Unmarshal(p, &wrapper); err != nil {
				log.Println("JSON parse error:", err)
				return
			}

			switch wrapper.Type {

			case "new_message":
				var msg WSAddMessage
				if err := json.Unmarshal(wrapper.Payload, &msg); err != nil {
					log.Println("new_message parse error:", err)
					continue
				}
				A.addMessageHandler(userID, msg)

			case "search":
				log.Println("Nu pozya")

				var payload WSSearchPayload
				if err := json.Unmarshal(wrapper.Payload, &payload); err != nil {
					log.Println("search payload parse error:", err)
					continue
				}
				A.searchUsersHandler(userID, payload.Text)

			case "create_chat":
				var createPayload WSCreatePayload

				if err := json.Unmarshal(wrapper.Payload, &createPayload); err != nil {
					log.Println("create_chat parse error:", err)
					continue
				}

				A.createPrivateChatHandler(userID, createPayload.UserId)

			case "create_folder":
				var createFolderPayload WSCreateFolderPayload

				if err := json.Unmarshal(wrapper.Payload, &createFolderPayload); err != nil {
					log.Println("create folder parse error: ", err)
					continue
				}

				A.createFolderHandler(createFolderPayload.UserId, createFolderPayload.Name)

			case "add_chat_to_folder":
				var AddChatToFolderPayload WSAddChatToFolderPayload

				if err := json.Unmarshal(wrapper.Payload, &AddChatToFolderPayload); err != nil {
					log.Println("create folder parse error: ", err)
					continue
				}

				A.addChatToFolderHandler(AddChatToFolderPayload.UserId, AddChatToFolderPayload.ChatId, AddChatToFolderPayload.Foldername)

			case "remove_chat_from_folder":
				var RemoveChatFromFolderPayload WSRemoveChatFromFolderPayload

				if err := json.Unmarshal(wrapper.Payload, &RemoveChatFromFolderPayload); err != nil {
					log.Println("remove chat from folder parse error: ", err)
					continue
				}

				A.removeChatFromFolderHandler(RemoveChatFromFolderPayload.UserId, RemoveChatFromFolderPayload.ChatId, RemoveChatFromFolderPayload.Foldername)

			case "toggle_chat_in_folder":
				var ToggleChatInFolderPayload WSToggleChatInFolderPayload

				if err := json.Unmarshal(wrapper.Payload, &ToggleChatInFolderPayload); err != nil {
					log.Println("toggle chat in folder parse error: ", err)
					continue
				}

				if ToggleChatInFolderPayload.IsChecked {
					A.addChatToFolderHandler(ToggleChatInFolderPayload.UserId, ToggleChatInFolderPayload.ChatId, ToggleChatInFolderPayload.Foldername)
				} else {
					A.removeChatFromFolderHandler(ToggleChatInFolderPayload.UserId, ToggleChatInFolderPayload.ChatId, ToggleChatInFolderPayload.Foldername)
				}

			case "edit_chat_name":
				var editChatNamePayload WSEditChatNamePayload

				if err := json.Unmarshal(wrapper.Payload, &editChatNamePayload); err != nil {
					log.Println("create folder parse error: ", err)
					continue
				}

				A.editChatNameHandler(editChatNamePayload.ChatId, editChatNamePayload.NewName)

			case "edit_folder_name":

				var editFolderNamePayload WSEditFolderNamePayload

				if err := json.Unmarshal(wrapper.Payload, &editFolderNamePayload); err != nil {
					log.Println("create folder parse error: ", err)
					continue
				}

				A.editFolderNameHandler(editFolderNamePayload.UserId, editFolderNamePayload.PrevName, editFolderNamePayload.NewName)

			case "edit_profile":

				var payload map[string]string

				if err := json.Unmarshal(wrapper.Payload, &payload); err != nil {
					log.Println("create folder parse error: ", err)
					continue
				}

				A.editProfileHandler(payload)

			default:
				log.Println("Unknown WS message type:", wrapper.Type)
			}
		}
	}()
}

func (A *API) WsHandler(w http.ResponseWriter, r *http.Request) {
	A.wsHandler(w, r)
}
