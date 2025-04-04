// src/App.js
import React, { useState } from 'react';
import Auth from '././components/auth/Auth';
import Dialogs from '././components/dialogs/Dialogs';

const App = () => {
    const [token, setToken] = useState('');

    const handleLogin = (token) => {
        setToken(token);
    };

    const handleLogout = () => {
        setToken('');
        localStorage.removeItem('token'); // Очистка токена из локального хранилища
    };

    return (
        <div className="App">
            <nav>
                <div className="nav-wrapper">
                    <a href="#" className="brand-logo">My App</a>
                    <ul id="nav-mobile" className="right hide-on-med-and-down">
                        {token && (
                            <li>
                                <button className="btn red" onClick={handleLogout}>
                                    Logout
                                </button>
                            </li>
                        )}
                    </ul>
                </div>
            </nav>
            <div className="container">
                <h1>WebSocket Dialogs</h1>
                {!token ? (
                    <Auth onLogin={handleLogin} />
                ) : (
                    <Dialogs token={token} />
                )}
            </div>
        </div>
    );
};

export default App;
