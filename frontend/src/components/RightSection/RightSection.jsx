import { useSelector } from "react-redux";
import { useState } from "react";
import ChatHeader from "./ChatHeader";
import Messages from "./Messages";
import ChatDialog from "../Dialog/ChatDialog";
import { selectPreparedMessages } from "../../store/slices/chatSlice";


const RightSection = ({ handleAddMessage }) => {
  const [inputValue, setInputValue] = useState('');
  const [editMode, setEditMode] = useState('');  

  const activeChat = useSelector((state) => state.chats.activeChat);
  const messages = useSelector(selectPreparedMessages);

  const handleSend = () => {
    if (!inputValue.trim()) return;
    handleAddMessage(inputValue, true)
  
    setInputValue('');
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter') handleSend();
  };

  if (activeChat == null) {
    return (
      <div className="right-section">
      </div>
    );
  }

  const handleEditSubmit = ({ chatId, name, description, photo, addedMembers, removedMembers }) => {
    
    // TODO compare with prev values
    sendWsAction('edit_chat', {
        chat_id: chatId,
        title: name,
        desc: description,
        avatar: photo,
        added_memeber_ids: addedMembers,
        removed_member_ids: removedMembers,
      });

    setCreateMode(null);
    closePanel();
  };

  return (
    <div className="right-section">
      <ChatHeader
          activeChat={activeChat}
          onEdit={() => setEditMode("edit")}
      />

      {editMode && (
        <ChatDialog
          mode="edit"
          type={activeChat.type}
          onClose={() => {setEditMode(null)}}
          onSubmit={handleEditSubmit}
        />
      )}

      <Messages messages={messages} />

      <div className="input-area">
        <img className="icon" src='icons/attach_icon.png' title="Attach file" alt="Attach" />
        <img className="icon" src='icons/voice_icon.png' title="Record voice" alt="Voice" />
        <input
          type="text"
          placeholder="Write a message..."
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyPress={handleKeyPress}
        />
        <button onClick={handleSend} title="Send">➤</button>
      </div>
    </div>
  );
};

export default RightSection;
