import { render, screen } from '@testing-library/react';
import About from './About';

describe('About Page', () => {
  test('renders the main heading and description', () => {
    render(<About />);

    // Check if the main heading is rendered
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('About This Application');

    // Check if the description paragraph is rendered
    expect(screen.getByText(/This application is designed to manage the inventory of a car dealership/i)).toBeInTheDocument();
  });

  test('renders all functionality sections', () => {
    render(<About />);

    // Check if the "Functionalities" section heading is rendered
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Functionalities');

    // Check for each functionality item
    const functionalityItems = [
      {
        title: '1. View Cars',
        description: /Navigate to the "Car Management" page to view a list of cars categorized into available, reserved, and sold/i
      },
      {
        title: '2. Add and Edit Cars',
        description: /You can add new cars to the inventory or edit existing car details/i
      },
      {
        title: '3. Reserve and Sell Cars',
        description: /For available cars, you can reserve or sell them to customers/i
      },
      {
        title: '4. Cancel Reservations and Delete Cars',
        description: /You can cancel reservations for cars that are currently reserved/i
      }
    ];

    functionalityItems.forEach(item => {
      expect(screen.getByRole('heading', { level: 3, name: item.title })).toBeInTheDocument();
      expect(screen.getByText((content, element) => element.tagName.toLowerCase() === 'p' && item.description.test(content))).toBeInTheDocument();
    });
  });
});