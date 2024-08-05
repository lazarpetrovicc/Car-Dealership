import { render, screen } from '@testing-library/react';
import { BrowserRouter as Router } from 'react-router-dom';
import Home from './Home';

describe('Home Page', () => {
  test('renders the main heading and introductory paragraph', () => {
    render(
      <Router>
        <Home />
      </Router>
    );

    // Check if the main heading is rendered
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Welcome to the Car Dealership App');

    // Check if the introductory paragraph is rendered
    expect(screen.getByText(/Manage your car inventory efficiently with our easy-to-use application/i)).toBeInTheDocument();
  });

  test('renders app features list', () => {
    render(
      <Router>
        <Home />
      </Router>
    );

    // List of features to check
    const features = [
      'Manage reservations and deletions of cars.',
      'Add new cars to your inventory and update existing details.',
      'View detailed information including car make, model, price, and status.',
      'Reserve or sell cars to customers with ease.',
      'Cancel reservations and delete available cars as needed.'
    ];

    features.forEach(feature => {
      expect(screen.getByText((content, element) => element.tagName.toLowerCase() === 'li' && content.includes(feature))).toBeInTheDocument();
    });
  });
});