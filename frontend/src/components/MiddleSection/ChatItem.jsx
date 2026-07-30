import React from 'react';

export const GetAvatar = ( {chat} ) => {
  if (!chat) {
    return <img className="chat-avatar" src="icons/me.jpg" />
  }

  if (!chat.avatarPreview || !chat.avatarPreview.Valid) {
    return (
    <div className="chat-avatar default-avatar">
      <span className="avatar-letter">{chat.title?.[0] || ''}</span>
    </div>
    );
  }

  return <img className="chat-avatar" src={chat.avatar.String} />
}

const ChatItem = ({ chat, isActive, onClick, onContextMenu }) => {
  
  return (
    <div className={`chat-item ${isActive ? 'active' : ''}`}
      onClick={onClick}
      onContextMenu={onContextMenu}
    >
      <GetAvatar chat={chat}/>

      <div className="chat-info">
        <div className="chat-top">
          <div className="chat-name">{chat.title}</div>
          <div className="chat-status">
            <span className="chat-time">{}</span>
          </div>
        </div>
        <div className="chat-last-message">{chat.lastMessage}</div>
      </div>
    </div>
  );
};

export default ChatItem;
