import React, { useState, useEffect } from 'react';
import axios from 'axios';
import '../styles.css';
import carStatuses from '../constants/carStatuses';

const CarForm = ({ carToEdit, onSubmit, onClose }) => {
  const apiUrl = process.env.REACT_APP_API_URL;
  const [make, setMake] = useState('');
  const [model, setModel] = useState('');
  const [year, setYear] = useState('');
  const [price, setPrice] = useState('');
  const [picture, setPicture] = useState(null);
  const [status, setStatus] = useState(carStatuses.StatusAvailable);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (carToEdit) {
      setMake(carToEdit.make || '');
      setModel(carToEdit.model || '');
      setYear(carToEdit.year || '');
      setPrice(carToEdit.price || '');
      setStatus(carToEdit.status || carStatuses.StatusAvailable);
    }
  }, [carToEdit]);

  const handleFileChange = (e) => {
    setPicture(e.target.files[0]);
  };

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (!make || !model || !year || !price || !picture) {
      setError('All fields are required.');
      return;
    }

    if (year < 1900) {
      setError('Year must be greater than or equal to 1900.');
      return;
    }

    if (price < 1) {
      setError('Price must be greater than or equal to 1.');
      return;
    }

    const formData = new FormData();
    formData.append('make', make);
    formData.append('model', model);
    formData.append('year', year);
    formData.append('price', price);
    formData.append('status', status);
    formData.append('picture', picture);

    try {
      let response;
      if (carToEdit) {
        // Update existing car
        response = await axios.put(`${apiUrl}/cars/${carToEdit.id}`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        });
        console.log('Car updated:', response.data); // Logging the response data
      } else {
        // Add new car
        response = await axios.post(`${apiUrl}/cars`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data', // Logging the response data
          },
        });
        console.log('Car created:', response.data);
      }

      onSubmit();
    } catch (error) {
      console.error('Error submitting form:', error);
      setError('There was an error submitting the form. Please try again.');
    }
  };

  return (
    <div className="modal">
      <div className="modal-content">
        <span className="close" onClick={onClose}>&times;</span>
        {error && <div className="error-message">{error}</div>}
        <form onSubmit={handleSubmit}>
          <label htmlFor="make">Make:</label>
          <input
            type="text"
            id="make"
            value={make}
            onChange={(e) => setMake(e.target.value)}
            required
          />
          <br />
          <label htmlFor="model">Model:</label>
          <input
            type="text"
            id="model"
            value={model}
            onChange={(e) => setModel(e.target.value)}
            required
          />
          <br />
          <label htmlFor="year">Year:</label>
          <input
            type="number"
            id="year"
            value={year}
            onChange={(e) => setYear(e.target.value)}
            min="1900"
            required
          />
          <br />
          <label htmlFor="price">Price:</label>
          <input
            type="number"
            id="price"
            step="0.01"
            value={price}
            onChange={(e) => setPrice(e.target.value)}
            min="1"
            required
          />
          <br />
          <label htmlFor="picture">Upload Picture:</label>
          <input
            type="file"
            id="picture"
            onChange={handleFileChange}
            accept=".jpg,.jpeg,.png"
            required
          />
          <br />
          <button type="submit">{carToEdit ? 'Update Car' : 'Add Car'}</button>
          <button type="button" onClick={onClose}>Cancel</button>
        </form>
      </div>
    </div>
  );
};

export default CarForm;