import { useState, useEffect, useCallback, useRef } from 'react';

const useWebSocket = (token, url, options = {}) => {
    const [isConnected, setIsConnected] = useState(false);
    const socketRef = useRef(null);
    const reconnectTimeoutRef = useRef();
    const optionsRef = useRef(options);

    // Обновляем ref с опциями при каждом изменении
    optionsRef.current = options;

    const send = useCallback((data) => {
        if (socketRef.current?.readyState === WebSocket.OPEN) {
            socketRef.current.send(data);
        }
    }, []);

    const connect = useCallback(() => {
        socketRef.current = new WebSocket(url);

        socketRef.current.onopen = (event) => {
            setIsConnected(true);
            if (typeof optionsRef.current.onOpen === 'function') {
                optionsRef.current.onOpen(event);
            }
            socketRef.current.send(JSON.stringify({
                action: 'auth',
                data:JSON.stringify({ token: token?.token }),
            }));
        }

        socketRef.current.onclose = (event) => {
            setIsConnected(false);
            if (typeof optionsRef.current.onClose === 'function') {
                optionsRef.current.onClose(event);
            }

            if (optionsRef.current.reconnect) {
                reconnectTimeoutRef.current = setTimeout(
                    connect,
                    optionsRef.current.reconnectInterval || 3000
                );
            }
        };

        socketRef.current.onerror = (event) => {
            if (typeof optionsRef.current.onError === 'function') {
                optionsRef.current.onError(event);
            }
        };

        socketRef.current.onmessage = (event) => {
            if (typeof optionsRef.current.onMessage === 'function') {
                try {
                    const data = JSON.parse(event.data);
                    optionsRef.current.onMessage(data);
                } catch (error) {
                    console.error('Error parsing WebSocket message:', error);
                }
            }
        };
    }, [url]);

    const disconnect = useCallback(() => {
        if (socketRef.current) {
            socketRef.current.close();
        }
        clearTimeout(reconnectTimeoutRef.current);
    }, []);

    useEffect(() => {
        connect();

        return () => {
            disconnect();
        };
    }, [connect, disconnect]);

    return {
        isConnected,
        send,
        disconnect,
        connect
    };
};

export default useWebSocket;