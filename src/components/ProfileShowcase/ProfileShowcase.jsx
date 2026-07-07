import React, { useEffect, useMemo, useRef, useState } from 'react';
import { Link } from 'react-router-dom';
import styles from './ProfileShowcase.module.scss';
import { resolveApiUrl } from '../../utils/api';

const DEFAULT_BATCH_SIZE = 4;
const DEFAULT_BUFFER_SIZE = 8;

const getInitials = (fullName, username, email) => {
  const baseValue = fullName || username || email || 'F';
  const parts = baseValue.trim().split(/\s+/).filter(Boolean);

  if (parts.length === 1) {
    return parts[0].slice(0, 2).toUpperCase();
  }

  return parts
    .slice(0, 2)
    .map((part) => part[0])
    .join('')
    .toUpperCase();
};

const formatCompactNumber = (value) => (
  new Intl.NumberFormat('bg-BG', {
    notation: 'compact',
    maximumFractionDigits: 1,
  }).format(value || 0)
);

function ProfileAvatar({ profile }) {
  const pictureUrl = resolveApiUrl(profile.profile_picture || profile.ProfilePicture || '');
  const fullName = profile.full_name || profile.FullName || '';
  const username = profile.username || '';
  const email = profile.email || '';

  if (pictureUrl) {
    return (
      <img
        src={pictureUrl}
        alt={fullName || username || 'Профил'}
        className={styles.avatarImage}
      />
    );
  }

  return (
    <div className={styles.avatarFallback}>
      {getInitials(fullName, username, email)}
    </div>
  );
}

export default function ProfileShowcase({
  title,
  subtitle,
  profiles,
  loading,
  error,
  emptyMessage,
  variant = 'discover',
}) {
  const sentinelRef = useRef(null);
  const [visibleCount, setVisibleCount] = useState(DEFAULT_BATCH_SIZE);
  const [bufferCount, setBufferCount] = useState(DEFAULT_BUFFER_SIZE);

  useEffect(() => {
    setVisibleCount(Math.min(DEFAULT_BATCH_SIZE, profiles.length));
    setBufferCount(Math.min(DEFAULT_BUFFER_SIZE, profiles.length));
  }, [profiles]);

  useEffect(() => {
    const sentinelElement = sentinelRef.current;

    if (!sentinelElement || visibleCount >= profiles.length) {
      return undefined;
    }

    const observer = new IntersectionObserver(
      (entries) => {
        const [entry] = entries;

        if (!entry?.isIntersecting) {
          return;
        }

        setVisibleCount((currentValue) => Math.min(profiles.length, currentValue + DEFAULT_BATCH_SIZE));
        setBufferCount((currentValue) => Math.min(profiles.length, Math.max(currentValue, visibleCount + DEFAULT_BATCH_SIZE)));
      },
      { rootMargin: '240px 0px' },
    );

    observer.observe(sentinelElement);

    return () => {
      observer.disconnect();
    };
  }, [profiles.length, visibleCount]);

  const visibleProfiles = useMemo(
    () => profiles.slice(0, Math.min(visibleCount, bufferCount)),
    [bufferCount, profiles, visibleCount],
  );

  return (
    <section className={`${styles.section} ${styles[variant]}`}>
      <div className={styles.header}>
        <h2 className={styles.title}>{title}</h2>
        {subtitle && <p className={styles.subtitle}>{subtitle}</p>}
      </div>

      {loading && (
        <div className={styles.feedbackCard}>Зареждаме профилите...</div>
      )}

      {!loading && error && (
        <div className={styles.feedbackCard}>{error}</div>
      )}

      {!loading && !error && visibleProfiles.length === 0 && (
        <div className={styles.feedbackCard}>{emptyMessage}</div>
      )}

      {!loading && !error && visibleProfiles.length > 0 && (
        <>
          <div className={styles.grid}>
            {visibleProfiles.map((profile) => {
              const profileId = profile.id || profile.ID;
              const fullName = profile.full_name || profile.FullName || 'Потребител';
              const username = profile.username || '';
              const bio = profile.bio || 'Все още няма добавена биография.';
              const postsCount = profile.stats?.posts || 0;
              const likesCount = profile.stats?.likes || 0;
              const commentsCount = profile.stats?.comments || 0;

              return (
                <Link
                  key={profileId}
                  to={`/profile/${profileId}`}
                  className={styles.card}
                >
                  <div className={styles.avatarShell}>
                    <ProfileAvatar profile={profile} />
                  </div>

                  <div className={styles.cardContent}>
                    <div className={styles.identity}>
                      <h3 className={styles.name}>{fullName}</h3>
                      {username && (
                        <p className={styles.username}>@{username}</p>
                      )}
                    </div>

                    <p className={styles.bio}>{bio}</p>

                    <div className={styles.statsRow}>
                      <span className={styles.statItem}>
                        <strong>{formatCompactNumber(postsCount)}</strong> поста
                      </span>
                      <span className={styles.statItem}>
                        <strong>{formatCompactNumber(likesCount)}</strong> харесвания
                      </span>
                      <span className={styles.statItem}>
                        <strong>{formatCompactNumber(commentsCount)}</strong> коментара
                      </span>
                    </div>
                  </div>
                </Link>
              );
            })}
          </div>

          {visibleCount < profiles.length && (
            <div ref={sentinelRef} className={styles.loadingTrigger}>
              Зареждаме още профили...
            </div>
          )}
        </>
      )}
    </section>
  );
}
