import { useMemo, useState } from 'react';
import { useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { getWebSocket } from '../../modules/ws';
import CreateDialog from './CreateDialog';

const LeftPanel = ({ closePanel }) => {
  const navigate = useNavigate();
  const chats = useSelector((state) => state.chats.chats);
  const [createMode, setCreateMode] = useState(null);

  const interactedUsers = useMemo(() => {
    return Object.values(chats).map((chat) => ({
      id: chat.user_id || chat.userId || chat.chat_id,
      name: chat.title || 'Unnamed',
      avatar: typeof chat.avatar === 'string' ? chat.avatar : chat.avatar?.String,
      subtitle: chat.lastMessage || 'Existing conversation',
    }));
  }, [chats]);

  const sendWsAction = (type, payload) => {
    const ws = getWebSocket();

    if (!ws || ws.readyState !== WebSocket.OPEN) {
      console.warn('WebSocket is not connected');
      return;
    }

    ws.send(JSON.stringify({ type, payload }));
  };

  const handleProfileClick = () => {
    navigate('/settings/profile');
    closePanel();
  };

  const handleSettingsClick = () => {
    navigate('/settings/safety');
    closePanel();
  };

  const handleCreateChat = () => {
    setCreateMode('chat');
  };

  const handleCreateChannel = () => {
    setCreateMode('channel');
  };

  const handleCreateSubmit = ({ name, description, photo, selectedUsers }) => {
    const selectedUserIds = selectedUsers.map((user) => user.id);

    if (createMode === 'chat') {
      sendWsAction('create_chat', {
        userId: selectedUserIds.length === 1 ? selectedUserIds[0] : undefined,
        title: name,
        description,
        avatar: photo,
        member_ids: selectedUserIds,
        members: selectedUsers,
      });
    }

    if (createMode === 'channel') {
      sendWsAction('create_channel', {
        title: name,
        description,
        avatar: photo,
        member_ids: selectedUserIds,
        members: selectedUsers,
      });
    }

    setCreateMode(null);
    closePanel();
  };

  return (
    <>
      <div className="left-panel-overlay" onClick={closePanel}>
        <aside className="left-panel" onClick={(event) => event.stopPropagation()}>
          <button type="button" className="left-panel-item" onClick={handleProfileClick}>
            <span className="left-panel-icon">P</span>
            <span>Profile</span>
          </button>

          <button type="button" className="left-panel-item" onClick={handleCreateChat}>
            <span className="left-panel-icon">+</span>
            <span>Add Chat</span>
          </button>

          <button type="button" className="left-panel-item" onClick={handleCreateChannel}>
            <span className="left-panel-icon">#</span>
            <span>Add Channel</span>
          </button>

          <button type="button" className="left-panel-item" onClick={handleSettingsClick}>
            <span className="left-panel-icon">S</span>
            <span>Settings</span>
          </button>
        </aside>
      </div>

      {createMode && (
        <CreateDialog
          mode={createMode}
          users={interactedUsers}
          onClose={() => setCreateMode(null)}
          onSubmit={handleCreateSubmit}
        />
      )}
    </>
  );
};

export default LeftPanel;
