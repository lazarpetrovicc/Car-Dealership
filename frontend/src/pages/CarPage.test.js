import React from 'react';
import { render, screen } from '@testing-library/react';
import CarPage from './CarPage';

describe('CarPage', () => {
  test('renders CarPage without crashing', () => {
    render(<CarPage />);

    // Check if the tabs are rendered
    expect(screen.getByText(/Available/i)).toBeInTheDocument();
    expect(screen.getByText(/Reserved/i)).toBeInTheDocument();
    expect(screen.getByText(/Sold/i)).toBeInTheDocument();
    
    // Check if the button to add a car is rendered
    expect(screen.getByText(/Add Car/i)).toBeInTheDocument();

    // Check if there is an error message element
    expect(screen.queryByText(/Failed to fetch cars/i)).toBeNull(); // Initially not visible
  });
});