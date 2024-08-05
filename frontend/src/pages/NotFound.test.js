import { render, screen } from '@testing-library/react';
import NotFound from './NotFound';

describe('NotFound Page', () => {
  test('renders 404 heading and message', () => {
    render(<NotFound />);

    // Check if the 404 heading is rendered
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('404 - Not Found');

    // Check if the not found message is rendered
    expect(screen.getByText(/The page you are looking for does not exist./i)).toBeInTheDocument();
  });
});