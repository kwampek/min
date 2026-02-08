import { useState } from 'react';
import { setUserInAuth } from "../store/slices/meSlice";
import { useNavigate } from 'react-router-dom';
import '../assets/styles/login_style.css'
import { useDispatch } from 'react-redux';


function Login() {
    const dispatch = useDispatch();

    const [error, setError] = useState("");
    const [login, setLogin] = useState("");
    const [password, setPassword] = useState("");

    const rerouting = useNavigate();

    const handleButton = async (type) => {
        const payload = {
            login: login,
            password: password
        }

        try {
            const res = await fetch("/api/" + type, {
                method: "POST",
                headers: {
                "Content-Type": "application/json",
                },
                credentials: "include", // для cookie
                body: JSON.stringify(payload),
            });

            if (!res.ok) {
                const text = await res.text();
                setError(text || "Ошибка при " + type);
                return;
            }

            const data = await res.json();

            dispatch(setUserInAuth(data))
            rerouting("/");


        } catch (err) {
            console.error(err);
            setError("Все сдохло (");
        }
    };



    return (
        <div className="login-box"> 
            <div className="text_min">
                <h1>MIN</h1>
            </div>
            
            <div className="input_grid">
                <input type="text" className="input_field" value={login} onChange={(e) => setLogin(e.target.value)} placeholder="Логин"></input>
            </div>

            <div className="input_grid">
                <input type="password" className="input_field" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="Пароль"></input>
            </div>

            {error && <div style={{ color: "red" }}>{error}</div>}

            <div className="input_grid">
                <button className="login_button" onClick={() => handleButton("login")}  >Войти</button>
            </div>

            <div className="input_grid">
                <button className="registration_button" onClick={() => handleButton("register")} >Зарегистрироваться</button>
            </div>
        </div>
    );
}

export default Login;
