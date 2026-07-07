import React from 'react';
import styles from './Hero.module.scss';

export default function Hero() {
  const heroStyle = {
    backgroundImage: `linear-gradient(135deg, rgba(92, 74, 59, 0.7) 0%, rgba(74, 60, 48, 0.7) 100%), url('${process.env.PUBLIC_URL}/sculp.jpg')`,
    backgroundSize: 'cover',
    backgroundPosition: 'center'
  };

  return (
    <section className={styles.hero} style={heroStyle}>
      <div className={styles.content}>
        <h1 className={styles.title}>φιλία</h1>
        <p className={styles.subtitle}>добра дигла солната крока</p>
      </div>
    </section>
  );
}

