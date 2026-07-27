import { GetAvatar } from "../MiddleSection/ChatItem.jsx";

const ChatHeader = ({ activeChat, onEdit }) => {
    if (!activeChat) return null;

    return (
        <div className="chat-header" onClick={onEdit}>
            <div className="chat-user">
                <GetAvatar chat={activeChat} />

                <div className="info">
                    <span className="name">{activeChat.title}</span>
                </div>
            </div>

            <div className="actions">
                <button
                    className="action-btn"
                    title="Call"
                    onClick={(e) => e.stopPropagation()}
                >
                    📞
                </button>

                <button
                    className="action-btn"
                    title="More"
                    onClick={(e) => e.stopPropagation()}
                >
                    ⋯
                </button>
            </div>
        </div>
    );
};

export default ChatHeader;