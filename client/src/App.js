// src/App.js
import React, { useState } from 'react';
import Navigation from "./components/navigation/Navigation";

// const App = () => {
//     const [token, setToken] = useState('');
//     const [userId, setUserID] = useState('');
//
//     const handleLogin = (token, userId) => {
//         setUserID(userId);
//         setToken(token);
//     };
//
//     const handleLogout = () => {
//         setToken('');
//         setUserID('');
//         localStorage.removeItem('token'); // Очистка токена из локального хранилища
//     };
//
//     return (
//         <div className="App">
//             <Navigation token={token} />
//             <Routes>
//                 <Route path="/" component={About} />
//                 <Route path="/auth" component={Auth} />
//                 <Route path="/dialogs" component={Dialogs} />
//             </Routes>
//         </div>
//     );
// };

function App() {

    return (
        <div>
            {/* Меню навигации */}
            <Navigation />
            {/* Контент страниц */}
        </div>
    );
}

export default App;
