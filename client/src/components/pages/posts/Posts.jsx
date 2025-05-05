import React, {useEffect, useState} from 'react';
import useWebSocket from "../../hook/WebSocket";

const Posts = (token) => {
    const [posts, setPosts] = useState([]);
    const [newPost, setNewPost] = useState({ user_id: '', text: '' });

    useEffect(() => {

    }, []);

    const { isConnected, send } = useWebSocket(token, 'ws://localhost:8080/ws', {
        onMessage: (data) => {
            if (Array.isArray(data.Post)) {
                // Если сервер прислал массив постов
                setPosts(data);
            } else if (data.action === 'posts/create') {
                // Если сервер прислал один новый пост
                setPosts(prev => [...prev, data.Post]);
            }
        },
        onOpen: () => {
            send(JSON.stringify({
                action: 'posts/list',
                data: JSON.stringify({
                    token: token,
                })
            }));
        },
        reconnect: true
    });

    const handleSubmit = (e) => {
        e.preventDefault();

        send(JSON.stringify({
            action: 'posts/create',
            data: JSON.stringify({
                text: newPost.text,
                user_id: newPost.user_id,
            })
        }));

        setNewPost({ user_id: '', text: '' });
    };

    return (
        <div className="container">
            <h2>Посты {isConnected ? '✓' : '⌛'}</h2>

            <h5>Сообщение отправлено</h5>
            <ul className="collection">
                {posts && posts.map((post) => (
                    <li key={post.id} className="collection-item">
                        <strong>ToUserId:</strong> {post?.user_id}
                        <strong>Text:</strong> {post?.text}
                    </li>
                ))}
            </ul>

            <div className="row">
                <form onSubmit={handleSubmit} className="col s12">
                    <div className="row">
                        <input placeholder="Text" type="text" className="validate"
                               onChange={(e) => setNewPost({...newPost, text: e.target.value})}/>
                    </div>
                    <div className="row">
                        <input placeholder="UserId" type="text" className="validate"
                               onChange={(e) => setNewPost({...newPost, user_id: e.target.value})}/>
                    </div>
                    <button className="btn waves-effect waves-light" disabled={!isConnected}>
                        Send
                    </button>
                </form>
            </div>
        </div>
    );
};

export default Posts;
