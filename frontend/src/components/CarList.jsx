import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import ReservationAndSaleForm from './ReservationAndSaleForm';
import ConfirmActionModal from './ConfirmActionModal';
import '../styles.css';
import carStatuses from '../constants/carStatuses';
import actions from '../constants/actions';

const CarList = ({ cars, onEdit, onActionComplete }) => {
  const apiUrl = process.env.REACT_APP_API_URL;
  const [action, setAction] = useState('');
  const [showCustomerForm, setShowCustomerForm] = useState(false);
  const [showConfirmActionModal, setShowConfirmActionModal] = useState(false);
  const [selectedCar, setSelectedCar] = useState(null);
  const [carImages, setCarImages] = useState({});

  // Memoize the fetchCarImage function to avoid unnecessary re-renders
  const fetchCarImage = useCallback((carId, pictureId) => {
    axios.get(`${apiUrl}/cars/image/${pictureId}`, { responseType: 'blob' })
      .then(response => {
        const imageUrl = URL.createObjectURL(response.data);
        setCarImages(prevState => ({ ...prevState, [carId]: imageUrl }));
      })
      .catch(error => {
        console.error('Error fetching car image:', error);
      });
  }, [apiUrl]);

  useEffect(() => {
    cars.forEach(car => {
      if (car.picture) {
        fetchCarImage(car.id, car.picture);
      }
    });
  }, [cars, fetchCarImage]);

  const openModal = (action, car, showCustomerForm) => {
    setAction(action);
    setSelectedCar(car);
    setShowCustomerForm(showCustomerForm);
    setShowConfirmActionModal(!showCustomerForm);
  };

  const closeModal = () => {
    setShowCustomerForm(false);
    setShowConfirmActionModal(false);
    setSelectedCar(null);
    setAction('');
  };

  const handleActionComplete = () => {
    onActionComplete();
  };

  const handleAction = (car, action) => {
    switch (action) {
      case 'edit':
        onEdit(car);
        break;
      case actions.reserveAction:
        openModal(actions.reserveAction, car, true);
        break;
      case actions.sellAction:
        openModal(actions.sellAction, car, true);
        break;
      case actions.deleteAction:
        openModal(actions.deleteAction, car, false);
        break;
      case actions.cancelAction:
        openModal(actions.cancelAction, car, false);
        break;
      default:
        break;
    }
  };

  const openFullSizeImage = (imageUrl) => {
    window.open(imageUrl, '_blank');
  };

  return (
    <div className="car-list">
      <ul>
        {cars.map(car => (
          <li key={car.id} className={`car-status-${car.status}`}>
            <a href={carImages[car.id]} target="_blank" rel="noopener noreferrer" onClick={() => openFullSizeImage(carImages[car.id])}>
              <img src={carImages[car.id]} alt={`${car.make} ${car.model}`} />
            </a>
            <div className="car-details">
              <strong>{car.make} {car.model} - ${car.price}</strong> ({car.year}) - {car.status}
            </div>
            <div className="car-actions">
              {car.status === carStatuses.StatusAvailable && (
                <>
                  <button className="action-button" onClick={() => handleAction(car, 'edit')}>Edit</button>
                  <button className="action-button" onClick={() => handleAction(car, actions.reserveAction)}>Reserve</button>
                  <button className="action-button" onClick={() => handleAction(car, actions.sellAction)}>Sell</button>
                  <button className="action-button" onClick={() => handleAction(car, actions.deleteAction)}>Delete</button>
                </>
              )}
              {car.status === carStatuses.StatusReserved && (
                <div className="reserved-details">
                  <span>
                    <strong>Reserved by:</strong> {car.customer.fullName}<br />
                    <strong>Email:</strong> {car.customer.email}<br />
                    <strong>Phone:</strong> {car.customer.phoneNumber}
                  </span>
                  <button className="action-button" onClick={() => handleAction(car, actions.cancelAction)}>Cancel Reservation</button>
                </div>
              )}
              {car.status === carStatuses.StatusSold && (
                <div className="sold-details">
                  <span>
                    <strong>Sold to:</strong> {car.customer.fullName}<br />
                    <strong>Email:</strong> {car.customer.email}<br />
                    <strong>Phone:</strong> {car.customer.phoneNumber}
                  </span>
                </div>
              )}
            </div>
          </li>
        ))}
      </ul>
      <ReservationAndSaleForm
        isOpen={showCustomerForm}
        onClose={closeModal}
        action={action}
        car={selectedCar}
        onActionComplete={handleActionComplete}
      />
      <ConfirmActionModal
        isOpen={showConfirmActionModal}
        onClose={closeModal}
        action={action}
        car={selectedCar}
        onActionComplete={handleActionComplete}
      />
    </div>
  );
};

export default CarList;