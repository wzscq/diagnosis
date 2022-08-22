import {HashRouter,Routes,Route} from "react-router-dom";
import DiagEventReport from './pages/DiagEventReport';

import './App.css';

function App() {
  return (
    <div className="App">
      <HashRouter>
        <Routes>
          <Route path="/" exact={true} element={<DiagEventReport/>} />
        </Routes>
      </HashRouter>
    </div>
  );
}

export default App;
