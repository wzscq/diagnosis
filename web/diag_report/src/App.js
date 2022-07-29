import {HashRouter,Routes,Route} from "react-router-dom";
import DiagReport from './pages/DiagReport';
import useFrame from './hook/useFrame';

import './App.css';

function App() {
  const sendMessageToParent=useFrame();

  return (
    <div className="App">
      <HashRouter>
        <Routes>
          <Route path="/" exact={true} element={<DiagReport sendMessageToParent={sendMessageToParent}/>} />
        </Routes>
      </HashRouter>
    </div>
  );
}

export default App;
