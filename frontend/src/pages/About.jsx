import React from 'react';

const About = () => {
  return (
    <div className="about-page">
      <h1>About This Application</h1>
      <p className="about-text">
        This application is designed to manage the inventory of a car dealership.
      </p>
      <div className="functionality">
        <h2>Functionalities</h2>
        <div className="functionality-item">
          <h3>1. View Cars</h3>
          <p>
            Navigate to the "Car Management" page to view a list of cars categorized
            into available, reserved, and sold. Each car displays its make, model,
            price, year, current status (available, reserved, or sold), and a picture.
            Click on the car image to view it in full size.
          </p>
        </div>
        <div className="functionality-item">
          <h3>2. Add and Edit Cars</h3>
          <p>
            You can add new cars to the inventory or edit existing car details
            (make, model, year, price, and picture). Click the "Add Car" button to
            open a form for adding a new car or click "Edit" on an available car to
            modify its details.
          </p>
        </div>
        <div className="functionality-item">
          <h3>3. Reserve and Sell Cars</h3>
          <p>
            For available cars, you can reserve or sell them to customers. Click
            "Reserve" or "Sell" to initiate the process, where you'll enter the
            customer's information for reservation or sale confirmation. Once a car
            is sold, neither the customer details nor the car details can be edited
            further, making the selling process irreversible.
          </p>
        </div>
        <div className="functionality-item">
          <h3>4. Cancel Reservations and Delete Cars</h3>
          <p>
            You can cancel reservations for cars that are currently reserved by
            clicking the "Cancel" button next to the reserved car. Only available
            cars can be deleted. To delete an available car, click "Delete" next to
            the car and confirm the deletion in the modal that appears.
          </p>
        </div>
      </div>
    </div>
  );
};

export default About;