import { useState, useEffect } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import styles from './Login.module.scss';
import Button from '../../components/Button/Button';
import { useAuth } from '../../contexts/AuthContext';
import { GOOGLE_LOGIN_URL } from '../../utils/api';

export default function Login() {
  const navigate = useNavigate();
  const location = useLocation();
  const { login, isAuthenticated } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loginError, setLoginError] = useState('');

  // Redirect if already authenticated
  useEffect(() => {
    if (isAuthenticated) {
      navigate('/feed');
    }
  }, [isAuthenticated, navigate]);

  useEffect(() => {
    if (location.state?.authError) {
      setLoginError(location.state.authError);
    }
  }, [location.state]);

  const handleLogin = async (e) => {
    e.preventDefault();
    setLoginError('');

    try {
      await login(email, password);
      navigate('/feed');
    } catch (error) {
      const errorMessage = error.message || error;
      if (errorMessage.includes("invalid email") || errorMessage.includes("invalid password")) {
        setLoginError("Invalid email or password");
      } else if (errorMessage.includes("user not found")) {
        setLoginError("User not found");
      } else {
        setLoginError(errorMessage || "Login failed. Please try again.");
      }
    }
  };

  const handleGoogleLogin = () => {
    window.location.assign(GOOGLE_LOGIN_URL);
  };

  return (
      <div className={styles.loginPage}>
        <div className={styles.container}>
          <div className={styles.imageSection}>
            <div className={styles.imageWrapper}>
              <img
                src="/sculp.jpg"
                alt="Classical sculpture"
                className={styles.image}
              />
            </div>
          </div>
          
          <div className={styles.formSection}>
          <h1 className={styles.title}>LOGIN</h1>
          <p className={styles.subtitle}>How to i get started lorem ipsum dolor at?</p>
          {loginError && <div className={styles.errorText}>{loginError}</div>}
          
          <form className={styles.form} onSubmit={handleLogin}>
            <div className={styles.inputGroup}>
              <span className={styles.icon}>✉️</span>
              <input type="email" placeholder="Email" className={styles.input} value={email} onChange={(e) => setEmail(e.target.value)} required />
            </div>
            
            <div className={styles.inputGroup}>
              <span className={styles.icon}>🔒</span>
              <input type="password" placeholder="Password" className={styles.input} value={password} onChange={(e) => setPassword(e.target.value)} required />
            </div>
            
            <Button className={styles.loginButton} type="submit">Login Now</Button>
          </form>
          
          <p className={styles.orText}>
            or <Link to="/signup" className={styles.link}>sign up</Link>
          </p>
          
          <p className={styles.socialHeader}>Login with Others</p>
          
          <Button className={styles.googleButton} type="button" onClick={handleGoogleLogin}>
            <span className={styles.googleLogo}>
              <svg width="20" height="20" viewBox="0 0 24 24">
                <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
              </svg>
            </span>
            Login with Google
          </Button>
          </div>
        </div>
      </div>
  );
}
