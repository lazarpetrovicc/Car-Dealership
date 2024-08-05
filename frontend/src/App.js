import React from 'react';
import { BrowserRouter as Router, Route, Routes, Link } from 'react-router-dom';
import './styles.css';
import Home from './pages/Home';
import CarPage from './pages/CarPage';
import About from './pages/About';
import NotFound from './pages/NotFound';

function App() {
  return (
    <Router>
      <div className="App">
        <header className="App-header">
          <div className="left-corner">
            <h1>Car Dealership App</h1>
          </div>
          <nav className="right-corner">
            <ul>
              <li><Link to="/" aria-label="Home">Home</Link></li>
              <li><Link to="/cars" aria-label="Car Management">Car Management</Link></li>
              <li><Link to="/about" aria-label="About">About</Link></li>
            </ul>
          </nav>
        </header>
        <main>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/cars" element={<CarPage />} />
            <Route path="/about" element={<About />} />
            {/* Route for handling 404 Not Found */}
            <Route path="*" element={<NotFound />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;