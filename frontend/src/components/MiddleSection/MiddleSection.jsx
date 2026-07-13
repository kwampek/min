import { useDispatch, useSelector} from 'react-redux';
import ChatItem from './ChatItem';
import { getWebSocket } from '../../modules/ws';
import { useEffect, useState } from 'react';
import { setSearchResults, removeSearchResult } from '../../store/slices/chatSlice';
import ChatContextMenu from './ChatContextMenu';

const MiddleSection = ({ updateChatClick }) => {  
  const [searchText, setSearchText] = useState('');
  const [contextMenu, setContextMenu] = useState(null);

  const dispatch = useDispatch();
  const {folders, chats, activeChat, activeFolderName, searchResults } = useSelector((state) => state.chats)
  const user = useSelector((state) => state.user)
  const chatsInFolder = folders[activeFolderName] || [];
  const chats_list = chatsInFolder.map((id) => chats[id]).filter(Boolean);


  useEffect(() => {
    const close = (e) => {
      if (e.target.closest('.context-menu')) {
        return;
      }
      setContextMenu(null);
    };
    window.addEventListener("click", close);
    window.addEventListener("contextmenu", close);
    return () => {
      window.removeEventListener("click", close);
      window.removeEventListener("contextmenu", close);
    };
  }, []);

  const sendSearchRequest = (query) => {
    const ws = getWebSocket();

    if (!ws) return;
    
    ws.send(JSON.stringify({
      type: "search",
      payload: { userId: user.userId, query: query },
    }));
  };


  const handleCreateChat = (user) => {
    
    const ws = getWebSocket();

    if (!ws) return;
    
    ws.send(JSON.stringify({
      type: "create_chat",
      payload: { userId: user.user_id }
    }));

  };


  useEffect(() => {
    if (!searchText) {
      dispatch(setSearchResults([]));
      return;
    }

    const timeout = setTimeout(() => {
      sendSearchRequest(searchText);
    }, 300);

    return () => clearTimeout(timeout);
  }, [searchText, dispatch]);

  // ujas
  const existingChats = chats_list.filter(chat => chat.title.toLowerCase().includes(searchText.toLowerCase()));
  const globalSearchResults = searchResults.filter(user => !user.have_chat);

  return (
    <div className="middle-section">
      <div className="search-bar">
        <input type="text" value={searchText} onChange={e => setSearchText(e.target.value)} placeholder="Search..."></input>
      </div>

      {contextMenu && (
        <ChatContextMenu
          x={contextMenu.x}
          y={contextMenu.y}
          chat={contextMenu.chat}
          userId={user.userId}
          onClose={() => setContextMenu(null)}
        />
      )}

      {existingChats.map(chat => (
        <ChatItem
          key={chat.chat_id}
          chat={chat}
          isActive={activeChat?.chat_id === chat.chat_id}
          onClick={() => updateChatClick(chat.chat_id)}
          onContextMenu={(e) => {
            e.preventDefault();
            e.stopPropagation();
            setContextMenu({
              x: e.clientX,
              y: e.clientY,
              chat
            });
          }}
        />
      ))}

      {globalSearchResults.length > 0 && <div className="search-separator">Глобальный поиск</div>}

      {globalSearchResults.map(user => (
        <ChatItem
          key={`new-${user.userId}`}
          chat={user}
          isActive={false}
          onClick={() => {
            dispatch(removeSearchResult(user.userId))
            handleCreateChat(user)}
          }
        />
      ))}

      {existingChats.length === 0 && globalSearchResults.length === 0 && searchText && (
        <div className="search-separator">Ничего не найдено</div>
      )}

    </div>
  );
};

export default MiddleSection;
