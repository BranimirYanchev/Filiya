import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import styles from './PostCard.module.scss';
import { resolveApiUrl } from '../../utils/api';

const DEFAULT_AVATAR = 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=50&h=50&fit=crop';

const formatTimestamp = (value) => {
  if (!value) {
    return 'Току-що';
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return 'Току-що';
  }

  const diffInMinutes = Math.round((Date.now() - date.getTime()) / 60000);

  if (diffInMinutes < 1) {
    return 'Току-що';
  }

  if (diffInMinutes < 60) {
    return `Преди ${diffInMinutes} мин`;
  }

  const diffInHours = Math.round(diffInMinutes / 60);
  if (diffInHours < 24) {
    return `Преди ${diffInHours} ч`;
  }

  const diffInDays = Math.round(diffInHours / 24);
  if (diffInDays < 7) {
    return `Преди ${diffInDays} дни`;
  }

  return new Intl.DateTimeFormat('bg-BG', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(date);
};

const getCollectionCount = (value) => (
  Array.isArray(value) ? value.length : 0
);

const truncateText = (value, limit = 220) => {
  if (!value || value.length <= limit) {
    return value;
  }

  return `${value.slice(0, limit).trim()}...`;
};

export default function PostCard({ post }) {
  const [isDetailsOpen, setIsDetailsOpen] = useState(false);

  if (!post) {
    return null;
  }

  const author = post.author || {};
  const attachments = Array.isArray(post.attachments) ? post.attachments : [];
  const categories = Array.isArray(post.categories) ? post.categories : [];

  const imageAttachments = attachments.filter((attachment) => (
    attachment?.kind === 'image'
    || attachment?.mime_type?.startsWith('image/')
  ));

  const coverImage = imageAttachments[0];
  const coverImageUrl = resolveApiUrl(
    coverImage?.file_url || coverImage?.file_path || '',
  );
  const authorName = author.full_name || author.FullName || author.username || author.email || 'Потребител';
  const avatarUrl = resolveApiUrl(author.profile_picture || author.ProfilePicture || DEFAULT_AVATAR);
  const title = post.title || 'Без заглавие';
  const content = post.content || '';
  const likesCount = getCollectionCount(post.likes);
  const commentsCount = getCollectionCount(post.comments);
  const authorId = author.id || author.ID;
  const profilePath = authorId ? `/profile/${authorId}` : null;
  const previewContent = truncateText(content);
  const hasMoreToShow = content.length > (previewContent?.length || 0) || attachments.length > 0;

  const renderAuthorBlock = (compact = false) => (
    profilePath ? (
      <Link to={profilePath} className={compact ? styles.authorLinkCompact : styles.authorLink}>
        <img
          src={avatarUrl}
          alt={authorName}
          className={compact ? styles.avatarCompact : styles.avatar}
        />
        <div className={styles.userDetails}>
          <span className={styles.userName}>{authorName}</span>
          <span className={styles.timestamp}>{formatTimestamp(post.created_at)}</span>
        </div>
      </Link>
    ) : (
      <div className={compact ? styles.authorLinkCompact : styles.authorLink}>
        <img
          src={avatarUrl}
          alt={authorName}
          className={compact ? styles.avatarCompact : styles.avatar}
        />
        <div className={styles.userDetails}>
          <span className={styles.userName}>{authorName}</span>
          <span className={styles.timestamp}>{formatTimestamp(post.created_at)}</span>
        </div>
      </div>
    )
  );

  return (
    <>
      <article className={styles.postCard}>
        <div className={styles.postHeader}>
          <div className={styles.userInfo}>
            {renderAuthorBlock(false)}
            <span className={styles.friendsIcon}>{post.is_private ? '🔒' : '🌍'}</span>
          </div>
        </div>

        <div className={styles.postContent}>
          <h3 className={styles.postText}>{title}</h3>
          {previewContent && <p className={styles.postBody}>{previewContent}</p>}

          {categories.length > 0 && (
            <div className={styles.categoryRow}>
              {categories.map((category) => (
                <span key={category.id} className={styles.categoryBadge}>
                  {category.name}
                </span>
              ))}
            </div>
          )}

          {coverImageUrl && (
            <img
              src={coverImageUrl}
              alt={title}
              className={styles.postImage}
            />
          )}
        </div>

        <div className={styles.postEngagement}>
          <div className={styles.engagementItem}>
            <span className={styles.engagementIcon}>❤️</span>
            <span className={styles.engagementCount}>{likesCount}</span>
          </div>
          <div className={styles.engagementItem}>
            <span className={styles.engagementIcon}>🖼️</span>
            <span className={styles.engagementCount}>{attachments.length}</span>
          </div>
          <div className={styles.engagementItem}>
            <span className={styles.engagementIcon}>💬</span>
            <span className={styles.engagementCount}>{commentsCount}</span>
          </div>
        </div>

        <div className={styles.profileActionRow}>
          {profilePath && (
            <Link to={profilePath} className={styles.profileActionLink}>
              Виж профила на автора
            </Link>
          )}
          {hasMoreToShow && (
            <button
              type="button"
              className={styles.seeMoreButton}
              onClick={() => setIsDetailsOpen(true)}
            >
              Виж още
            </button>
          )}
        </div>
      </article>

      {isDetailsOpen && (
        <div className={styles.modalOverlay} onClick={() => setIsDetailsOpen(false)}>
          <div
            className={styles.modalCard}
            onClick={(event) => event.stopPropagation()}
          >
            <button
              type="button"
              className={styles.modalCloseButton}
              onClick={() => setIsDetailsOpen(false)}
            >
              ×
            </button>

            <div className={styles.modalCategoryRow}>
              {categories.length > 0 ? (
                categories.map((category) => (
                  <span key={category.id} className={styles.modalCategoryBadge}>
                    {category.name}
                  </span>
                ))
              ) : (
                <span className={styles.modalCategoryBadge}>Пост</span>
              )}
            </div>

            <h2 className={styles.modalTitle}>{title}</h2>
            <p className={styles.modalTimestamp}>{formatTimestamp(post.created_at)}</p>

            <div className={styles.modalAuthorCard}>
              <div className={styles.modalAuthorMeta}>
                {renderAuthorBlock(true)}
                <p className={styles.modalAuthorHint}>Отвори профила на автора, за да разгледаш повече от него.</p>
              </div>
              {profilePath && (
                <Link to={profilePath} className={styles.modalProfileLink}>
                  Към профила
                </Link>
              )}
            </div>

            {content && <p className={styles.modalContent}>{content}</p>}

            {attachments.length > 0 && (
              <div className={styles.attachmentsSection}>
                <div className={styles.attachmentsHeader}>
                  <h3 className={styles.attachmentsTitle}>Прикачени файлове</h3>
                  <span className={styles.attachmentsCount}>{attachments.length}</span>
                </div>

                {imageAttachments.length > 0 && (
                  <div className={styles.attachmentGroup}>
                    <h4 className={styles.attachmentGroupTitle}>Снимки</h4>
                    <div className={styles.attachmentGrid}>
                      {imageAttachments.map((attachment) => {
                        const imageUrl = resolveApiUrl(
                          attachment?.file_url || attachment?.file_path || '',
                        );

                        return (
                          <a
                            key={attachment.id || attachment.file_name}
                            href={imageUrl}
                            target="_blank"
                            rel="noreferrer"
                            className={styles.attachmentCard}
                          >
                            <img
                              src={imageUrl}
                              alt={attachment.file_name || title}
                              className={styles.attachmentImage}
                            />
                            <span className={styles.attachmentName}>
                              {attachment.file_name || 'Снимка'}
                            </span>
                          </a>
                        );
                      })}
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      )}
    </>
  );
}
