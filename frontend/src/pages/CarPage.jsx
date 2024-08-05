import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import CarList from '../components/CarList';
import CarForm from '../components/CarForm';
import Tab from '../components/Tab';
import '../styles.css';
import carStatuses from '../constants/carStatuses';

const CarPage = () => {
  const apiUrl = process.env.REACT_APP_API_URL;
  const [cars, setCars] = useState([]);
  const [currentTab, setCurrentTab] = useState(carStatuses.StatusAvailable); // State to manage tabs (available, reserved, sold)
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [carToEdit, setCarToEdit] = useState(null);
  const [error, setError] = useState(null);

  const fetchCars = useCallback(async (status) => {
    try {
      const response = await axios.get(`${apiUrl}/cars/${status}`);
      setCars(response.data);
    } catch (error) {
      console.error('Error fetching cars:', error);
      setError('Failed to fetch cars. Please try again later.');
    }
  }, [apiUrl]);

  useEffect(() => {
    fetchCars(currentTab);
  }, [currentTab, fetchCars]);

  const handleTabChange = (tab) => {
    setCurrentTab(tab);
  };

  const handleAddCar = () => {
    setIsModalOpen(true);
    setCarToEdit(null);
  };

  const handleEditCar = (car) => {
    setIsModalOpen(true);
    setCarToEdit(car);
  };

  const handleFormSubmit = () => {
    setIsModalOpen(false);
    fetchCars(currentTab);
  };

  const handleActionComplete = () => {
    fetchCars(currentTab);
  };

  return (
    <div>
      <div className="tabs">
        <Tab label="Available" onClick={() => handleTabChange(carStatuses.StatusAvailable)} active={currentTab === carStatuses.StatusAvailable} />
        <Tab label="Reserved" onClick={() => handleTabChange(carStatuses.StatusReserved)} active={currentTab === carStatuses.StatusReserved} />
        <Tab label="Sold" onClick={() => handleTabChange(carStatuses.StatusSold)} active={currentTab === carStatuses.StatusSold} />
      </div>
      <div>
        {/* Button to add new car */}
        {currentTab === carStatuses.StatusAvailable && (
          <button className="add-car-button" onClick={handleAddCar}>Add Car</button>
        )}
        {/* Modal for adding/editing cars */}
        {isModalOpen && (
          <div className="modal">
            <CarForm carToEdit={carToEdit} onSubmit={handleFormSubmit} onClose={() => setIsModalOpen(false)} />
          </div>
        )}
      </div>
      <div className="tab-content">
        {error && <div className="error-message">{error}</div>}
        {cars && 
          <CarList cars={cars} onEdit={handleEditCar} onActionComplete={handleActionComplete} />
        }
      </div>
    </div>
  );
};

export default CarPage;