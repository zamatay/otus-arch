// src/Dialogs.js
import React, {useEffect, useRef, useState} from 'react';
import { jwtDecode } from 'jwt-decode';
import InputWithDebounce from "../../input/Debounce";

const Dialogs = ({ token, userId }) => {
    const [dialogs, setDialogs] = useState([]);
    const [newMessage, setNewMessage] = useState('');
    const [toUserID, setToUserID] = useState('');
    const [isConnected, setIsConnected] = useState(false);
    const ws = useRef(null);

    const handleInputChange = (value) => {
        setToUserID(value)
        if (!toUserID || !isConnected || !ws.current) {
            return
        }
        ws.current.send(JSON.stringify({ action: 'dialogs/list', data: JSON.stringify({ user_id: toUserID, token}) }));
    };

    useEffect(() => {
        ws.current = new WebSocket('ws://localhost:8080/ws');
        setIsConnected(true);

        ws.current.onopen = () => {
            setIsConnected(true);
            console.log('WebSocket connected');
            //websocket.send(JSON.stringify({ action: 'dialogs/list', data: JSON.stringify({ user_id: toUserID, token}) }));
        };

        ws.current.onmessage = (event) => {
            const data = JSON.parse(event.data);
            if (data?.action === 'dialogs/list') {
                setDialogs(data.Dialogs);
            } else if (data?.action === 'dialogs/post') {
                setDialogs((prevDialogs) => [...prevDialogs, data.Dialogs]);
            }
        };      

        ws.current.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        ws.current.onclose = () => {
            setIsConnected(false);
            console.log('WebSocket disconnected');
        };
        //setWs(websocket);

        return () => {
            ws.current.close();
        };
    }, [token]);

    const sendMessage = () => {
        if (ws && newMessage.trim()) {
            const messageData = {
                action: 'dialogs/post',
                data: JSON.stringify({
                    text: newMessage,
                    from_user_id: jwtDecode(token).id,
                    to_user_id: Number(toUserID),
                })
            };
            ws.current.send(JSON.stringify(messageData));
            setNewMessage('');
        }
    };

    return (
        <div className="container">
            <h2>Dialogs</h2>
            <ul className="collection">
                {dialogs && dialogs.map((dialog) => (
                    <li key={dialog.id} className="collection-item">
                        <strong>From:</strong> {dialog.from_user_id} <br />
                        <strong>To:</strong> {dialog.to_user_id} <br />
                        <strong>Text:</strong> {dialog.text}
                    </li>
                ))}
            </ul>
            <div className="row">
                <input
                    type="text"
                    id="newMessage"
                    value={newMessage}
                    onChange={(e) => setNewMessage(e.target.value)}
                    required
                />
                <label htmlFor="newMessage">Type a message</label>
            </div>
            <div className="row">
                <InputWithDebounce onValueChange={handleInputChange} />
                <label htmlFor="toUserId">user id</label>
            </div>
            <button className="btn waves-effect waves-light" onClick={sendMessage}>
                Send
            </button>
        </div>
    );
};

export default Dialogs;
