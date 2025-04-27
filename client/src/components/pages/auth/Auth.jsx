import React, { useState } from 'react';
import axios from 'axios';

const Auth = ({ onLogin }) => {
    const [login, setLogin] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');

    const handleLogin = async (e) => {
        e.preventDefault();
        try {
            const response = await axios.post('http://localhost:8080/auth/login', {
                login,
                password,
            });
            const token = response.data.token;
            const userId = response.data.user_id;
            localStorage.setItem('token', token);
            localStorage.setItem('user_id', userId);
            onLogin(token, userId);
        } catch (err) {
            setError('Invalid credentials');
        }
    };

    return (
        <div className="container">
            <h2>Login</h2>
            {error && <p className="red-text">{error}</p>}
            <form onSubmit={handleLogin}>
                <div className="input-field">
                    <input
                        type="text"
                        id="username"
                        value={login}
                        onChange={(e) => setLogin(e.target.value)}
                        required
                    />
                    <label htmlFor="username">Username</label>
                </div>
                <div className="input-field">
                    <input
                        type="password"
                        id="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                    />
                    <label htmlFor="password">Password</label>
                </div>
                <button className="btn waves-effect waves-light" type="submit">Login</button>
            </form>
        </div>
    );
};

export default Auth;