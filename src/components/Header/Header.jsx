import React, { useState, useEffect, useRef } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import styles from './Header.module.scss';
import Button from '../Button/Button';
import { useAuth } from '../../contexts/AuthContext';

export function UnauthenticatedHeader() {
  const navigate = useNavigate();
  const [isScrolled, setIsScrolled] = useState(false);
  const [isVisible, setIsVisible] = useState(true);
  const [lastScrollY, setLastScrollY] = useState(0);

  useEffect(() => {
    const handleScroll = () => {
      const currentScrollY = window.scrollY;
      
      // Add background when scrolled
      setIsScrolled(currentScrollY > 50);
      
      // Hide/show navbar on scroll
      if (currentScrollY > lastScrollY && currentScrollY > 100) {
        setIsVisible(false);
      } else {
        setIsVisible(true);
      }
      
      setLastScrollY(currentScrollY);
    };

    window.addEventListener('scroll', handleScroll, { passive: true });
    return () => window.removeEventListener('scroll', handleScroll);
  }, [lastScrollY]);

  return (
    <header className={`${styles.header} ${isScrolled ? styles.scrolled : ''} ${!isVisible ? styles.hidden : ''}`}>
      <nav className={styles.nav}>
        <Link to="/" className={styles.logo}>Filia</Link>
        <div className={styles.navLinks}>
          <Link to="/#contacts" className={styles.link}>Контакти</Link>
          <Link to="/login" className={styles.link}>Вход</Link>
          <Button 
            className={styles.registerButton}
            onClick={() => navigate('/signup')}
          >
            Регистрация
          </Button>
        </div>
      </nav>
    </header>
  );
}

export function AuthenticatedHeader() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [isScrolled, setIsScrolled] = useState(false);
  const [isVisible, setIsVisible] = useState(true);
  const [lastScrollY, setLastScrollY] = useState(0);
  const [showUserMenu, setShowUserMenu] = useState(false);
  const userMenuRef = useRef(null);

  useEffect(() => {
    const handleScroll = () => {
      const currentScrollY = window.scrollY;
      
      setIsScrolled(currentScrollY > 50);
      
      if (currentScrollY > lastScrollY && currentScrollY > 100) {
        setIsVisible(false);
      } else {
        setIsVisible(true);
      }
      
      setLastScrollY(currentScrollY);
    };

    window.addEventListener('scroll', handleScroll, { passive: true });
    return () => window.removeEventListener('scroll', handleScroll);
  }, [lastScrollY]);

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (userMenuRef.current && !userMenuRef.current.contains(event.target)) {
        setShowUserMenu(false);
      }
    };

    if (showUserMenu) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [showUserMenu]);

  const handleLogout = async () => {
    await logout();
    navigate('/');
    setShowUserMenu(false);
  };

  return (
    <header className={`${styles.header} ${isScrolled ? styles.scrolled : ''} ${!isVisible ? styles.hidden : ''}`}>
      <nav className={styles.nav}>
        <Link to="/" className={styles.logo}>Filia</Link>
        <div className={styles.navLinks}>
          <Link to="/#contacts" className={styles.link}>Контакти</Link>
          <Link to="/feed" className={styles.link}>Постове и Дебати</Link>
          <div className={styles.userMenu} ref={userMenuRef}>
            <Button 
              className={styles.profileButton}
              onClick={() => setShowUserMenu(!showUserMenu)}
            >
              {user?.full_name || user?.FullName || user?.email || 'Профил'}
            </Button>
            {showUserMenu && (
              <div className={styles.dropdown}>
                <Link 
                  to={`/profile${user?.id || user?.ID ? `/${user.id || user.ID}` : ''}`}
                  className={styles.dropdownItem}
                  onClick={() => setShowUserMenu(false)}
                >
                  Профил
                </Link>
                <button 
                  className={styles.dropdownItem}
                  onClick={handleLogout}
                >
                  Изход
                </button>
              </div>
            )}
          </div>
        </div>
      </nav>
    </header>
  );
}