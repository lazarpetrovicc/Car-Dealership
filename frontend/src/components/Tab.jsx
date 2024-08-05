import React from 'react';
import '../styles.css';

const Tab = ({ label, active, onClick }) => {
  return (
    <button
      className={`tab ${active ? 'active' : ''}`}
      onClick={onClick}
      aria-selected={active}
      role="tab"
    >
      {label}
    </button>
  );
};

export default Tab;