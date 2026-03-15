import React, { useMemo, useState } from 'react';
import Navbar from '../../components/Navbar/Navbar';
import Footer from '../../components/Footer/Footer';
import styles from './Circles.module.scss';
import { FEED_WEIGHTS } from '../../recommendation/feedConfig';
import {
  buildSpecificQualificationKey,
  CIRCLE_CATEGORIES_STORAGE_KEY,
  getSavedCirclePreferences,
  PERIODS,
  QUALIFICATION_GROUPS
} from '../Posts/categoryFilters';

const getInitialState = () => {
  try {
    const saved = getSavedCirclePreferences();
    return saved || {
      periods: [],
      generalQualifications: [],
      specificQualifications: []
    };
  } catch (error) {
    // Ignore malformed state.
    return {
      periods: [],
      generalQualifications: [],
      specificQualifications: []
    };
  }
};

export default function Circles({ onNavigate }) {
  const initialState = useMemo(() => getInitialState(), []);
  const [selectedPeriods, setSelectedPeriods] = useState(initialState.periods);
  const [selectedGeneralQualifications, setSelectedGeneralQualifications] = useState(initialState.generalQualifications);
  const [selectedSpecificQualifications, setSelectedSpecificQualifications] = useState(initialState.specificQualifications);
  const [feedbackMessage, setFeedbackMessage] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  const togglePeriod = (periodName) => {
    setSelectedPeriods((prev) => (
      prev.includes(periodName)
        ? prev.filter((item) => item !== periodName)
        : [...prev, periodName]
    ));
  };

  const toggleGeneralQualification = (groupName) => {
    const group = QUALIFICATION_GROUPS.find((entry) => entry.name === groupName);
    if (!group) return;

    setSelectedGeneralQualifications((prev) => {
      const isSelected = prev.includes(groupName);
      if (!isSelected) return [...prev, groupName];

      setSelectedSpecificQualifications((specificPrev) => (
        specificPrev.filter((item) => !item.startsWith(`${groupName}::`))
      ));
      return prev.filter((item) => item !== groupName);
    });
  };

  const toggleSpecificQualification = (groupName, itemName) => {
    const key = buildSpecificQualificationKey(groupName, itemName);

    setSelectedSpecificQualifications((prev) => (
      prev.includes(key)
        ? prev.filter((item) => item !== key)
        : [...prev, key]
    ));

    setSelectedGeneralQualifications((prev) => (
      prev.includes(groupName) ? prev : [...prev, groupName]
    ));
  };

  const handleSave = () => {
    if (selectedGeneralQualifications.length === 0) {
      setErrorMessage('Трябва да избереш поне една обща квалификация.');
      setFeedbackMessage('');
      return;
    }

    const payload = {
      periods: selectedPeriods,
      generalQualifications: selectedGeneralQualifications,
      specificQualifications: selectedSpecificQualifications
    };

    window.localStorage.setItem(CIRCLE_CATEGORIES_STORAGE_KEY, JSON.stringify(payload));
    setErrorMessage('');
    setFeedbackMessage('Изборът е запазен успешно.');
  };

  return (
    <div className={styles.circlesPage}>
      <Navbar currentPage="circles" onNavigate={onNavigate} />
      <main className={styles.main}>
        <section className={styles.hero}>
          <h1>Кръгове и Интереси</h1>
          <p>
            Периодизацията е вертикалата, квалификациите са хоризонталата.
            Можеш да избираш много елементи, но минимум една обща квалификация е задължителна.
          </p>
        </section>

        <section className={styles.matrix}>
          <aside className={styles.periodColumn}>
            <h2>Периодизация</h2>
            <div className={styles.optionList}>
              {PERIODS.map((period) => (
                <label key={period} className={styles.optionItem}>
                  <input
                    type="checkbox"
                    checked={selectedPeriods.includes(period)}
                    onChange={() => togglePeriod(period)}
                  />
                  <span>{period}</span>
                </label>
              ))}
            </div>
          </aside>

          <div className={styles.qualificationColumn}>
            <h2>Квалификация</h2>
            <div className={styles.generalRow}>
              {QUALIFICATION_GROUPS.map((group) => (
                <label key={group.name} className={styles.generalCard}>
                  <input
                    type="checkbox"
                    checked={selectedGeneralQualifications.includes(group.name)}
                    onChange={() => toggleGeneralQualification(group.name)}
                  />
                  <span>{group.name}</span>
                </label>
              ))}
            </div>

            <div className={styles.groupsGrid}>
              {QUALIFICATION_GROUPS.map((group) => (
                <article className={styles.groupCard} key={group.name}>
                  <h3>{group.name}</h3>
                  <div className={styles.optionList}>
                    {group.items.map((item) => (
                      <label key={`${group.name}-${item}`} className={styles.optionItem}>
                        <input
                          type="checkbox"
                          checked={selectedSpecificQualifications.includes(buildSpecificQualificationKey(group.name, item))}
                          onChange={() => toggleSpecificQualification(group.name, item)}
                        />
                        <span>{item}</span>
                      </label>
                    ))}
                  </div>
                </article>
              ))}
            </div>
          </div>
        </section>

        <section className={styles.algorithmCard}>
          <h2>Предложен Feed Алгоритъм</h2>
          <p>
            {FEED_WEIGHTS.interest}% интереси (включително съдържание от приятели),
            {` ${FEED_WEIGHTS.fresh}%`} нови постове и {` ${FEED_WEIGHTS.random}%`} напълно случайно съдържание.
          </p>
        </section>

        <section className={styles.actions}>
          <button type="button" onClick={handleSave}>Запази избор</button>
          {errorMessage ? <p className={styles.error}>{errorMessage}</p> : null}
          {feedbackMessage ? <p className={styles.success}>{feedbackMessage}</p> : null}
        </section>
      </main>
      <Footer />
    </div>
  );
}
