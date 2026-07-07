import React from 'react';
import styles from './Bio.module.scss';

export default function Bio({ user }) {
  const bio = user?.bio?.trim()
    || 'Този човек все още не е добавил биография.';

  return (
    <div className={styles.bio}>
      <h2 className={styles.heading}>Био</h2>
      <p className={styles.text}>{bio}</p>
    </div>
  );
}
