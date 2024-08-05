import React, { useState } from 'react';
import axios from 'axios';
import '../styles.css';
import actions from '../constants/actions';

const ConfirmActionModal = ({ isOpen, onClose, action, car, onActionComplete }) => {
  const apiUrl = process.env.REACT_APP_API_URL;
  const [error, setError] = useState(null);

  const handleSubmit = () => {
    setError(null); // Clear any previous errors
    if (action === actions.deleteAction) {
      deleteCar(car);
    } else if (action === actions.cancelAction) {
      cancelReservation(car);
    }
  };

  const deleteCar = (car) => {
    axios.delete(`${apiUrl}/cars/${car.id}`)
      .then(response => {
        console.log('Car deleted:', response.data);
        onActionComplete(); // Callback to refresh or update the car list
        onClose();
      })
      .catch(error => {
        console.error('Error deleting car:', error);
        setError('Failed to delete car. Please try again.');
      });
  };

  const cancelReservation = (car) => {
    axios.post(`${apiUrl}/cars/${car.id}/cancel-reservation`)
      .then(response => {
        console.log('Reservation canceled:', response.data);
        onActionComplete(); // Callback to refresh or update the car list
        onClose();
      })
      .catch(error => {
        console.error('Error canceling reservation:', error);
        setError('Failed to cancel reservation. Please try again.');
      });
  };

  const title = action === actions.deleteAction ? 'Confirm Car Deletion' : 'Confirm Cancellation';
  const message = action === actions.deleteAction ? `Are you sure you want to delete ${car?.make} ${car?.model}?` : `Are you sure you want to cancel the reservation for ${car?.make} ${car?.model} by ${car?.customer?.fullName}?`;

  return (
    <>
      {isOpen && (
        <div className="modal" role="dialog" aria-modal="true" aria-labelledby="modal-title">
          <div className="modal-content">
            <span className="close" onClick={onClose}>&times;</span>
            <h2 id="modal-title" className="modal-title">{title}</h2>
            <p className="modal-message">{message}</p>
            {error && <p className="error-message">{error}</p>}
            <div className="modal-buttons">
              <button className="modal-button" onClick={handleSubmit}>Confirm</button>
              <button className="modal-button" onClick={onClose}>Cancel</button>
            </div>
          </div>
        </div>
      )}
    </>
  );
};

export default ConfirmActionModal;