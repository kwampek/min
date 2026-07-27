import { useEffect, useRef } from "react";

const DateDivider = ({ date }) => {
    const messageDate = new Date(date);
    const today = new Date();

    let text = messageDate.toLocaleDateString("ru-RU", {
        day: "numeric",
        month: "long",
        year:
            messageDate.getFullYear() !== today.getFullYear()
                ? "numeric"
                : undefined,
    });

    return (
        <div className="date-wrapper">
            <div className="date">{text}</div>
        </div>
    );
};

const formatTime = (date) =>
    new Date(date).toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
    });

const Message = ({ msg }) => {
    return (
        <div className={`message-row ${msg.mine ? "mine" : "other"}`}>

            {!msg.mine && (
                <div className="avatar-slot">
                    {msg.showAvatar && (
                        <img
                            className="avatar"
                            src={msg.avatar}
                            alt=""
                        />
                    )}
                </div>
            )}

            <div className={`message ${msg.mine ? "mine" : "other"} ${msg.group}`}>
                {msg.text}
                <span className="time">{formatTime(msg.time)}</span>
            </div>

            {msg.mine && (
                <div className="avatar-slot">
                    {msg.showAvatar && (
                        <img
                            className="avatar"
                            src={msg.avatar}
                            alt=""
                        />
                    )}
                </div>
            )}
        </div>
    );
};

const Messages = ({ messages }) => {
    const refBottom = useRef(null);

    useEffect(() => {
        refBottom.current?.scrollIntoView({ behavior: "smooth" });
    }, [messages]);

    return (
        <div className="messages">
            {messages.map((msg) => (
                <div key={msg.id}>
                    {msg.showDateDivider && (
                        <DateDivider date={msg.time} />
                    )}

                    <Message msg={msg} />
                </div>
            ))}

            <div ref={refBottom} />
        </div>
    );
};


export default Messages;