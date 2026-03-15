import React from 'react';
import styles from './App.module.scss';
import Home from './pages/Home/Home';
import Posts from './pages/Posts/Posts';
import Auth from './pages/Auth/Auth';
import Profile from './pages/Profile/Profile';
import {
  ensureValidSession,
  isAuthenticated
} from './auth/session';

const PROTECTED_PAGES = new Set(['posts', 'profile']);

const resolvePageFromPath = (path) => {
  if (path === '/circles' || path === '/posts') return 'posts';
  if (path === '/login') return 'login';
  if (path === '/profile') return 'profile';
  return 'home';
};

function App() {
  const [page, setPage] = React.useState(resolvePageFromPath(window.location.pathname));

  React.useEffect(() => {
    const handlePopState = () => {
      setPage(resolvePageFromPath(window.location.pathname));
    };

    window.addEventListener('popstate', handlePopState);
    return () => window.removeEventListener('popstate', handlePopState);
  }, []);

  React.useEffect(() => {
    const guardProtectedRoute = async () => {
      if (!PROTECTED_PAGES.has(page)) {
        return;
      }

      if (!isAuthenticated()) {
        if (window.location.pathname !== '/login') {
          window.history.replaceState({}, '', '/login');
        }
        setPage('login');
        return;
      }

      try {
        await ensureValidSession({ force: true });
      } catch (error) {
        if (window.location.pathname !== '/login') {
          window.history.replaceState({}, '', '/login');
        }
        setPage('login');
      }
    };

    guardProtectedRoute();
  }, [page]);

  React.useEffect(() => {
    const handleSessionExpired = () => {
      if (window.location.pathname !== '/login') {
        window.history.replaceState({}, '', '/login');
      }
      setPage('login');
      window.scrollTo({ top: 0, behavior: 'smooth' });
    };

    const handleAuthChanged = (event) => {
      if (!event.detail?.authenticated && PROTECTED_PAGES.has(resolvePageFromPath(window.location.pathname))) {
        handleSessionExpired();
      }
    };

    const intervalId = window.setInterval(async () => {
      if (!isAuthenticated()) {
        return;
      }

      try {
        await ensureValidSession({ force: true });
      } catch (error) {
        handleSessionExpired();
      }
    }, 5 * 60 * 1000);

    window.addEventListener('filia:session-expired', handleSessionExpired);
    window.addEventListener('filia:auth-changed', handleAuthChanged);

    return () => {
      window.clearInterval(intervalId);
      window.removeEventListener('filia:session-expired', handleSessionExpired);
      window.removeEventListener('filia:auth-changed', handleAuthChanged);
    };
  }, []);

  const navigateTo = (nextPage, options = {}) => {
    const basePath = nextPage === 'posts'
      ? '/posts'
      : nextPage === 'login'
        ? '/login'
        : nextPage === 'profile'
          ? '/profile'
        : '/';
    const nextPath = options.search ? `${basePath}${options.search}` : basePath;

    if (PROTECTED_PAGES.has(nextPage) && !isAuthenticated()) {
      if (window.location.pathname !== '/login') {
        window.history.pushState({}, '', '/login');
      }
      setPage('login');
      window.scrollTo({ top: 0, behavior: 'smooth' });
      return;
    }

    if (`${window.location.pathname}${window.location.search}` !== nextPath) {
      window.history.pushState({}, '', nextPath);
    }

    setPage(nextPage);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  return (
    <div className={styles.app}>
      {page === 'posts' ? (
        <Posts onNavigate={navigateTo} />
      ) : page === 'login' ? (
        <Auth onNavigate={navigateTo} />
      ) : page === 'profile' ? (
        <Profile onNavigate={navigateTo} />
      ) : (
        <Home onNavigate={navigateTo} />
      )}
    </div>
  );
}

export default App;
