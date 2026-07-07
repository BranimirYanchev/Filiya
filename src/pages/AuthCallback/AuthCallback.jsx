import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';

export default function AuthCallback() {
  const navigate = useNavigate();
  const { completeGoogleAuth } = useAuth();
  const [message, setMessage] = useState('Завършваме Google входа...');

  useEffect(() => {
    let isCancelled = false;

    const finishAuth = async () => {
      const params = new URLSearchParams(window.location.search);
      const authStatus = params.get('auth');

      if (authStatus !== 'success') {
        navigate('/login', {
          replace: true,
          state: { authError: 'Google входът не беше завършен успешно.' }
        });
        return;
      }

      try {
        await completeGoogleAuth();
        if (!isCancelled) {
          navigate('/feed', { replace: true });
        }
      } catch (error) {
        if (!isCancelled) {
          setMessage('Google входът не успя. Пренасочваме те към страницата за вход...');
          window.setTimeout(() => {
            navigate('/login', {
              replace: true,
              state: { authError: error.message || 'Google входът не успя.' }
            });
          }, 900);
        }
      }
    };

    finishAuth();

    return () => {
      isCancelled = true;
    };
  }, [completeGoogleAuth, navigate]);

  return (
    <div style={{ minHeight: '100vh', display: 'grid', placeItems: 'center', padding: '2rem', textAlign: 'center' }}>
      <div>
        <h1 style={{ marginBottom: '1rem' }}>Filia</h1>
        <p>{message}</p>
      </div>
    </div>
  );
}
