import {HashRouter,Routes,Route} from "react-router-dom";
import DiagReport from './pages/DiagReport';

import './App.css';

function App() {
  return (
    <div className="App">
      <HashRouter>
        <Routes>
          <Route path="/" exact={true} element={<DiagReport/>} />
        </Routes>
      </HashRouter>
    </div>
  );
}

export default App;
