import { render, screen } from '@testing-library/react';
import App from './App';

describe('App', () => {
  test('renders the header with the correct title', () => {
    render(<App />);

    // Check if the app title is rendered
    expect(screen.getByText('Car Dealership App')).toBeInTheDocument();
  });

  test('renders navigation links', () => {
    render(<App />);

    // Check if navigation links are rendered
    const homeLink = screen.getByLabelText(/Home/i);
    const carManagementLink = screen.getByLabelText(/Car Management/i);
    const aboutLink = screen.getByLabelText(/About/i);
    
    expect(homeLink).toBeInTheDocument();
    expect(carManagementLink).toBeInTheDocument();
    expect(aboutLink).toBeInTheDocument();
  });
});