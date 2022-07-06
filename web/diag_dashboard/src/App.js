import { Provider } from 'react-redux';
import {HashRouter,Routes,Route} from "react-router-dom"; 

import {store} from './redux';
import Dashboard from './pages/Dashboard';

function App() {
  return (
    <Provider store={store}>
      <HashRouter>
        <Routes>
          <Route path="/" exact={true} element={<Dashboard/>} />
        </Routes>
      </HashRouter>
    </Provider>
  );
}

export default App;
