// Object.freeze() ensures immutability of the constants.
const actions = Object.freeze({
  reserveAction: "reserve",
  sellAction: "sell",
  deleteAction: "delete",
  cancelAction: "cancel"
});

export default actions;