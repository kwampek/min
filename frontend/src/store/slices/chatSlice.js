import { createSlice } from '@reduxjs/toolkit';

const initialState = {
  messages: {}, // chat_id: [{ id, chat_id, text, time, mine, media, status }]
  chats: {},    // chat_id: { from_chat_member,  id, name, avatar, lastMessage, unread }
  folders: {},  // folderName: [chat_id, ...]
  activeChat: null,
  activeFolderName: 'All',
  searchResults: [],
};

const chatSlice = createSlice({
  name: 'chats',
  initialState,
  reducers: {
    setChats(state, action) {
      const { messages, chats, folders } = action.payload;
      state.messages = messages;
      state.chats = chats;
      state.folders = folders;

      if (!state.folders["All"]) state.folders["All"] = [];

      if (!state.activeChat && folders["All"]?.length > 0) {
        const firstChatId = folders["All"][0];
        state.activeChat = chats[firstChatId];
      }
    },
    setActiveChat(state, action) {
      const chatId = action.payload;
      state.activeChat = state.chats[chatId] || null;
    },
    addMessage(state, action) {
      const { chat_id, message } = action.payload;
      if (!state.messages[chat_id]) {
        state.messages[chat_id] = [];
      }
      state.messages[chat_id].push(message);

      if (state.chats[chat_id]) {
        state.chats[chat_id].lastMessage = message.text;
      }

      if (state.activeChat?.chat_id === chat_id) {
        state.activeChat = {
          ...state.activeChat,
          lastMessage: message.text,
        };
      }
    },
    setActiveFolder(state, action) {
      state.activeFolderName = action.payload;
    },
    addFolder(state, action) {
      state.folders[action.payload] = [];
    },
    addChatToFolder(state, action) {
      const { folder, user_preview } = action.payload;
      if (!state.chatsByFolders[folder]) state.chatsByFolders[folder] = [];
      state.chatsByFolders[folder].push(user_preview);
    },
    setSearchResults(state, action) {
      state.searchResults = Array.isArray(action.payload) ? action.payload : [];
    },
    addNewChat(state, action) {
      const chat = action.payload;
      state.chats[chat.chat_id] = chat;

      if (!state.folders["All"]) state.folders["All"] = [];
      state.folders["All"].push(chat.chat_id);

      setActiveFolder(state, "All");
      setActiveChat(state, chat.chat_id)
    },
    removeSearchResult(state, action) {
      const userId = action.payload;
      state.searchResults = state.searchResults.filter(u => u.userId !== userId);
    },
    editFolderName(state, action) {
      const { prev_name, new_name } = action.payload;

      if (!state.folders[prev_name]) return;

      state.folders[new_name] = state.folders[prev_name];
      delete state.folders[prev_name];

      if (state.activeFolderName === prev_name) {
        state.activeFolderName = new_name;
      }
    },

    editChatName(state, action) {
      const { chat_id, new_name } = action.payload;
      if (!state.chats[chat_id]) return;

      state.chats[chat_id].name = new_name;

      if (state.activeChat?.chat_id === chat_id) {
        state.activeChat = {
          ...state.activeChat,
          name: new_name,
        };
      }
    },
    moveChatToFolder(state, action) {
      const { chat_id, fromFolder, toFolder } = action.payload;

      if (!state.folders[fromFolder] || !state.folders[toFolder]) return;

      state.folders[fromFolder] =
        state.folders[fromFolder].filter(id => id !== chat_id);

      if (!state.folders[toFolder].includes(chat_id)) {
        state.folders[toFolder].push(chat_id);
      }
    },
    toggleChatInFolder(state, action) {
      const { chat_id, folderName, isChecked } = action.payload;

      if (!state.folders[folderName]) return;

      if (isChecked) {
        if (!state.folders[folderName].includes(chat_id)) {
          state.folders[folderName].push(chat_id);
        }
      } else {
        state.folders[folderName] = state.folders[folderName].filter(id => id !== chat_id);
      }
    }
  }
});

export const { setChats, setActiveChat, addMessage, setActiveFolder, addFolder, addChatToFolder, setSearchResults, addNewChat, removeSearchResult, editFolderName, editChatName, moveChatToFolder, toggleChatInFolder } = chatSlice.actions;
export default chatSlice.reducer;

export const normalizeChats = (data) => {

  const messages = {};
  const chats = {};
  const folders = {};

  Object.entries(data.chatsByFolders).forEach(([folderName, chatList]) => {
    folders[folderName] = [];
    chatList.forEach(chat => {
      chats[chat.chat_id] = {
        chat_member_id: chat.chat_member_id,
        chat_id: chat.chat_id,
        name: chat.name,
        avatar: chat.avatar.String,
        unread: chat.unread || 0
      };

      folders[folderName].push(chat.chat_id);
    });
  });

  Object.entries(data.messagesByChats).forEach(([chat_id, messagesByChat]) => {
      messages[chat_id] = (messagesByChat || []).map(msg => ({
          id: msg.message_id,
          text: msg.message_text,
          time: msg.send_time,
          from_chat_member_id: msg.from_chat_member_id,
          mine: msg.from_chat_member_id === chats[chat_id].chat_member_id,
          media: msg.media_id || null,
          status: msg.status
        }));
      
        if (messages[chat_id].length > 0) {
          chats[chat_id].lastMessage = messages[chat_id][messages[chat_id].length - 1].text;
        } else {
          chats[chat_id].lastMessage = "";
        }
  });

  return { messages, chats, folders };
};
