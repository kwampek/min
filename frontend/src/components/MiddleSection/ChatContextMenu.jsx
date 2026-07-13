import { useDispatch, useSelector } from "react-redux";
import { editChatName, toggleChatInFolder } from "../../store/slices/chatSlice";
import { getWebSocket } from "../../modules/ws";
import { useEffect, useRef, useState } from "react";

const ChatContextMenu = ({ x, y, chat, userId, onClose }) => {
  const dispatch = useDispatch();
  const folders = useSelector(state => state.chats.folders);
  const ws = getWebSocket();
  const menuRef = useRef(null);
  const closeTimeoutRef = useRef(null);
  const [position, setPosition] = useState({ top: y, left: x });
  
  const getChatFolders = () => {
    const chatFolders = [];
    Object.keys(folders).forEach(folderName => {
      if (folderName !== "All" && folders[folderName].includes(chat.chat_id)) {
        chatFolders.push(folderName);
      }
    });
    return chatFolders;
  };
  
  const [selectedFolders, setSelectedFolders] = useState(new Set(getChatFolders()));
  
  useEffect(() => {
    const chatFolders = [];
    Object.keys(folders).forEach(folderName => {
      if (folderName !== "All" && folders[folderName].includes(chat.chat_id)) {
        chatFolders.push(folderName);
      }
    });
    setSelectedFolders(new Set(chatFolders));
  }, [folders, chat.chat_id]);

  const handleRename = () => {
    const newName = prompt("Новое имя чата:", chat.title);
    if (!newName) return;

    dispatch(editChatName({ chat_id: chat.chat_id, new_name: newName }));

    ws.send(JSON.stringify({
      type: "edit_chat_name",
      payload: { chat_id: chat.chat_id, new_name: newName }
    }));

    onClose();
  };

  const handleToggleFolder = (folderName) => {
    const isChecked = selectedFolders.has(folderName);
    const newSelectedFolders = new Set(selectedFolders);
    
    if (isChecked) {
      newSelectedFolders.delete(folderName);
    } else {
      newSelectedFolders.add(folderName);
    }
    
    setSelectedFolders(newSelectedFolders);

    dispatch(toggleChatInFolder({
      chat_id: chat.chat_id,
      folderName,
      isChecked: !isChecked
    }));

    ws.send(JSON.stringify({
      type: "toggle_chat_in_folder",
      payload: { 
        user_id: userId,
        chat_id: chat.chat_id, 
        folder_name: folderName,
        is_checked: !isChecked
      }
    }));
  };

  const handleMouseLeave = () => {
    closeTimeoutRef.current = setTimeout(() => {
    onClose();
    }, 150);
  };

  const handleMouseEnter = () => {
    if (closeTimeoutRef.current) {
      clearTimeout(closeTimeoutRef.current);
      closeTimeoutRef.current = null;
    }
  };

  useEffect(() => {
    return () => {
      if (closeTimeoutRef.current) {
        clearTimeout(closeTimeoutRef.current);
      }
    };
  }, []);

  useEffect(() => {
    if (menuRef.current) {
      const rect = menuRef.current.getBoundingClientRect();
      const windowWidth = window.innerWidth;
      const windowHeight = window.innerHeight;
      let newLeft = x;
      let newTop = y;

      if (rect.right > windowWidth) {
        newLeft = windowWidth - rect.width - 10;
      }
      if (rect.bottom > windowHeight) {
        newTop = windowHeight - rect.height - 10;
      }
      if (rect.left < 0) {
        newLeft = 10;
      }
      if (rect.top < 0) {
        newTop = 10;
      }

      if (newLeft !== position.left || newTop !== position.top) {
        setPosition({ top: newTop, left: newLeft });
      }
    }
  }, [x, y, position.left, position.top]);

  const menuStyle = {
    position: 'fixed',
    top: position.top,
    left: position.left,
    zIndex: 9999
  };

  return (
    <div
      ref={menuRef}
      className="context-menu"
      style={menuStyle}
      onMouseLeave={handleMouseLeave}
      onMouseEnter={handleMouseEnter}
    >
      <div className="menu-item" onClick={handleRename}>
        <span>Переименовать</span>
      </div>

      {Object.keys(folders).filter(f => f !== "All").length > 0 && (
        <>
          <div className="menu-divider"></div>
          <div className="menu-sub-header">
            <span>Папки</span>
          </div>
          {Object.keys(folders).filter(f => f !== "All").map(folder => (
          <div
            key={folder}
              className="menu-item sub checkbox-item"
              onClick={(e) => {
                e.stopPropagation();
                handleToggleFolder(folder);
              }}
            >
              <input
                type="checkbox"
                checked={selectedFolders.has(folder)}
                onChange={() => {}} // Обработка через onClick родителя
                onClick={(e) => e.stopPropagation()}
              />
              <span>{folder}</span>
          </div>
        ))}
        </>
      )}
    </div>
  );
};

export default ChatContextMenu;
