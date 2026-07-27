const MembersList = ({ users, removeMember }) => {

    return (
        <div className="create-dialog-members">

            <div className="members-header">
                <span>Participants</span>
                <small>{users.length}</small>
            </div>

            <div className="members-list">

                {users.map(user => (

                    <button
                        key={user.id}
                        type="button"
                        className="member-row"
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

                        <button
                            className="remove-member-btn"
                            type="button"
                            onClick={() => removeMember(user.id)}
                        >
                            Remove
                        </button>

                    </button>

                ))}

            </div>

        </div>
    );
};

export default MembersList;