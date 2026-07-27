const SelectableUsersList = ({
    users,
    selectedUserIds,
    toggleUser,
}) => {

    return (
      <div className="create-dialog-members">

        <div className="members-header">
          <span>People</span>
          <small>{selectedUserIds.length} selected</small>
        </div>

        <div className="members-list">

          {users.map(user => {

            const selected =
              selectedUserIds.includes(user.id);

            return (

              <button
                key={user.id}
                type="button"
                className={`member-row ${selected ? "selected" : ""}`}
                onClick={() => toggleUser(user.id)}
              >

                <span className="member-avatar">
                  {user.avatar
                    ? <img src={user.avatar} alt="" />
                    : user.name[0].toUpperCase()}
                </span>

                <span className="member-info">
                  <strong>{user.name}</strong>
                  <small>{user.subtitle}</small>
                </span>

                <span className="member-check">
                  {selected ? "Added" : "Add"}
                </span>

              </button>

            );

          })}

        </div>

      </div>
    );
};

export default SelectableUsersList;