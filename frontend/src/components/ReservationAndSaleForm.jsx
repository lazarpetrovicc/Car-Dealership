import React, { useState } from 'react';
import axios from 'axios';
import '../styles.css';
import actions from '../constants/actions';

const ReservationAndSaleForm = ({ isOpen, onClose, action, car, onActionComplete }) => {
  const apiUrl = process.env.REACT_APP_API_URL;
  const [customerFullName, setCustomerFullName] = useState('');
  const [customerEmail, setCustomerEmail] = useState('');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [error, setError] = useState(null);

  if (!isOpen) return null;

  const handleSubmit = (e) => {
    e.preventDefault();

    const customer = {
      fullname: customerFullName,
      email: customerEmail,
      phonenumber: phoneNumber,
    };

    if (action === actions.reserveAction) {
      reserveCar(car, customer);
    } else if (action === actions.sellAction) {
      sellCar(car, customer);
    }
    
    setCustomerFullName('');
    setCustomerEmail('');
    setPhoneNumber('');
  };

  const reserveCar = (car, customer) => {
    axios.post(`${apiUrl}/cars/${car.id}/reserve`, customer)
      .then(response => {
        console.log('Car reserved:', response.data);
        onActionComplete();
        onClose();
      })
      .catch(error => {
        console.error('Error reserving car:', error);
        setError('Failed to reserve car. Please try again.');
      });
  };

  const sellCar = (car, customer) => {
    axios.post(`${apiUrl}/cars/${car.id}/sell`, customer)
      .then(response => {
        console.log('Car sold:', response.data);
        onActionComplete();
        onClose();
      })
      .catch(error => {
        console.error('Error selling car:', error);
        setError('Failed to sell car. Please try again.');
      });
  };

  const title = action === actions.reserveAction ? 'Confirm Reservation' : 'Confirm Sale';
  const message = action === actions.reserveAction ? `Are you sure you want to reserve ${car?.make} ${car?.model}?` : `Are you sure you want to sell ${car?.make} ${car?.model}?`;
  const warningMessage = action === actions.sellAction ? 'This action is irreversible.' : '';

  return (
    <div className="modal" role="dialog" aria-modal="true" aria-labelledby="modal-title">
      <div className="modal-content">
        <span className="close" onClick={onClose}>&times;</span>
        <h2 id="modal-title" className="modal-title">{title}</h2>
        <p className="modal-message">{message}</p>
        {error && <p className="error-message">{error}</p>}
        <p className="modal-message">Enter the customer details below:</p>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Full name:</label>
            <input
              type="text"
              value={customerFullName}
              onChange={e => setCustomerFullName(e.target.value)}
              required
            />
          </div>
          <div className="form-group">
            <label>Email:</label>
            <input
              type="email"
              value={customerEmail}
              onChange={e => setCustomerEmail(e.target.value)}
              required
            />
          </div>
          <div className="form-group">
            <label>Phone number (numbers only):</label>
            <input
              type="tel"
              value={phoneNumber}
              onChange={e => setPhoneNumber(e.target.value)}
              required
              pattern="\d*" // HTML5 pattern attribute to accept only numeric input
            />
          </div>
          {warningMessage && <p className="warning-message">{warningMessage}</p>}
          <button type="submit">Confirm</button>
          <button type="button" onClick={onClose}>Cancel</button>
        </form>
      </div>
    </div>
  );
};

export default ReservationAndSaleForm;