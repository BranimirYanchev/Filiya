import React from 'react';
import styles from './Button.module.scss';

export default function Button({ children, onClick, className, ...props }) {
  return (
    <button 
      className={`${styles.button} ${className || ''}`} 
      onClick={onClick}
      {...props}
    >
      {children}
    </button>
  );
}
