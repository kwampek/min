import { useDispatch, useSelector } from 'react-redux';
import LeftSidebar from '../components/LeftSidebar/LeftSidebar';
import LeftPanel from '../components/LeftPanel/LeftPanel';
import MiddleSection from '../components/MiddleSection/MiddleSection';
import RightSection from '../components/RightSection/RightSection';
import '../assets/styles/main_style.css';
import { setActiveChat, addMessage } from '../store/slices/chatSlice';
import { getWebSocket, initWebSocket } from '../modules/ws';
import { useEffect, useState } from 'react';
import { useNavigate } from "react-router-dom";

const App = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { activeChat } = useSelector((state) => state.chats);
  const { userId } = useSelector((state) => state.user)

  const [isLeftPanelOpen, setLeftPanelOpen] = useState(false);

  const closePanel = () => setLeftPanelOpen(false);


  useEffect(() => {
      initWebSocket({
          onUnauthorized: () => navigate("/login"),
      });
  }, [navigate]);


  const ws = getWebSocket();

  const updateChatClick = (chat_id) => {
    dispatch(setActiveChat(chat_id));
  };

  const handleAddMessage = (text, mine = true) => {
    if (activeChat === null) return;

    const message = {
      id: Date.now(),
      chat_id: activeChat.chat_id,
      from_chat_member_id: activeChat.chat_member_id,
      text,
      time: new Date().toLocaleTimeString().slice(0, 5),
      mine,
    };

    ws.send(JSON.stringify({
      type: 'new_message',
      payload: message,
    }));

    dispatch(addMessage({ chat_id: activeChat.chat_id, message}));
  };


  const date = '17 декабря'; // okak

  return (
    <div className="container">
      <LeftSidebar
        isLeftPanelOpen={isLeftPanelOpen}
        setLeftPanelOpen={setLeftPanelOpen}
      />

      {isLeftPanelOpen && (
          <LeftPanel
              closePanel={closePanel}
          />
      )}

      <div className="resize-bar-left"></div>
      <MiddleSection activeChatId={activeChat} updateChatClick={updateChatClick} />
      <RightSection date={date} handleAddMessage={handleAddMessage}/>
    </div>
  );
};

export default App;
