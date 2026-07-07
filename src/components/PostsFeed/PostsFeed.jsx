import React, { useEffect, useMemo, useState } from 'react';
import styles from './PostsFeed.module.scss';
import PostCard from '../PostCard/PostCard';
import { api } from '../../utils/api';

const extractPosts = (payload) => {
  if (Array.isArray(payload?.data)) {
    return payload.data;
  }

  if (Array.isArray(payload)) {
    return payload;
  }

  return [];
};

const getPostAuthorId = (post) => (
  String(post?.author_id || post?.author?.id || post?.author?.ID || '')
);

export default function PostsFeed({ user, isOwnProfile = false }) {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    let isCancelled = false;

    const loadPosts = async () => {
      try {
        setLoading(true);
        setError('');

        const response = await api.get('/posts');

        if (isCancelled) {
          return;
        }

        const authorId = String(user?.id || user?.ID || '');
        const allPosts = extractPosts(response);
        const filteredPosts = allPosts
          .filter((post) => getPostAuthorId(post) === authorId)
          .sort((leftPost, rightPost) => (
            new Date(rightPost?.created_at || 0).getTime()
            - new Date(leftPost?.created_at || 0).getTime()
          ));

        setPosts(filteredPosts);
      } catch (requestError) {
        if (!isCancelled) {
          setError(requestError.message || 'Не успяхме да заредим постовете.');
        }
      } finally {
        if (!isCancelled) {
          setLoading(false);
        }
      }
    };

    if (user?.id || user?.ID) {
      loadPosts();
      return () => {
        isCancelled = true;
      };
    }

    setPosts([]);
    setLoading(false);
    return () => {
      isCancelled = true;
    };
  }, [user]);

  const authorName = user?.full_name || user?.FullName || user?.username || user?.email || 'този профил';

  const sectionTitle = useMemo(
    () => (isOwnProfile ? 'Моите постове' : `Постове на ${authorName}`),
    [authorName, isOwnProfile],
  );

  const emptyMessage = isOwnProfile
    ? 'Все още нямаш публикувани постове.'
    : `${authorName} все още няма публикувани постове.`;

  return (
    <section className={styles.section}>
      <div className={styles.header}>
        <div>
          <h2 className={styles.title}>{sectionTitle}</h2>
          <p className={styles.subtitle}>
            {posts.length > 0
              ? `${posts.length} ${posts.length === 1 ? 'публикация' : 'публикации'}`
              : 'Показваме всички налични публикации за този профил.'}
          </p>
        </div>
      </div>

      {loading && <div className={styles.feedback}>Зареждаме постовете...</div>}

      {!loading && error && <div className={styles.feedback}>{error}</div>}

      {!loading && !error && posts.length === 0 && (
        <div className={styles.feedback}>{emptyMessage}</div>
      )}

      {!loading && !error && posts.length > 0 && (
        <div className={styles.postsFeed}>
          {posts.map((post) => (
            <PostCard key={post.id || post.ID} post={post} />
          ))}
        </div>
      )}
    </section>
  );
}
