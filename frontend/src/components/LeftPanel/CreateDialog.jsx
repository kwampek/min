import { useEffect, useMemo, useState } from 'react';

const CreateDialog = ({ mode, users, onClose, onSubmit }) => {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [photo, setPhoto] = useState(null);
  const [selectedUserIds, setSelectedUserIds] = useState([]);
  const reader = new FileReader();

  const isChat = mode === 'chat';
  const title = isChat ? 'Create Chat' : 'Create Channel';

  const photoPreview = useMemo(() => {
    if (!photo) return null;
    return URL.createObjectURL(photo);
  }, [photo]);

  useEffect(() => {
    return () => {
      if (photoPreview) {
        URL.revokeObjectURL(photoPreview);
      }
    };
  }, [photoPreview]);

  const toggleUser = (userId) => {
    setSelectedUserIds((currentIds) => {
      if (currentIds.includes(userId)) {
        return currentIds.filter((id) => id !== userId);
      }

      return [...currentIds, userId];
    });
  };

  const handleSubmit = (event) => {
    event.preventDefault();

    const selectedUsers = users.filter((user) => selectedUserIds.includes(user.id));
    onSubmit({
      name: name.trim(),
      description: description.trim(),
      photo: photo
        ? {
            name: photo.name,
            type: photo.type,
            size: photo.size,
            data: reader.result.split(",")[1],
          }
        : null,
      selectedUsers,
    });
  };

  return (
    <div className="create-dialog-overlay" onClick={onClose}>
      <form className="create-dialog" onClick={(event) => event.stopPropagation()} onSubmit={handleSubmit}>
        <div className="create-dialog-header">
          <div>
            <h2>{title}</h2>
            <p>{isChat ? 'Choose chat details and members.' : 'Choose channel details and members.'}</p>
          </div>
          <button type="button" className="create-dialog-close" onClick={onClose} aria-label="Close">
            x
          </button>
        </div>

        <label className="photo-picker">
          <input
            type="file"
            accept="image/*"
            onChange={(event) => setPhoto(event.target.files?.[0] || null)}
          />
          <span className="photo-preview">
            {photoPreview ? <img src={photoPreview} alt="" /> : name.trim().charAt(0).toUpperCase() || '+'}
          </span>
          <span>
            <strong>Photo</strong>
            <small>{photo ? photo.name : 'Choose image'}</small>
          </span>
        </label>

        <div className="create-dialog-fields">
          <label>
            <span>Name</span>
            <input
              type="text"
              value={name}
              onChange={(event) => setName(event.target.value)}
              placeholder={isChat ? 'Design team' : 'Product updates'}
              maxLength={64}
              required
            />
          </label>

          <label>
            <span>Description</span>
            <textarea
              value={description}
              onChange={(event) => setDescription(event.target.value)}
              placeholder="Short description"
              rows={3}
              maxLength={240}
            />
          </label>
        </div>

        <div className="create-dialog-members">
          <div className="members-header">
            <span>People</span>
            <small>{selectedUserIds.length} selected</small>
          </div>

          <div className="members-list">
            {users.length === 0 ? (
              <div className="members-empty">No interacted users yet</div>
            ) : (
              users.map((user) => {
                const isSelected = selectedUserIds.includes(user.id);

                return (
                  <button
                    type="button"
                    key={user.id}
                    className={`member-row ${isSelected ? 'selected' : ''}`}
                    onClick={() => toggleUser(user.id)}
                  >
                    <span className="member-avatar">
                      {user.avatar ? <img src={user.avatar} alt="" /> : user.name.charAt(0).toUpperCase()}
                    </span>
                    <span className="member-info">
                      <strong>{user.name}</strong>
                      <small>{user.subtitle}</small>
                    </span>
                    <span className="member-check">{isSelected ? 'Added' : 'Add'}</span>
                  </button>
                );
              })
            )}
          </div>
        </div>

        <div className="create-dialog-actions">
          <button type="button" className="create-dialog-secondary" onClick={onClose}>
            Cancel
          </button>
          <button type="submit" className="create-dialog-primary">
            Create
          </button>
        </div>
      </form>
    </div>
  );
};

export default CreateDialog;
