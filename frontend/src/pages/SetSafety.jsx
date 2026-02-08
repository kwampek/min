import { Link } from 'react-router-dom';

import Safety from './Safety';
import '../assets/styles/settings.css';

const SetSafety = () => {
    return (
        <div className="settings-box">
            <Safety />
            <div className="content">
                <h1>Безопасность и приватность</h1>
                <div className="section">
                    <Link to="settings/blacklist.html">Черный список</Link>
                    <Link to="settings/privacy.html">Конфиденциальность</Link>
                </div>
            </div>
        </div>
    );
};

export default SetSafety;
