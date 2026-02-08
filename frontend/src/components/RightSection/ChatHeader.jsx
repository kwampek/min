import { GetAvatar } from "../MiddleSection/ChatItem.jsx"

const ChatHeader = ( {activeChat} ) => {
    if (activeChat == null) {
        return;
    }

    return (
        <div className="chat-header">
            <div className="chat-user">

            <GetAvatar chat={activeChat}/>

            <div className="info">
                <span className="name">{activeChat.name}</span>
            </div>
            </div>
            <div className="actions">
            <button className="action-btn" title="Call">📞</button>
            <button className="action-btn" title="More">⋯</button>
            </div>
        </div>
    );
}

export default ChatHeader;