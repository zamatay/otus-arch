// src/Dialogs.js
import React, { useEffect, useState } from 'react';
import jwt_decode from 'jwt-decode';

const Dialogs = ({ token }) => {
    const [dialogs, setDialogs] = useState([]);
    const [newMessage, setNewMessage] = useState('');
    const [ws, setWs] = useState(null);

    useEffect(() => {
        const websocket = new WebSocket('ws://your-websocket-url');

        websocket.onopen = () => {
            console.log('WebSocket connected');
            websocket.send(JSON.stringify({ action: 'authenticate', token }));
            websocket.send(JSON.stringify({ action: 'dialogs/list' }));
        };

        websocket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            if (data.action === 'dialogs/list') {
                setDialogs(data.dialogs);
            } else if (data.action === 'dialogs/post') {
                setDialogs((prevDialogs) => [...prevDialogs, data.dialog]);
            }
        };

        setWs(websocket);

        return () => {
            websocket.close();
        };
    }, [token]);

    const sendMessage = () => {
        if (ws && newMessage.trim()) {
            // const messageData = {
            //     action: 'dialogs/post',
            //     text: newMessage,
            //     from_user_id: jwt_decode(token).user_id,
            //     to_user_id: 'recipient_user_id', // Замените на реальный ID получателя
            // };
            // ws.send(JSON.stringify(messageData));
            // setNewMessage('');
        }
    };

    return (
        <div className="container">
            <h2>Dialogs</h2>
            <ul className="collection">
                {dialogs.map((dialog) => (
                    <li key={dialog.id} className="collection-item">
                        <strong>From:</strong> {dialog.from_user_id} <br />
                        <strong>To:</strong> {dialog.to_user_id} <br />
                        <strong>Text:</strong> {dialog.text}
                    </li>
                ))}
            </ul>
            <div className="input-field">
                <input
                    type="text"
                    id="newMessage"
                    value={newMessage}
                    onChange={(e) => setNewMessage(e.target.value)}
                    required
                />
                <label htmlFor="newMessage">Type a message</label>
            </div>
            <button className="btn waves-effect waves-light" onClick={sendMessage}>
                Send
            </button>
        </div>
    );
};

export default Dialogs;
