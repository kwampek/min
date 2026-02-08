import { useState } from 'react';
import Safety from './Safety';
import '../assets/styles/settings.css';

const SetProfile = () => {

    const [profile, _] = useState({
        username: '@KWAMPEK',
        email: 'chesnokov.as@phystech.edu',
        phone: '79050000000',
        birthdate: '2005-11-18',
        gender: 'Мужской',
    });

    const handleSaveProfile = () => {
        console.log(profile);
    };

    return (
        <div className="settings-profile-box">
            <Safety />
            <div className="settings-content">
                <div className="profile-image-box">
                    <img src="/icons/me.jpg" />
                </div>
                <div className="profile-info-box">  
                    <input type="text" placeholder={profile.username} />
                    <input type="text" placeholder={profile.email} />
                    <input type="text" placeholder={profile.phone} />
                    <input type="text" placeholder={profile.birthdate} />
                    <input type="text" placeholder={profile.gender} />
                    
                    <button onClick={handleSaveProfile}>Сохранить</button>
                </div>
            </div>
        </div>
    );
}

export default SetProfile;