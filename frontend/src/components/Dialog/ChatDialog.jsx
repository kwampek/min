import { useEffect, useMemo, useState } from "react";
import { useSelector } from "react-redux";

import MembersList from "./MembersList";
import SelectableUsersList from "./SelectableUsersList";

const ChatDialog = ({
    mode,
    type,
    onClose,
    onSubmit,
}) => {

  const {
      loading,
      chat,
      users,
  } = useSelector(state => state.chatDialog);

  const isChat = type === "chat";
  const isEdit = mode === "edit";

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [photo, setPhoto] = useState(null);
  const [selectedUserIds, setSelectedUserIds] = useState([]);

  useEffect(() => {
    setName(chat?.title ?? "");
    setDescription(chat?.description ?? "");
    setPhoto(null);

    setSelectedUserIds(
      chat?.members?.map(u => u.id) ?? []
    );
  }, [chat]);

  const photoPreview = useMemo(() => {
    if (!photo)
      return chat?.avatar ?? null;

    return URL.createObjectURL(photo);
  }, [photo, chat]);

  useEffect(() => {
      return () => {
          if (photo && photoPreview)
              URL.revokeObjectURL(photoPreview);
      };
  }, [photo, photoPreview]);

  const dialogTitle =
    isEdit
      ? (isChat ? "Edit Chat" : "Edit Channel")
      : (isChat ? "Create Chat" : "Create Channel");

  const dialogSubtitle =
    isEdit
      ? "Edit chat information."
      : isChat
        ? "Choose chat details and members."
        : "Choose channel details and members.";

  const removeMember = (user_id) => {
    setSelectedUserIds(current =>
      current.includes(user_id)
        ? current.filter(x => x !== user_id)
        : [...current, user_id]
    );
  }

  const toggleUser = (user_id) => {
    setSelectedUserIds(current =>
      current.includes(user_id)
        ? current.filter(x => x !== user_id)
        : [...current, user_id]
    );
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    onSubmit({
      chatId: chat?.chat_id,
      name,
      description,
      photo,
      selectedUserIds,
    });
  };

  return (
      <div className="create-dialog-overlay" onClick={onClose}>
          <form
              className="create-dialog"
              onClick={e => e.stopPropagation()}
              onSubmit={handleSubmit}
          >

              <div className="create-dialog-header">
                  <div>
                      <h2>{dialogTitle}</h2>
                      <p>{dialogSubtitle}</p>
                  </div>

                  <button
                      type="button"
                      className="create-dialog-close"
                      onClick={onClose}
                  >
                      ✕
                  </button>
              </div>

              <label className="photo-picker">
                  <input
                      type="file"
                      accept="image/*"
                      onChange={(e) =>
                          setPhoto(e.target.files?.[0] ?? null)
                      }
                  />

                  <span className="photo-preview">
                      {photoPreview
                          ? <img src={photoPreview} alt="" />
                          : (name[0]?.toUpperCase() ?? "+")}
                  </span>
              </label>

              <div className="create-dialog-fields">
                  <label>
                      <span>Name</span>

                      <input
                          value={name}
                          onChange={(e) => setName(e.target.value)}
                          required
                      />
                  </label>

                  <label>
                      <span>Description</span>

                      <textarea
                          rows={3}
                          value={description}
                          onChange={(e) => setDescription(e.target.value)}
                      />
                  </label>

              </div>

              {loading.users ? (
                  <div className="members-loading">
                      Loading...
                  </div>
              ) : isEdit ? (
                  <MembersList
                      users={users}
                      removeMember={removeMember}
                  />
              ) : (
                  <SelectableUsersList
                      users={users}
                      selectedUserIds={selectedUserIds}
                      toggleUser={toggleUser}
                  />
              )}

              <div className="create-dialog-actions">

                  <button
                      type="button"
                      className="create-dialog-secondary"
                      onClick={onClose}
                  >
                      Cancel
                  </button>

                  <button
                      type="submit"
                      className="create-dialog-primary"
                  >
                      {isEdit ? "Save" : "Create"}
                  </button>

              </div>

          </form>
      </div>
  );
};

export default ChatDialog;