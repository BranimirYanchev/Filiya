import React, { useEffect, useState } from 'react';
import {
  ArrowLeft,
  Globe,
  CirclePlus,
  CircleEllipsis,
  Mail,
  MessageCircle,
  Heart,
  Link2,
  MessageSquareMore,
  PencilLine,
  Share2,
  Send,
  UserRound,
  X
} from 'lucide-react';
import styles from './Profile.module.scss';
import {
  fetchCurrentUser,
  fetchCurrentUserPosts,
  fetchUserFriends,
  updateProfileBasic,
  updateProfileSensitive
} from '../../auth/session';

const EMPTY_POSTS = [];

const THEMES = [
  {
    name: 'Пясък',
    surface: '#e9dfcf',
    accent: '#d9c6a6',
    accentStrong: '#c1a270',
    coverBlend: 'linear-gradient(135deg, rgba(70, 36, 16, 0.16), rgba(70, 36, 16, 0.42))'
  },
  {
    name: 'Маслина',
    surface: '#e3e3d7',
    accent: '#c8ccb2',
    accentStrong: '#8a9770',
    coverBlend: 'linear-gradient(135deg, rgba(44, 58, 32, 0.18), rgba(44, 58, 32, 0.44))'
  },
  {
    name: 'Камък',
    surface: '#e4e1dd',
    accent: '#d1cbc3',
    accentStrong: '#8e8173',
    coverBlend: 'linear-gradient(135deg, rgba(49, 42, 36, 0.14), rgba(49, 42, 36, 0.42))'
  }
];

const PROFILE_STORAGE_KEY = 'filia-profile-draft';

const SOCIAL_FIELDS = [
  { key: 'website', label: 'Уебсайт', icon: Globe, placeholder: 'https://example.com' },
  { key: 'facebook', label: 'Facebook', icon: Link2, placeholder: 'https://facebook.com/...' },
  { key: 'instagram', label: 'Instagram', icon: Link2, placeholder: 'https://instagram.com/...' },
  { key: 'discord', label: 'Discord', icon: MessageSquareMore, placeholder: 'https://discord.gg/...' },
  { key: 'linkedin', label: 'LinkedIn', icon: Link2, placeholder: 'https://linkedin.com/in/...' },
  { key: 'github', label: 'GitHub', icon: Link2, placeholder: 'https://github.com/...' },
  { key: 'youtube', label: 'YouTube', icon: Link2, placeholder: 'https://youtube.com/@...' },
  { key: 'telegram', label: 'Telegram', icon: Send, placeholder: 'https://t.me/...' },
  { key: 'x', label: 'X / Twitter', icon: Link2, placeholder: 'https://x.com/...' },
  { key: 'tiktok', label: 'TikTok', icon: Link2, placeholder: 'https://tiktok.com/@...' }
];

const DEFAULT_PROFILE = {
  name: '',
  email: '',
  bio: '',
  socials: Object.fromEntries(SOCIAL_FIELDS.map((field) => [field.key, '']))
};

const normalizeText = (value, fallback = '') => {
  if (typeof value !== 'string') return fallback;
  const trimmed = value.trim();
  return trimmed || fallback;
};

const normalizeSocialLink = (value) => normalizeText(value, '');

const toExternalUrl = (value) => {
  const trimmed = normalizeText(value, '');
  if (!trimmed) return '';
  if (/^https?:\/\//i.test(trimmed)) return trimmed;
  return `https://${trimmed}`;
};

const buildDisplayName = (user) => {
  const fullName = normalizeText(user?.full_name);
  if (fullName) return fullName;

  const email = normalizeText(user?.email);
  if (!email) return DEFAULT_PROFILE.name;

  const localPart = email.split('@')[0]?.replace(/[._-]+/g, ' ')?.trim();
  return localPart || DEFAULT_PROFILE.name;
};

const formatRelativeTime = (value) => {
  if (!value) return 'преди малко';
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) return 'преди малко';

  const diffMs = Date.now() - parsed.getTime();
  const diffMinutes = Math.max(1, Math.floor(diffMs / 60000));

  if (diffMinutes < 60) return `преди ${diffMinutes} мин`;

  const diffHours = Math.floor(diffMinutes / 60);
  if (diffHours < 24) return `преди ${diffHours} ч`;

  const diffDays = Math.floor(diffHours / 24);
  return `преди ${diffDays} дни`;
};

const mapUserToProfile = (user, persistedExtras) => ({
  id: user?.id,
  name: buildDisplayName(user),
  email: normalizeText(user?.email, DEFAULT_PROFILE.email),
  bio: normalizeText(user?.bio, DEFAULT_PROFILE.bio),
  profilePicture: normalizeText(user?.profile_picture, ''),
  socials: Object.fromEntries(
    SOCIAL_FIELDS.map((field) => [field.key, normalizeSocialLink(persistedExtras?.socials?.[field.key])])
  )
});

const mapPosts = (items) => (
  Array.isArray(items)
    ? items.map((post) => ({
        id: post.id,
        author: normalizeText(post?.author?.full_name, normalizeText(post?.author?.email, 'Автор')),
        time: formatRelativeTime(post.created_at),
        text: normalizeText(post.content, 'Без съдържание.'),
        image: null,
        likes: Array.isArray(post.likes) ? post.likes.length : 0,
        shares: Array.isArray(post.tags) ? post.tags.length : 0,
        comments: Array.isArray(post.comments) ? post.comments.length : 0
      }))
    : EMPTY_POSTS
);

const mapFriends = (items) => (
  Array.isArray(items)
    ? items.map((friend, index) => ({
        id: friend.id || `${friend.email}-${index}`,
        name: normalizeText(friend.full_name, normalizeText(friend.email, 'Потребител')),
        email: normalizeText(friend.email, 'Без имейл')
      }))
    : []
);

export default function Profile({ onNavigate }) {
  const [activeTab, setActiveTab] = useState('posts');
  const [profile, setProfile] = useState(DEFAULT_PROFILE);
  const [friends, setFriends] = useState([]);
  const [posts, setPosts] = useState([]);
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [draft, setDraft] = useState({ ...DEFAULT_PROFILE, oldPassword: '' });
  const [themeIndex] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [loadError, setLoadError] = useState('');
  const [saveError, setSaveError] = useState('');

  useEffect(() => {
    let isMounted = true;

    const loadProfileData = async () => {
      setIsLoading(true);
      setLoadError('');

      try {
        const raw = window.localStorage.getItem(PROFILE_STORAGE_KEY);
        const persistedExtras = raw ? JSON.parse(raw) : {};

        const user = await fetchCurrentUser();
        const nextProfile = mapUserToProfile(user, persistedExtras);
        const [friendItems, postItems] = await Promise.allSettled([
          fetchUserFriends(),
          fetchCurrentUserPosts(user?.id)
        ]);

        if (!isMounted) return;

        setProfile(nextProfile);
        setDraft({ ...nextProfile, oldPassword: '' });
        setFriends(friendItems.status === 'fulfilled' ? mapFriends(friendItems.value) : []);
        setPosts(postItems.status === 'fulfilled' ? mapPosts(postItems.value) : EMPTY_POSTS);
      } catch (error) {
        if (!isMounted) return;
        setLoadError(error.message || 'Неуспешно зареждане на профила.');
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    };

    loadProfileData();

    return () => {
      isMounted = false;
    };
  }, []);

  useEffect(() => {
    try {
      window.localStorage.setItem(PROFILE_STORAGE_KEY, JSON.stringify({
        socials: profile.socials
      }));
    } catch (error) {
      // Ignore local storage write failures.
    }
  }, [profile.socials]);

  const theme = THEMES[themeIndex];
  const activePosts = activeTab === 'posts' ? posts : EMPTY_POSTS;
  const activeSocials = SOCIAL_FIELDS
    .map((field) => ({
      ...field,
      href: toExternalUrl(profile.socials[field.key])
    }))
    .filter((field) => field.href);

  const openEditModal = () => {
    setSaveError('');
    setDraft({ ...profile, oldPassword: '' });
    setIsEditOpen(true);
  };

  const closeEditModal = () => {
    setIsEditOpen(false);
    setSaveError('');
    setDraft({ ...profile, oldPassword: '' });
  };

  const handleDraftChange = (field) => (event) => {
    const value = event.target.value;
    setDraft((current) => ({
      ...current,
      [field]: value
    }));
  };

  const handleSocialChange = (field) => (event) => {
    const value = event.target.value;
    setDraft((current) => ({
      ...current,
      socials: {
        ...current.socials,
        [field]: value
      }
    }));
  };

  const saveProfile = (event) => {
    event.preventDefault();
    setSaveError('');
    setIsSaving(true);

    const nextProfile = {
      ...profile,
      name: normalizeText(draft.name, profile.name),
      email: normalizeText(draft.email, profile.email),
      bio: normalizeText(draft.bio, profile.bio),
      socials: {
        ...Object.fromEntries(
          SOCIAL_FIELDS.map((field) => [field.key, normalizeSocialLink(draft.socials[field.key])])
        )
      }
    };

    const bioChanged = nextProfile.bio !== profile.bio;
    const sensitiveChanged = nextProfile.name !== profile.name || nextProfile.email !== profile.email;

    const runSave = async () => {
      if (bioChanged) {
        await updateProfileBasic({ bio: nextProfile.bio });
      }

      if (sensitiveChanged) {
        if (!draft.oldPassword.trim()) {
          throw new Error('За смяна на име или имейл въведи текущата си парола.');
        }

        await updateProfileSensitive({
          fullName: nextProfile.name !== profile.name ? nextProfile.name : undefined,
          email: nextProfile.email !== profile.email ? nextProfile.email : undefined,
          oldPassword: draft.oldPassword.trim()
        });
      }

      setProfile(nextProfile);
      setDraft({ ...nextProfile, oldPassword: '' });
      setIsEditOpen(false);
    };

    runSave()
      .catch((error) => {
        setSaveError(error.message || 'Неуспешно запазване на профила.');
      })
      .finally(() => {
        setIsSaving(false);
      });
  };

  return (
    <main
      className={styles.page}
      style={{
        '--profile-surface': theme.surface,
        '--profile-accent': theme.accent,
        '--profile-accent-strong': theme.accentStrong,
        '--profile-cover-blend': theme.coverBlend
      }}
    >
      <div className={styles.shell}>
        <section className={styles.heroCard}>
          <div className={styles.heroTopBar}>
            <button type="button" className={styles.backButton} onClick={() => onNavigate?.('posts')}>
              <ArrowLeft size={28} />
            </button>
            <h1>{profile.name || 'Профил'}</h1>
          </div>

          <div className={styles.cover} />

          <div className={styles.profileRow}>
            <div className={styles.avatar}>
              <UserRound size={74} />
            </div>

            <div className={styles.identityBlock}>
              <h2>{profile.name || 'Профил'}</h2>
              {profile.bio ? <p>{profile.bio}</p> : null}
            </div>

            <div className={styles.heroActions}>
              <button type="button" className={styles.actionButton} onClick={openEditModal}>
                <PencilLine size={20} />
                Редактирай профил
              </button>
            </div>
          </div>
        </section>

        <section className={styles.grid}>
          <aside className={styles.aboutCard}>
            <h3>Информация</h3>
            {loadError ? <p className={styles.statusMessage}>{loadError}</p> : null}
            {isLoading ? <p className={styles.statusMessage}>Зареждане на профила...</p> : null}
            <div className={styles.infoList}>
              {profile.email ? (
                <div className={styles.infoItem}>
                  <Mail size={24} />
                  <span>{profile.email}</span>
                </div>
              ) : null}
              <div className={styles.infoItem}>
                <MessageCircle size={24} />
                <span>{posts.length} публикувани поста</span>
              </div>
              <div className={styles.infoItem}>
                <CirclePlus size={24} />
                <span>{friends.length} приятели</span>
              </div>
            </div>

            <div className={styles.bioBlock}>
              <div className={styles.bioHeader}>
                <h3>Биография</h3>
              </div>
              {profile.bio ? (
                <p>{profile.bio}</p>
              ) : (
                <p className={styles.statusMessage}>Добави кратко описание от редакцията на профила.</p>
              )}
            </div>

            <div className={styles.socialBlock}>
              <div className={styles.bioHeader}>
                <h3>Социални мрежи</h3>
              </div>
              <div className={styles.socialActions}>
                {activeSocials.length ? activeSocials.map((item) => {
                const Icon = item.icon;
                return (
                  <a
                    key={item.key}
                    className={styles.socialButton}
                    href={item.href}
                    target="_blank"
                    rel="noreferrer"
                    aria-label={item.label}
                    title={item.label}
                  >
                    <Icon size={18} />
                    <span>{item.label}</span>
                  </a>
                );
              }) : (
                  <p className={styles.statusMessage}>Все още няма добавени социални мрежи.</p>
                )}
              </div>
            </div>
          </aside>

          <section className={styles.postsCard}>
            <div className={styles.tabs}>
              <div className={styles.tabList}>
                <button
                  type="button"
                  className={activeTab === 'posts' ? `${styles.tab} ${styles.tabActive}` : styles.tab}
                  onClick={() => setActiveTab('posts')}
                >
                  Постове
                </button>
                <span>|</span>
                <button
                  type="button"
                  className={activeTab === 'tagged' ? `${styles.tab} ${styles.tabActive}` : styles.tab}
                  onClick={() => setActiveTab('tagged')}
                >
                  Тагнати в постове
                </button>
              </div>
              <button
                type="button"
                className={styles.postsActionButton}
                onClick={() => onNavigate?.('posts', { search: '?compose=1' })}
              >
                <CirclePlus size={18} />
                Добави пост
              </button>
            </div>

            <div className={styles.postsFeed}>
              {activePosts.length ? activePosts.map((post) => (
                <article key={post.id} className={styles.post}>
                  <div className={styles.postHeader}>
                    <div className={styles.postAvatar}>
                      <UserRound size={30} />
                    </div>
                    <div>
                      <strong>{post.author}</strong>
                      <span>{post.time}</span>
                    </div>
                    <button type="button" className={styles.moreButton}>
                      <CircleEllipsis size={20} />
                    </button>
                  </div>

                  <p className={styles.postText}>{post.text}</p>
                  {post.image ? (
                    <img src={post.image} alt={post.author} className={styles.postImage} />
                  ) : null}

                  <div className={styles.postStats}>
                    <span><Heart size={22} /> {post.likes.toLocaleString()}</span>
                    <span><Share2 size={22} /> {post.shares.toLocaleString()}</span>
                    <span><MessageCircle size={22} /> {post.comments}</span>
                  </div>

                  <button type="button" className={styles.seeMore}>Виж още</button>
                </article>
              )) : (
                <p className={styles.statusMessage}>
                  {activeTab === 'posts' ? 'Все още няма публикувани постове.' : 'Няма тагнати постове.'}
                </p>
              )}
            </div>
          </section>

          <aside className={styles.friendsCard}>
            <h3>Приятели</h3>
            <div className={styles.friendList}>
              {friends.length ? friends.map((friend) => (
                <div key={friend.id} className={styles.friendItem}>
                  <div className={styles.friendAvatar}>
                    <UserRound size={22} />
                  </div>
                  <div>
                    <strong>{friend.name}</strong>
                    <span>{friend.email}</span>
                  </div>
                </div>
              )) : (
                <p className={styles.statusMessage}>Все още няма добавени приятели.</p>
              )}
            </div>
          </aside>
        </section>
      </div>

      {isEditOpen ? (
        <div className={styles.modalOverlay} onClick={closeEditModal} role="presentation">
          <section
            className={styles.modalCard}
            onClick={(event) => event.stopPropagation()}
            role="dialog"
            aria-modal="true"
            aria-labelledby="edit-profile-title"
          >
            <div className={styles.modalHeader}>
              <div>
                <h2 id="edit-profile-title">Редактирай профила</h2>
              </div>
              <button type="button" className={styles.iconButton} onClick={closeEditModal} aria-label="Затвори редакцията">
                <X size={20} />
              </button>
            </div>

            <form className={styles.editForm} onSubmit={saveProfile}>
              <div className={styles.formGrid}>
                <label className={styles.field}>
                  <span>Име</span>
                  <input value={draft.name} onChange={handleDraftChange('name')} />
                </label>
                <label className={styles.field}>
                  <span>Имейл</span>
                  <input type="email" value={draft.email} onChange={handleDraftChange('email')} />
                </label>
                <label className={styles.field}>
                  <span>Текуща парола</span>
                  <input type="password" value={draft.oldPassword} onChange={handleDraftChange('oldPassword')} />
                </label>
                {SOCIAL_FIELDS.map((field) => (
                  <label key={field.key} className={styles.field}>
                    <span>{field.label}</span>
                    <input
                      type="url"
                      placeholder={field.placeholder}
                      value={draft.socials[field.key]}
                      onChange={handleSocialChange(field.key)}
                    />
                  </label>
                ))}
              </div>

              <label className={`${styles.field} ${styles.fieldFull}`}>
                <span>Биография</span>
                <textarea rows="6" value={draft.bio} onChange={handleDraftChange('bio')} />
              </label>

              {saveError ? <p className={styles.modalError}>{saveError}</p> : null}

              <div className={styles.modalActions}>
                <button type="button" className={styles.secondaryButton} onClick={closeEditModal}>
                  Отказ
                </button>
                <button type="submit" className={styles.primaryButton} disabled={isSaving}>
                  {isSaving ? 'Запазване...' : 'Запази промените'}
                </button>
              </div>
            </form>
          </section>
        </div>
      ) : null}
    </main>
  );
}
