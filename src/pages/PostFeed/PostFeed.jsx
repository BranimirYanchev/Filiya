import React, { useEffect, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import styles from './PostFeed.module.scss';
import Button from '../../components/Button/Button';
import PostCard from '../../components/PostCard/PostCard';
import FriendsList from '../../components/FriendsList/FriendsList';
import { useAuth } from '../../contexts/AuthContext';
import { api, resolveApiUrl } from '../../utils/api';

const extractList = (payload) => {
  if (Array.isArray(payload?.data)) {
    return payload.data;
  }

  if (Array.isArray(payload)) {
    return payload;
  }

  return [];
};

const getPostCategoryIds = (post) => {
  if (Array.isArray(post?.category_ids)) {
    return post.category_ids;
  }

  if (Array.isArray(post?.categories)) {
    return post.categories.map((category) => category.id).filter(Boolean);
  }

  return [];
};

const buildPostFormData = ({ title, content, categoryIds, imageFile }) => {
  const formData = new FormData();

  formData.append('title', title.trim());
  formData.append('content', content.trim());
  formData.append('is_private', 'false');

  categoryIds.forEach((categoryId) => {
    formData.append('category_ids', String(categoryId));
  });

  if (imageFile) {
    formData.append('attachments', imageFile);
  }

  return formData;
};

export default function PostFeed() {
  const navigate = useNavigate();
  const { user, isAuthenticated } = useAuth();
  const [activeTab, setActiveTab] = useState('posts');
  const [sortBy, setSortBy] = useState('recent');
  const [posts, setPosts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [selectedCategoryId, setSelectedCategoryId] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [isComposerOpen, setIsComposerOpen] = useState(false);
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [selectedComposerCategories, setSelectedComposerCategories] = useState([]);
  const [selectedImage, setSelectedImage] = useState(null);
  const [imagePreviewUrl, setImagePreviewUrl] = useState('');
  const [submitError, setSubmitError] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    let isCancelled = false;

    const loadData = async () => {
      try {
        setLoading(true);
        setError('');

        const [postsResponse, categoriesResponse] = await Promise.all([
          api.get('/posts'),
          api.get('/categories'),
        ]);

        if (isCancelled) {
          return;
        }

        setPosts(extractList(postsResponse));
        setCategories(extractList(categoriesResponse));
      } catch (err) {
        if (!isCancelled) {
          setError(err.message || 'Неуспешно зареждане на постовете.');
        }
      } finally {
        if (!isCancelled) {
          setLoading(false);
        }
      }
    };

    loadData();

    return () => {
      isCancelled = true;
    };
  }, []);

  useEffect(() => {
    if (!selectedImage) {
      setImagePreviewUrl('');
      return undefined;
    }

    const nextPreviewUrl = window.URL.createObjectURL(selectedImage);
    setImagePreviewUrl(nextPreviewUrl);

    return () => {
      window.URL.revokeObjectURL(nextPreviewUrl);
    };
  }, [selectedImage]);

  const rootCategories = useMemo(
    () => categories.filter((category) => !category.parent_id),
    [categories],
  );

  const visiblePosts = useMemo(() => {
    const filteredByTab = posts.filter((post) => {
      if (activeTab !== 'tagged') {
        return true;
      }

      const taggedUsers = Array.isArray(post?.tagged_users) ? post.tagged_users : [];
      const currentUserId = user?.id || user?.ID;

      if (!currentUserId) {
        return false;
      }

      return taggedUsers.some((taggedUser) => (
        String(taggedUser?.id || taggedUser?.ID) === String(currentUserId)
      ));
    });

    const filteredByCategory = filteredByTab.filter((post) => {
      if (!selectedCategoryId) {
        return true;
      }

      return getPostCategoryIds(post).includes(selectedCategoryId);
    });

    const sortedPosts = [...filteredByCategory];

    sortedPosts.sort((leftPost, rightPost) => {
      if (sortBy === 'popular') {
        const leftScore = (
          (Array.isArray(leftPost.likes) ? leftPost.likes.length : 0)
          + (Array.isArray(leftPost.comments) ? leftPost.comments.length : 0)
        );
        const rightScore = (
          (Array.isArray(rightPost.likes) ? rightPost.likes.length : 0)
          + (Array.isArray(rightPost.comments) ? rightPost.comments.length : 0)
        );

        if (leftScore !== rightScore) {
          return rightScore - leftScore;
        }
      }

      const leftDate = new Date(leftPost.created_at || 0).getTime();
      const rightDate = new Date(rightPost.created_at || 0).getTime();
      return rightDate - leftDate;
    });

    return sortedPosts;
  }, [activeTab, posts, selectedCategoryId, sortBy, user]);

  const refreshPosts = async () => {
    const postsResponse = await api.get('/posts');
    setPosts(extractList(postsResponse));
  };

  const resetComposer = () => {
    setTitle('');
    setContent('');
    setSelectedComposerCategories([]);
    setSelectedImage(null);
    setSubmitError('');
  };

  const openComposer = () => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }

    setIsComposerOpen(true);
  };

  const closeComposer = () => {
    setIsComposerOpen(false);
    resetComposer();
  };

  const toggleComposerCategory = (categoryId) => {
    setSelectedComposerCategories((currentCategories) => (
      currentCategories.includes(categoryId)
        ? currentCategories.filter((currentCategoryId) => currentCategoryId !== categoryId)
        : [...currentCategories, categoryId]
    ));
  };

  const handleCreatePost = async (event) => {
    event.preventDefault();

    if (!title.trim() || !content.trim()) {
      setSubmitError('Заглавието и съдържанието са задължителни.');
      return;
    }

    try {
      setIsSubmitting(true);
      setSubmitError('');

      const formData = buildPostFormData({
        title,
        content,
        categoryIds: selectedComposerCategories,
        imageFile: selectedImage,
      });

      await api.postForm('/posts', formData);
      await refreshPosts();
      closeComposer();
      setActiveTab('posts');
    } catch (err) {
      setSubmitError(err.message || 'Не успяхме да публикуваме поста.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const avatarUrl = resolveApiUrl(
    user?.profile_picture
    || user?.ProfilePicture
    || 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=40&h=40&fit=crop',
  );

  return (
    <div className={styles.postFeedPage}>
      <div className={styles.topHeader}>
        <span className={styles.headerText}>Постове</span>
      </div>

      <header className={styles.mainHeader}>
        <div className={styles.headerContent}>
          <div className={styles.logo}>Filia</div>

          <div className={styles.centerIcon}>
            <div className={styles.birdIcon}>✦</div>
          </div>

          <div className={styles.headerActions}>
            <Button className={styles.createButton} onClick={openComposer}>
              + Нов пост
            </Button>
            <img
              src={avatarUrl}
              alt={user?.full_name || user?.FullName || 'Потребител'}
              className={styles.userAvatar}
            />
          </div>
        </div>
      </header>

      <div className={styles.container}>
        <aside className={styles.leftSidebar}>
          <nav className={styles.navigation}>
            <button
              className={`${styles.navItem} ${sortBy === 'popular' ? styles.navItemActive : ''}`}
              onClick={() => setSortBy('popular')}
            >
              <span className={styles.navIcon}>🔥</span>
              Популярни
            </button>
            <button
              className={`${styles.navItem} ${sortBy === 'recent' ? styles.navItemActive : ''}`}
              onClick={() => setSortBy('recent')}
            >
              <span className={styles.navIcon}>🕐</span>
              Най-нови
            </button>
          </nav>

          <div className={styles.categories}>
            <h3 className={styles.categoriesTitle}>Категории</h3>
            <div className={styles.categoryList}>
              <button
                className={`${styles.categoryItem} ${selectedCategoryId === null ? styles.categoryItemActive : ''}`}
                onClick={() => setSelectedCategoryId(null)}
              >
                Всички постове
              </button>
              {rootCategories.map((category) => (
                <button
                  key={category.id}
                  className={`${styles.categoryItem} ${selectedCategoryId === category.id ? styles.categoryItemActive : ''}`}
                  onClick={() => setSelectedCategoryId(category.id)}
                >
                  {category.name}
                </button>
              ))}
            </div>
          </div>
        </aside>

        <main className={styles.mainContent}>
          <div className={styles.tabs}>
            <button
              className={`${styles.tab} ${activeTab === 'posts' ? styles.active : ''}`}
              onClick={() => setActiveTab('posts')}
            >
              Постове
            </button>
            <button
              className={`${styles.tab} ${activeTab === 'tagged' ? styles.active : ''}`}
              onClick={() => setActiveTab('tagged')}
            >
              Тагнати постове
            </button>
          </div>

          {loading && (
            <div className={styles.feedbackCard}>Зареждаме постовете...</div>
          )}

          {!loading && error && (
            <div className={styles.feedbackCard}>{error}</div>
          )}

          {!loading && !error && (
            <div className={styles.postsContainer}>
              {visiblePosts.length > 0 ? (
                visiblePosts.map((post) => (
                  <PostCard key={post.id || post.ID} post={post} />
                ))
              ) : (
                <div className={styles.feedbackCard}>
                  {activeTab === 'tagged'
                    ? 'Все още няма тагнати постове за този профил.'
                    : 'Все още няма постове за тези филтри.'}
                </div>
              )}
            </div>
          )}
        </main>

        <aside className={styles.rightSidebar}>
          <FriendsList />
        </aside>
      </div>

      {isComposerOpen && (
        <div className={styles.modalOverlay} onClick={closeComposer}>
          <div
            className={styles.modalCard}
            onClick={(event) => event.stopPropagation()}
          >
            <div className={styles.modalHeader}>
              <h2 className={styles.modalTitle}>Нов пост</h2>
              <button className={styles.closeButton} onClick={closeComposer}>
                ×
              </button>
            </div>

            <form className={styles.composerForm} onSubmit={handleCreatePost}>
              <label className={styles.field}>
                <span className={styles.fieldLabel}>Заглавие</span>
                <input
                  className={styles.textInput}
                  type="text"
                  value={title}
                  onChange={(event) => setTitle(event.target.value)}
                  placeholder="Напиши кратко заглавие"
                  maxLength={160}
                />
              </label>

              <label className={styles.field}>
                <span className={styles.fieldLabel}>Съдържание</span>
                <textarea
                  className={styles.textArea}
                  value={content}
                  onChange={(event) => setContent(event.target.value)}
                  placeholder="Какво искаш да споделиш?"
                  rows={6}
                />
              </label>

              <div className={styles.field}>
                <span className={styles.fieldLabel}>Категории</span>
                <div className={styles.categoryPicker}>
                  {rootCategories.map((category) => (
                    <button
                      key={category.id}
                      type="button"
                      className={`${styles.categoryChip} ${selectedComposerCategories.includes(category.id) ? styles.categoryChipActive : ''}`}
                      onClick={() => toggleComposerCategory(category.id)}
                    >
                      {category.name}
                    </button>
                  ))}
                </div>
              </div>

              <label className={styles.field}>
                <span className={styles.fieldLabel}>Снимка</span>
                <input
                  className={styles.fileInput}
                  type="file"
                  accept="image/*"
                  onChange={(event) => setSelectedImage(event.target.files?.[0] || null)}
                />
              </label>

              {imagePreviewUrl && (
                <div className={styles.previewCard}>
                  <img
                    src={imagePreviewUrl}
                    alt="Преглед на качената снимка"
                    className={styles.previewImage}
                  />
                </div>
              )}

              {submitError && (
                <div className={styles.submitError}>{submitError}</div>
              )}

              <div className={styles.modalActions}>
                <Button type="button" className={styles.secondaryButton} onClick={closeComposer}>
                  Отказ
                </Button>
                <Button type="submit" disabled={isSubmitting}>
                  {isSubmitting ? 'Публикуване...' : 'Публикувай'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
