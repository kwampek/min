import { useEffect, useState, useRef } from 'react';

const Messages = ({messages}) => {
    const refNiz = useRef(null);

    useEffect(() => {
      refNiz.current?.scrollIntoView({ behavior: 'smooth' });
    }, [messages]);


    return (
      <div className="messages">
        {messages.map((msg) => (
          <div key={msg.id} className={`message ${msg.mine ? 'mine' : 'other'}`}>
            {msg.text} <span className="time">{msg.time}</span>
          </div>
        ))}
        <div ref={refNiz} />
      </div>
    );
}

export default Messages;