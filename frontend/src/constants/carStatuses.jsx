// Object.freeze() ensures immutability of the constants.
const carStatuses = Object.freeze({
  StatusAvailable: "available",    // Indicates that the car is available.
  StatusReserved: "reserved",      // Indicates that the car has been reserved by a customer.
  StatusSold: "sold"               // Indicates that the car has been sold.
});

export default carStatuses;