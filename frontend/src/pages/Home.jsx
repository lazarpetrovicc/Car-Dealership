import React from 'react';
import { Link } from 'react-router-dom';

const Home = () => {
  return (
    <div className="home-page">
      <div className="welcome-content">
        <h1>Welcome to the Car Dealership App</h1>
        <p className="welcome-text">
          Manage your car inventory efficiently with our easy-to-use application.
        </p>
        <div className="action-buttons">
          <Link to="/cars" className="home-button">View Car Inventory</Link>
          <Link to="/about" className="home-button">Learn More About Us</Link>
        </div>
      </div>
      <div className="app-info">
        <h2>App Features</h2>
        <ul className="features-list">
          <li>- Manage reservations and deletions of cars.</li>
          <li>- Add new cars to your inventory and update existing details.</li>
          <li>- View detailed information including car make, model, price, and status.</li>
          <li>- Reserve or sell cars to customers with ease.</li>
          <li>- Cancel reservations and delete available cars as needed.</li>
        </ul>
      </div>
    </div>
  );
};

export default Home;