import React from 'react';
import Dialogs from "./components/dialogs/Dialogs";
import 'materialize-css/dist/css/materialize.min.css';

const App = () => {
  return (
      <div className="App">
        <h1>WebSocket Dialogs</h1>
        <Dialogs />
      </div>
  );
};

export default App;