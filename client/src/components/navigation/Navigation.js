import {Link, Route, Routes, useNavigate} from 'react-router-dom';
import React, {useState} from "react";
import Auth from "../pages/auth/Auth";
import Dialogs from "../pages/dialogs/Dialogs";
import About from "../pages/about/About";
import Posts from "../pages/posts/Posts";

const Navigation = ({}) => {
    const [token, setToken] = useState('');
    const [userId, setUserID] = useState('');
    const navigate = useNavigate();

    const handleLogin = (token, userId) => {
        setUserID(userId);
        setToken(token);
        navigate("/");
    };

    const handleLogout = () => {
        setToken('');
        setUserID('');
        localStorage.removeItem('token'); // Очистка токена из локального хранилища
    };

    return (
        <>
            <nav>
                <ul style={{ listStyle: 'none', display: 'flex', gap: '20px' }}>
                    {token ? (
                        <>
                            <li>
                                <Link to="/">Главная</Link>
                            </li>
                            <li>
                                <Link to="/dialogs">Диалоги</Link>
                            </li>
                            <li>
                                <Link to="/posts">Посты</Link>
                            </li>
                            <li>
                                <Link onClick={handleLogout} to="/logout">Выйти</Link>
                            </li>
                        </>
                    ):(
                        <>
                            <li>
                                <Link to="/auth">Авторизация</Link>
                            </li>
                        </>
                    )}
                </ul>
            </nav>
            <Routes>
                <Route path="/" element={<About />} />
                <Route path="/auth" element={<Auth onLogin={handleLogin} />} />
                <Route path="/dialogs" element={<Dialogs token={token} userId={userId} />} />
                <Route path="/posts" element={<Posts token={token} userId={userId} />} />
            </Routes>
        </>

    );
};

export default Navigation;