import React, { useEffect, useMemo, useRef, useState } from 'react';
import styles from './PostShowcase.module.scss';
import PostCard from '../PostCard/PostCard';

const INITIAL_VISIBLE_COUNT = 4;
const INITIAL_BUFFER_COUNT = 8;
const LOAD_STEP = 4;

export default function PostShowcase({
  title,
  subtitle,
  posts,
  loading,
  error,
  emptyMessage,
}) {
  const sentinelRef = useRef(null);
  const [visibleCount, setVisibleCount] = useState(INITIAL_VISIBLE_COUNT);
  const [bufferCount, setBufferCount] = useState(INITIAL_BUFFER_COUNT);

  useEffect(() => {
    setVisibleCount(Math.min(INITIAL_VISIBLE_COUNT, posts.length));
    setBufferCount(Math.min(INITIAL_BUFFER_COUNT, posts.length));
  }, [posts]);

  useEffect(() => {
    const sentinelElement = sentinelRef.current;

    if (!sentinelElement || visibleCount >= posts.length) {
      return undefined;
    }

    const observer = new IntersectionObserver(
      (entries) => {
        const [entry] = entries;

        if (!entry?.isIntersecting) {
          return;
        }

        setBufferCount((currentValue) => Math.min(posts.length, currentValue + LOAD_STEP));
        setVisibleCount((currentValue) => Math.min(posts.length, currentValue + LOAD_STEP));
      },
      { rootMargin: '240px 0px' },
    );

    observer.observe(sentinelElement);

    return () => {
      observer.disconnect();
    };
  }, [posts.length, visibleCount]);

  const visiblePosts = useMemo(
    () => posts.slice(0, Math.min(visibleCount, bufferCount)),
    [bufferCount, posts, visibleCount],
  );

  return (
    <section className={styles.section}>
      <div className={styles.header}>
        <h2 className={styles.title}>{title}</h2>
        {subtitle && <p className={styles.subtitle}>{subtitle}</p>}
      </div>

      {loading && <div className={styles.feedbackCard}>Зареждаме постовете...</div>}

      {!loading && error && <div className={styles.feedbackCard}>{error}</div>}

      {!loading && !error && visiblePosts.length === 0 && (
        <div className={styles.feedbackCard}>{emptyMessage}</div>
      )}

      {!loading && !error && visiblePosts.length > 0 && (
        <>
          <div className={styles.postsList}>
            {visiblePosts.map((post) => (
              <PostCard key={post.id || post.ID} post={post} />
            ))}
          </div>

          {visibleCount < posts.length && (
            <div ref={sentinelRef} className={styles.loadingTrigger}>
              Зареждаме още постове...
            </div>
          )}
        </>
      )}
    </section>
  );
}
