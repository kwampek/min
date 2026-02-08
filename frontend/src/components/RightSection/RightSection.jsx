import { useSelector } from "react-redux";
import { useState } from "react";
import ChatHeader from "./ChatHeader";
import Messages from "./Messages";


const RightSection = ({ date, handleAddMessage }) => {
  const [inputValue, setInputValue] = useState('');

  const activeChat = useSelector((state) => state.chats.activeChat);
  const messages = useSelector(
    (state) => (activeChat ? state.chats.messages[activeChat.chat_id] || [] : [])
  );

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

  return (
    <div className="right-section">
      <ChatHeader activeChat={activeChat}/>

      <div className="date-wrapper">
        <div className="date">{date}</div>
      </div>

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
