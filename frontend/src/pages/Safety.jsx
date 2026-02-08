import { Link } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import '../assets/styles/settings.css';
import { useNavigate } from 'react-router-dom';

const Settings = () => {
    const dispatch = useDispatch();
    const navigate = useNavigate();

    const handleLogout = () => {
        localStorage.removeItem('userData');
        dispatch({ type: 'LOGOUT' });
        navigate('/login', { replace: true });
    };


    return (
        <div className="sidebar">
            <Link to="/" className="back-button">← Назад</Link>
            <div className="settings-profile">
                <img src='/icons/me.jpg' />
                <div className="username">
                    KWAMPEK <br/>
                    <small>@kwampek</small>
                </div>
            </div>
            <div className="search-box">
                <input type="text" placeholder="Поиск" />
            </div>
            <div className="menu">
                <Link to="/settings/profile">Профиль</Link>
                <Link to="/settings/gosuslugi">Госуслуги</Link>
                <button className="logout-button" onClick={handleLogout}>
                Выйти
                </button>
            </div>
        </div>
    );
};

export default Settings;
