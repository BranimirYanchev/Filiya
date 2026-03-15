import React, { useState } from 'react';
import { LockKeyhole, Mail, UserRound } from 'lucide-react';
import styles from './Auth.module.scss';
import { loginUser, registerUser } from '../../auth/session';

export default function Auth({ onNavigate }) {
  const [mode, setMode] = useState('login');
  const [loginForm, setLoginForm] = useState({ email: '', password: '' });
  const [registerForm, setRegisterForm] = useState({
    fullName: '',
    email: '',
    password: '',
    repeatedPassword: ''
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');

  const updateLoginField = (field, value) => {
    setLoginForm((prev) => ({ ...prev, [field]: value }));
  };

  const updateRegisterField = (field, value) => {
    setRegisterForm((prev) => ({ ...prev, [field]: value }));
  };

  const handleLogin = async (event) => {
    event.preventDefault();
    setErrorMessage('');
    setIsSubmitting(true);

    try {
      await loginUser(loginForm);
      onNavigate?.('posts');
    } catch (error) {
      setErrorMessage(error.message || 'Неуспешен вход.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleRegister = async (event) => {
    event.preventDefault();
    setErrorMessage('');

    if (registerForm.password !== registerForm.repeatedPassword) {
      setErrorMessage('Паролите не съвпадат.');
      return;
    }

    setIsSubmitting(true);

    try {
      await registerUser(registerForm);
      onNavigate?.('posts');
    } catch (error) {
      setErrorMessage(error.message || 'Неуспешна регистрация.');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <main className={styles.authPage}>
      <div className={mode === 'signup' ? `${styles.stage} ${styles.stageSignup}` : styles.stage}>
        <div className={styles.visualOverlay} />
        <button
          type="button"
          className={styles.brandButton}
          onClick={() => onNavigate?.('home')}
        >
          φιλία
        </button>
        <div className={mode === 'signup' ? `${styles.visualCopyStack} ${styles.visualCopyStackSignup}` : styles.visualCopyStack}>
          <div className={mode === 'login' ? `${styles.visualCopy} ${styles.visualCopyActive}` : styles.visualCopy}>
            <p>Класическа естетика, съвременна общност.</p>
            <h1>Влез в пространството за идеи, текстове и кръгове.</h1>
          </div>
          <div className={mode === 'signup' ? `${styles.visualCopy} ${styles.visualCopyActive}` : styles.visualCopy}>
            <p>Форма, движение и ритъм.</p>
            <h1>Създай профил и влез плавно в мрежата на интересите си.</h1>
          </div>
        </div>

        <div className={mode === 'signup' ? `${styles.formDock} ${styles.formDockSignup}` : styles.formDock}>
          <div className={mode === 'signup' ? `${styles.formTrack} ${styles.formTrackSignup}` : styles.formTrack}>
            <section className={styles.formPanel}>
              <div className={styles.formShell}>
                <span className={styles.eyebrow}>Добре дошъл</span>
                <h2>Вход</h2>
                <p>Влез в профила си и продължи към съдържанието, което следиш.</p>

                <form className={styles.form} onSubmit={handleLogin}>
                  <label className={styles.field}>
                    <Mail size={22} />
                    <input
                      type="email"
                      placeholder="Имейл"
                      value={loginForm.email}
                      onChange={(event) => updateLoginField('email', event.target.value)}
                    />
                  </label>
                  <label className={styles.field}>
                    <LockKeyhole size={22} />
                    <input
                      type="password"
                      placeholder="Парола"
                      value={loginForm.password}
                      onChange={(event) => updateLoginField('password', event.target.value)}
                    />
                  </label>
                  {errorMessage && mode === 'login' ? <p className={styles.authError}>{errorMessage}</p> : null}
                  <button type="submit" className={styles.primaryButton} disabled={isSubmitting}>
                    {isSubmitting ? 'Изчакване...' : 'Влез'}
                  </button>
                </form>

                <div className={styles.switchLine}>
                  <span>Нямаш профил?</span>
                  <button type="button" onClick={() => setMode('signup')}>
                    Регистрация
                  </button>
                </div>

                <div className={styles.separator}>
                  <span>Вход с други профили</span>
                </div>

                <button type="button" className={styles.googleButton}>
                  <span className={styles.googleMark}>G</span>
                  Продължи с Google
                </button>
              </div>
            </section>

            <section className={styles.formPanel}>
              <div className={styles.formShell}>
                <span className={styles.eyebrow}>Ново начало</span>
                <h2>Регистрация</h2>
                <p>Създай профил и започни да подреждаш своя интелектуален кръг.</p>

                <form className={styles.form} onSubmit={handleRegister}>
                  <label className={styles.field}>
                    <UserRound size={22} />
                    <input
                      type="text"
                      placeholder="Име"
                      value={registerForm.fullName}
                      onChange={(event) => updateRegisterField('fullName', event.target.value)}
                    />
                  </label>
                  <label className={styles.field}>
                    <Mail size={22} />
                    <input
                      type="email"
                      placeholder="Имейл"
                      value={registerForm.email}
                      onChange={(event) => updateRegisterField('email', event.target.value)}
                    />
                  </label>
                  <label className={styles.field}>
                    <LockKeyhole size={22} />
                    <input
                      type="password"
                      placeholder="Парола"
                      value={registerForm.password}
                      onChange={(event) => updateRegisterField('password', event.target.value)}
                    />
                  </label>
                  <label className={styles.field}>
                    <LockKeyhole size={22} />
                    <input
                      type="password"
                      placeholder="Повтори паролата"
                      value={registerForm.repeatedPassword}
                      onChange={(event) => updateRegisterField('repeatedPassword', event.target.value)}
                    />
                  </label>
                  {errorMessage && mode === 'signup' ? <p className={styles.authError}>{errorMessage}</p> : null}
                  <button type="submit" className={styles.primaryButton} disabled={isSubmitting}>
                    {isSubmitting ? 'Изчакване...' : 'Създай профил'}
                  </button>
                </form>

                <div className={styles.switchLine}>
                  <span>Имаш профил?</span>
                  <button type="button" onClick={() => setMode('login')}>
                    Вход
                  </button>
                </div>

                <div className={styles.separator}>
                  <span>Регистрация с други профили</span>
                </div>

                <button type="button" className={styles.googleButton}>
                  <span className={styles.googleMark}>G</span>
                  Регистрирай се с Google
                </button>
              </div>
            </section>
          </div>
        </div>
      </div>
    </main>
  );
}
