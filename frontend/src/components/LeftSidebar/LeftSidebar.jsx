import React, { useState, useRef, useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { setActiveFolder, addFolder, editFolderName } from '../../store/slices/chatSlice';
import { getWebSocket } from '../../modules/ws';

const LeftSidebar = ({ isLeftPanelOpen, setLeftPanelOpen}) => {
  const dispatch = useDispatch();

  const { folders, activeFolderName} = useSelector((state) => state.chats);
  const user_id = useSelector((state) => state.user.userId);
  const folders_list = Object.keys(folders);
  const [editingFolder, setEditingFolder] = useState(null);
  const [editValue, setEditValue] = useState('');
  const inputRef = useRef(null);

  const handleFolderClick = (folderName, e) => {
    if (e && e.target.closest('.folder-name')) {
      return;
    }
    dispatch(setActiveFolder(folderName));
  };

  const handleFolderNameClick = (folderName, e) => {
    e.stopPropagation();
    if (folderName === 'All') return;
    setEditingFolder(folderName);
    setEditValue(folderName);
  };

  const handleFolderNameBlur = () => {
    if (editingFolder) {
      saveFolderName();
    }
  };

  const handleFolderNameKeyDown = (e) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      saveFolderName();
    } else if (e.key === 'Escape') {
      setEditingFolder(null);
      setEditValue('');
    }
  };

  const saveFolderName = () => {
    if (!editingFolder) return;
    
    const trimmedName = editValue.trim();
    
    if (!trimmedName || trimmedName === editingFolder) {
      setEditingFolder(null);
      setEditValue('');
      return;
    }

    if (folders_list.includes(trimmedName)) {
      alert('Папка с таким именем уже существует');
      setEditingFolder(null);
      setEditValue('');
      return;
    }

    const ws = getWebSocket();
    dispatch(editFolderName({ prev_name: editingFolder, new_name: trimmedName }));

    ws.send(JSON.stringify({
      type: "edit_folder_name",
      payload: {
        user_id: user_id,
        prev_name: editingFolder,
        new_name: trimmedName
      }
    }));

    setEditingFolder(null);
    setEditValue('');
  };

  useEffect(() => {
    if (editingFolder && inputRef.current) {
      inputRef.current.focus();
      inputRef.current.select();
    }
  }, [editingFolder]);

  const handleAddFolder = () => {
    const newName = prompt('Введите имя новой папки:');
    if (!newName) return;

    const trimmedName = newName.trim();
    if (!folders_list.includes(trimmedName)) {
      dispatch(addFolder(trimmedName));
      const ws = getWebSocket();

      ws.send(JSON.stringify({
        type: "create_folder",
        payload: {
          user_id: user_id,
          folder_name: trimmedName
        }
      }));
    }
  };
  
  // TODO fix img

  return (
    <div className="left-sidebar">
      <div className="min-label">MIN</div>

      <div className="folders">
        {folders_list.map((folder, index) => (
          <div
            key={index}
            className={`folder ${folder === activeFolderName ? 'active' : ''}`}
            onClick={(e) => handleFolderClick(folder, e)}
            title={folder}
          >
            <img src="icons/folder.png"/>
            {editingFolder === folder ? (
              <input
                ref={inputRef}
                type="text"
                className="folder-name-input"
                value={editValue}
                onChange={(e) => setEditValue(e.target.value)}
                onBlur={handleFolderNameBlur}
                onKeyDown={handleFolderNameKeyDown}
                onClick={(e) => e.stopPropagation()}
                maxLength={20}
              />
            ) : (
              <div 
                className="folder-name"
                onClick={(e) => handleFolderNameClick(folder, e)}
              >
              {folder.length > 8 ? folder.slice(0, 7) + '…' : folder}
            </div>
            )}
          </div>
        ))}
      </div>

      <button className="edit-button" onClick={handleAddFolder} title="Добавить папку">
        <div className="add-button">+</div>
      </button>

      <div className="profile" onClick={() => setLeftPanelOpen(!isLeftPanelOpen)}>
        <img src="icons/me.jpg" />  
      </div>
    </div>
  );
};

export default LeftSidebar;
