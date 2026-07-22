import { useNavigate } from 'react-router-dom';

const SettingsDrawer = ({ setCreateMode }) => {
    const navigate = useNavigate();

    return (
        <div className="settings-drawer">
            <div onClick={() => navigate("/settings/profile")}>
                Профиль
            </div>

            <div onClick={() => {setCreateMode("chat")}}>
                Создать чат
            </div>

            <div onClick={() => {setCreateMode("channel")}}>
                Создать канал
            </div>

            <div onClick={() => navigate("/settings/safety")}>
                Безопасность
            </div>

            <div onClick={() => navigate("/settings/devices")}>
                Устройства
            </div>
        </div>
    );
};


export default SettingsDrawer;
