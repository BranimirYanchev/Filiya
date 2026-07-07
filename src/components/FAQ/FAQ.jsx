import React, { useState } from 'react';
import styles from './FAQ.module.scss';
import Button from '../Button/Button';

const questions = [
  {
    id: 1,
    question: 'The expense windows adapted sin: Wrong video dram.',
    answer: 'Offending belonging protection gentleman tube oh you called or michele in. Wrongness unlooked bodyly. She met hammed sir breeding her.',
  },
  {
    id: 2,
    question: 'Sin curiosity day assurance bed necessary?',
    answer: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
  }, 
  {
    id: 3,
    question: 'Produce say the ice moments parties?',
    answer: 'Ut enim ad minim veniam, quis nostrud exercitation ullamco.',
  },
  {
    id: 4,
    question: 'Simple innate summer fat appear basket his desire joy?',
    answer: 'Duis aute irure dolor in reprehenderit in voluptate velit esse cillum.',
  },
  {
    id: 5,
    question: 'Outward clothes promise at gravity do curried?',
    answer: 'Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia.',
  },
];

export default function FAQ() {
  const [expandedId, setExpandedId] = useState(1);

  const toggleQuestion = (id) => {
    setExpandedId(expandedId === id ? null : id);
  };

  return (
    <div className={styles.container}>
      <div className={styles.separator}></div>
      <h2 className={styles.title}>Fequently asked questions</h2>
      <div className={styles.content}>
        <div className={styles.accordion}>
          {questions.map((item) => (
            <div key={item.id} className={styles.questionItem}>
              <button
                className={styles.questionButton}
                onClick={() => toggleQuestion(item.id)}
                aria-expanded={expandedId === item.id}
              >
                <span className={styles.questionText}>{item.question}</span>
                <span className={styles.icon}>{expandedId === item.id ? '−' : '+'}</span>
              </button>
              {expandedId === item.id && (
                <div className={styles.answer}>
                  <p>{item.answer}</p>
                </div>
              )}
            </div>
          ))}
        </div>
        <div className={styles.ctaCard}>
          <div className={styles.iconBox}></div>
          <h3 className={styles.ctaTitle}>Do you have more questions?</h3>
          <p className={styles.ctaText}>
            Lorem ipsum dolor sit amet, consectetur adipiscing elit. Ut elit tellus, luctus nec
            ullamcorper mattis, pulvinar dapibus leo.
          </p>
          <Button className={styles.ctaButton}>Check a Demo Call</Button>
        </div>
      </div>
    </div>
  );
}

