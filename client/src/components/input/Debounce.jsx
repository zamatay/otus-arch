import React, { useState, useEffect } from 'react';

const InputWithDebounce = ({ onValueChange }) => {
    const [inputValue, setInputValue] = useState('');
    const [debouncedValue, setDebouncedValue] = useState(inputValue);

    const handleInputChange = (event) => {
        setInputValue(event.target.value);
    };

    useEffect(() => {
        const handler = setTimeout(() => {
            setDebouncedValue(inputValue);
            onValueChange(inputValue);
        }, 500);

        return () => {
            clearTimeout(handler);
        };
    }, [inputValue]);

    return (
            <input
                type="text"
                id="input"
                value={inputValue}
                onChange={handleInputChange}
            />
    );
};

export default InputWithDebounce;