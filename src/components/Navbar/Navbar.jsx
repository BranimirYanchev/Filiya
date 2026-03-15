import styles from './Navbar.module.scss';
import React, { useEffect, useState } from 'react';

const Navbar = ({ currentPage = 'home', onNavigate, alwaysScrolled = false }) => {
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    if (alwaysScrolled) {
      setScrolled(true);
      return undefined;
    }

    const handleScroll = () => {
      const sectionHeight = document.querySelector(`section`)?.offsetHeight || 0;
      setScrolled(window.scrollY > sectionHeight - 55);
    };

    handleScroll();
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, [alwaysScrolled]);

  const handleNavigate = (event, nextPage) => {
    event.preventDefault();
    if (typeof onNavigate === 'function') {
      onNavigate(nextPage);
    }
  };

  return (
    <nav className={scrolled || alwaysScrolled ? styles.navbarScrolled : styles.navbar}>
      <div className={styles.logo}>
        <p>φιλία</p>
      </div>
      <ul>
        <li>
          <a
            href="/"
            onClick={(event) => handleNavigate(event, 'home')}
            className={currentPage === 'home' ? styles.activeLink : ''}
          >
            Начало
          </a>
        </li>
        <li>
          <a
            href="/posts"
            onClick={(event) => handleNavigate(event, 'posts')}
            className={currentPage === 'posts' ? styles.activeLink : ''}
          >
            Кръгове
          </a>
        </li>
        <li>
          <a
            href="/login"
            onClick={(event) => handleNavigate(event, 'login')}
            className={styles.login_button}
          >
            Вход
          </a>
        </li>
      </ul>
    </nav>
  );
}

export default Navbar;
