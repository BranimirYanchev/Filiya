import React, { useEffect, useMemo, useRef, useState } from 'react';
import {
  BookOpen,
  ChevronLeft,
  ChevronRight,
  Search,
  Ellipsis,
  ExternalLink,
  Flame,
  FileText,
  Heart,
  Home,
  Landmark,
  Compass,
  Image,
  List,
  MessageCircle,
  Paperclip,
  Smile,
  Shapes,
  Repeat2,
  X,
  UsersRound
} from 'lucide-react';
import styles from './Posts.module.scss';
import { FEED_WEIGHTS } from '../../recommendation/feedConfig';
import {
  createPost,
  fetchCategories,
  fetchCurrentUser,
  fetchPosts,
  logoutUser
} from '../../auth/session';
import {
  buildSpecificQualificationKey,
  getSavedCirclePreferences,
  PERIODS,
  QUALIFICATION_GROUPS
} from './categoryFilters';

const VALID_PERIODS = new Set(PERIODS);
const QUALIFICATION_ITEMS_BY_GROUP = new Map(
  QUALIFICATION_GROUPS.map((group) => [group.name, new Set(group.items)])
);
const QUALIFICATION_GROUP_NAMES = new Set(QUALIFICATION_GROUPS.map((group) => group.name));

const EMPTY_FILTERS = {
  periods: [],
  generalQualifications: [],
  specificQualifications: []
};

const CATEGORY_LABELS = {
  Philosophy: 'Философия',
  Literature: 'Литература',
  History: 'История',
  Art: 'Изкуство',
  'Humanities and Social Sciences': 'Хуманитарни и социални науки',
  Antiquity: 'Античност',
  'Middle Ages': 'Средновековие',
  Renaissance: 'Ренесанс',
  Baroque: 'Барок',
  Enlightenment: 'Просвещение',
  'XIX century': 'XIX век',
  Modernism: 'Модернизъм',
  Vanguard: 'Авангард',
  Postmodernity: 'Постмодерност',
  Modernity: 'Съвременност',
  'A history of ideas': 'История на идеите',
  Ethics: 'Етика',
  Metaphysics: 'Метафизика',
  'Political philosophy': 'Политическа философия',
  Aesthetics: 'Естетика',
  Poetry: 'Поезия',
  Prose: 'Проза',
  Drama: 'Драма',
  Essays: 'Есеистика',
  'Political history': 'Политическа история',
  'Social history': 'Социална история',
  'Cultural history': 'Културна история',
  'Military history': 'Военна история',
  'Visual arts': 'Визуални изкуства',
  Musics: 'Музика',
  Theater: 'Театър',
  Cinema: 'Кино',
  'Architecture and design': 'Архитектура и дизайн',
  Anthropology: 'Антропология',
  Sociology: 'Социология',
  Psychology: 'Психология',
  Linguistics: 'Лингвистика',
  'Political science': 'Политология'
};

const sameSelection = (left, right) => {
  if (left.length !== right.length) {
    return false;
  }

  const compare = new Set(left);
  return right.every((item) => compare.has(item));
};

const FRIENDS = [
  { name: 'Eddie Lobanovskiy', email: 'lobanovskiy@gmail.com' },
  { name: 'Alexey Stave', email: 'alexeyts@gmail.com' },
  { name: 'Anton Tkacheve', email: 'tkacheveanton@gmail.com' },
  { name: 'Maria Gonzalez', email: 'maria.gonzalez@email.com' },
  { name: 'Jin Park', email: 'jin.park@email.com' },
  { name: 'Liam O\'Brien', email: 'liam.obrien@email.com' },
  { name: 'Aisha Khan', email: 'aisha.khan@email.com' },
  { name: 'Ravi Patel', email: 'ravi.patel@email.com' },
  { name: 'Sofia Rossi', email: 'sofia.rossi@email.com' },
  { name: 'David Kim', email: 'david.kim@email.com' }
];

const takeUnique = (source, count, selectedIds) => {
  const result = [];
  for (const item of source) {
    if (result.length >= count) break;
    if (selectedIds.has(item.id)) continue;
    selectedIds.add(item.id);
    result.push(item);
  }
  return result;
};

const shuffle = (arr) => [...arr].sort(() => Math.random() - 0.5);

const getInitialFilterState = () => {
  try {
    const saved = getSavedCirclePreferences();
    if (!saved) {
      return EMPTY_FILTERS;
    }

    const normalizedPeriods = Array.isArray(saved.periods)
      ? [...new Set(saved.periods.filter((value) => VALID_PERIODS.has(value)))]
      : [];

    const normalizedGeneral = Array.isArray(saved.generalQualifications)
      ? [...new Set(saved.generalQualifications.filter((value) => QUALIFICATION_GROUP_NAMES.has(value)))]
      : [];

    const normalizedSpecific = Array.isArray(saved.specificQualifications)
      ? [
          ...new Set(
            saved.specificQualifications.filter((value) => {
              const [groupName, subName] = value.split('::');
              return (
                typeof groupName === 'string'
                && typeof subName === 'string'
                && QUALIFICATION_GROUP_NAMES.has(groupName)
                && QUALIFICATION_ITEMS_BY_GROUP.get(groupName)?.has(subName)
              );
            })
          )
        ]
      : [];

    return {
      periods: normalizedPeriods,
      generalQualifications: normalizedGeneral,
      specificQualifications: normalizedSpecific
    };
  } catch (error) {
    return EMPTY_FILTERS;
  }
};

const translateCategoryName = (name) => CATEGORY_LABELS[name] || name || 'Без категория';

const getAttachmentUrl = (attachment) => `https://filiya-backend.onrender.com${attachment.file_url}`;

const countImageAttachments = (attachments = []) => attachments.filter((attachment) => attachment.kind === 'image').length;
const getImageAttachments = (attachments = []) => attachments.filter((attachment) => attachment.kind === 'image');

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

const mapApiPost = (post) => ({
  id: post.id,
  author: post?.author?.full_name || post?.author?.email || 'Потребител',
  handle: post?.author?.email ? `@${post.author.email.split('@')[0]}` : '',
  time: formatRelativeTime(post?.created_at),
  category: translateCategoryName(post?.categories?.[0]?.name),
  title: post?.title || '',
  text: post?.content || '',
  likes: Array.isArray(post?.likes) ? post.likes.length : 0,
  shares: Array.isArray(post?.attachments) ? post.attachments.length : 0,
  comments: Array.isArray(post?.comments) ? post.comments.length : 0,
  isFriend: false,
  isNew: true,
  periodizations: Array.isArray(post?.categories) ? post.categories.map((category) => translateCategoryName(category.name)) : [],
  qualifications: Array.isArray(post?.categories) ? post.categories.map((category) => translateCategoryName(category.name)) : [],
  attachments: Array.isArray(post?.attachments) ? post.attachments : []
});

export default function Posts({ onNavigate }) {
  const [searchText, setSearchText] = useState('');
  const [previewImage, setPreviewImage] = useState(null);
  const [selectedPost, setSelectedPost] = useState(null);
  const [composerOpen, setComposerOpen] = useState(false);
  const [composerTitle, setComposerTitle] = useState('');
  const [composerText, setComposerText] = useState('');
  const [composerFiles, setComposerFiles] = useState([]);
  const [composerCategoryIds, setComposerCategoryIds] = useState([]);
  const [composerError, setComposerError] = useState('');
  const [composerSubmitting, setComposerSubmitting] = useState(false);
  const [profileMenuOpen, setProfileMenuOpen] = useState(false);
  const [feedPosts, setFeedPosts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [currentUser, setCurrentUser] = useState(null);
  const initialFilters = useMemo(() => getInitialFilterState(), []);
  const [selectedPeriods, setSelectedPeriods] = useState(initialFilters.periods);
  const [selectedGeneralQualifications, setSelectedGeneralQualifications] = useState(initialFilters.generalQualifications);
  const [selectedSpecificQualifications, setSelectedSpecificQualifications] = useState(initialFilters.specificQualifications);

  const [draftPeriods, setDraftPeriods] = useState(initialFilters.periods);
  const [draftGeneralQualifications, setDraftGeneralQualifications] = useState(initialFilters.generalQualifications);
  const [draftSpecificQualifications, setDraftSpecificQualifications] = useState(initialFilters.specificQualifications);
  const [activeFilterMenu, setActiveFilterMenu] = useState('periods');
  const [feedMode, setFeedMode] = useState('default');
  const imageInputRef = useRef(null);
  const fileInputRef = useRef(null);

  useEffect(() => {
    let isMounted = true;

    const loadData = async () => {
      try {
        const [postsResponse, categoriesResponse, userResponse] = await Promise.all([
          fetchPosts({ limit: 50, offset: 0 }),
          fetchCategories({ limit: 200, offset: 0 }),
          fetchCurrentUser().catch(() => null)
        ]);

        if (!isMounted) return;

        setFeedPosts(Array.isArray(postsResponse) ? postsResponse.map(mapApiPost) : []);
        setCategories(Array.isArray(categoriesResponse) ? categoriesResponse : []);
        setCurrentUser(userResponse);
      } catch (error) {
        if (!isMounted) return;
      }
    };

    loadData();

    return () => {
      isMounted = false;
    };
  }, []);

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    if (params.get('compose') !== '1') {
      return;
    }

    setComposerOpen(true);

    params.delete('compose');
    const nextSearch = params.toString();
    const nextUrl = `${window.location.pathname}${nextSearch ? `?${nextSearch}` : ''}`;
    window.history.replaceState({}, '', nextUrl);
  }, []);

  const hasPendingChanges = useMemo(() => (
    !sameSelection(selectedPeriods, draftPeriods)
    || !sameSelection(selectedGeneralQualifications, draftGeneralQualifications)
    || !sameSelection(selectedSpecificQualifications, draftSpecificQualifications)
  ), [selectedPeriods, selectedGeneralQualifications, selectedSpecificQualifications, draftPeriods, draftGeneralQualifications, draftSpecificQualifications]);

  const hasActiveFilters = selectedPeriods.length > 0 || selectedGeneralQualifications.length > 0 || selectedSpecificQualifications.length > 0;

  const toggleDraftPeriod = (periodName) => {
    setDraftPeriods((prev) => (
      prev.includes(periodName)
        ? prev.filter((item) => item !== periodName)
        : [...prev, periodName]
    ));
  };

  const toggleDraftGeneralQualification = (groupName) => {
    setDraftGeneralQualifications((prev) => {
      const isSelected = prev.includes(groupName);
      if (!isSelected) {
        return [...prev, groupName];
      }

      setDraftSpecificQualifications((specificPrev) => (
        specificPrev.filter((item) => !item.startsWith(`${groupName}::`))
      ));
      return prev.filter((item) => item !== groupName);
    });
  };

  const toggleDraftSpecificQualification = (groupName, itemName) => {
    const key = buildSpecificQualificationKey(groupName, itemName);

    setDraftSpecificQualifications((prev) => (
      prev.includes(key) ? prev.filter((item) => item !== key) : [...prev, key]
    ));

    setDraftGeneralQualifications((prev) => (
      prev.includes(groupName) ? prev : [...prev, groupName]
    ));
  };

  const clearFilters = () => {
    setDraftPeriods([]);
    setDraftGeneralQualifications([]);
    setDraftSpecificQualifications([]);
    setSelectedPeriods([]);
    setSelectedGeneralQualifications([]);
    setSelectedSpecificQualifications([]);
  };

  const applyFilters = () => {
    if (!hasPendingChanges) {
      return;
    }

    setSelectedPeriods(draftPeriods);
    setSelectedGeneralQualifications(draftGeneralQualifications);
    setSelectedSpecificQualifications(draftSpecificQualifications);
  };

  const displayedPosts = useMemo(() => {
    const q = searchText.trim().toLowerCase();
    const hasAnyQualificationFilter = selectedGeneralQualifications.length > 0 || selectedSpecificQualifications.length > 0;

    const base = feedPosts.filter((post) => {
      const matchesQuery = q.length === 0
        || post.title.toLowerCase().includes(q)
        || post.text.toLowerCase().includes(q)
        || post.author.toLowerCase().includes(q)
        || post.category.toLowerCase().includes(q);

      if (!matchesQuery) {
        return false;
      }

      const postQualificationSet = new Set(post.qualifications || []);
      const hasPeriod = !selectedPeriods.length || (post.periodizations || []).some((period) => selectedPeriods.includes(period));

      const hasQualification = !hasAnyQualificationFilter
        ? true
        : selectedGeneralQualifications.some((groupName) => postQualificationSet.has(groupName))
          || selectedSpecificQualifications.some((key) => postQualificationSet.has(key.split('::')[1]));

      return hasPeriod && hasQualification;
    });

    if (feedMode === 'popular') {
      return [...base].sort((left, right) => {
        const leftScore = left.likes + left.shares + left.comments;
        const rightScore = right.likes + right.shares + right.comments;
        return rightScore - leftScore;
      });
    }

    const interestPool = shuffle(base.filter((post) => post.isFriend));
    const freshPool = shuffle(base.filter((post) => post.isNew));
    const randomPool = shuffle(base);

    const total = Math.min(base.length, 8);
    const interestCount = Math.round(total * (FEED_WEIGHTS.interest / 100));
    const freshCount = Math.round(total * (FEED_WEIGHTS.fresh / 100));
    const randomCount = Math.max(0, total - interestCount - freshCount);

    const selectedIds = new Set();
    const fromInterest = takeUnique(interestPool, interestCount, selectedIds);
    const fromFresh = takeUnique(freshPool, freshCount, selectedIds);
    const fromRandom = takeUnique(randomPool, randomCount, selectedIds);

    return [...fromInterest, ...fromFresh, ...fromRandom];
  }, [feedMode, feedPosts, searchText, selectedPeriods, selectedGeneralQualifications, selectedSpecificQualifications]);

  const resetComposer = () => {
    setComposerOpen(false);
    setComposerTitle('');
    setComposerText('');
    setComposerFiles([]);
    setComposerCategoryIds([]);
    setComposerError('');
    setComposerSubmitting(false);
  };

  const toggleComposerCategory = (categoryId) => {
    setComposerCategoryIds((prev) => (
      prev.includes(categoryId)
        ? prev.filter((item) => item !== categoryId)
        : [...prev, categoryId]
    ));
  };

  const handleComposerFiles = (fileList) => {
    const files = Array.from(fileList || []);
    setComposerFiles((prev) => [...prev, ...files]);
  };

  const removeComposerFile = (indexToRemove) => {
    setComposerFiles((prev) => prev.filter((_, index) => index !== indexToRemove));
  };

  const openPostDetails = (post) => {
    setSelectedPost(post);
  };

  const closePostDetails = () => {
    setSelectedPost(null);
  };

  const openImagePreview = (attachments, startIndex) => {
    const images = getImageAttachments(attachments);
    if (!images.length) return;

    setPreviewImage({
      images,
      index: Math.max(0, Math.min(startIndex, images.length - 1))
    });
  };

  const closeImagePreview = () => setPreviewImage(null);

  const showPreviousPreviewImage = () => {
    if (!previewImage?.images?.length) return;
    setPreviewImage((current) => ({
      ...current,
      index: current.index === 0 ? current.images.length - 1 : current.index - 1
    }));
  };

  const showNextPreviewImage = () => {
    if (!previewImage?.images?.length) return;
    setPreviewImage((current) => ({
      ...current,
      index: current.index === current.images.length - 1 ? 0 : current.index + 1
    }));
  };

  const submitComposer = async () => {
    if (!composerText.trim()) {
      setComposerError('Добави съдържание за поста.');
      return;
    }

    if (!composerCategoryIds.length) {
      setComposerError('Избери поне една категория.');
      return;
    }

    setComposerSubmitting(true);
    setComposerError('');

    try {
      await createPost({
        title: composerTitle.trim() || composerText.trim().slice(0, 80),
        content: composerText.trim(),
        categoryIds: composerCategoryIds,
        attachments: composerFiles
      });

      const refreshedPosts = await fetchPosts({ limit: 50, offset: 0 });
      setFeedPosts(Array.isArray(refreshedPosts) ? refreshedPosts.map(mapApiPost) : []);

      resetComposer();
    } catch (error) {
      setComposerError(error.message || 'Неуспешно публикуване на поста.');
      setComposerSubmitting(false);
    }
  };

  return (
    <div className={styles.page}>
      <main className={styles.shell}>
        <section className={styles.board}>
          <aside className={styles.leftPanel}>
            <div className={styles.categories}>
              <div className={styles.primaryMenu}>
                <button
                  type="button"
                  className={styles.primaryMenuItem}
                  onClick={() => onNavigate?.('home')}
                >
                  <span className={styles.filterMenuLabel}>
                    <Home size={18} />
                    Начало
                  </span>
                </button>
                <button
                  type="button"
                  className={feedMode === 'default' ? `${styles.primaryMenuItem} ${styles.filterMenuItemActive}` : styles.primaryMenuItem}
                  onClick={() => {
                    setFeedMode('default');
                    setActiveFilterMenu('qualifications');
                  }}
                >
                  <span className={styles.filterMenuLabel}>
                    <Compass size={18} />
                    Открий
                  </span>
                </button>
                <button
                  type="button"
                  className={feedMode === 'popular' ? `${styles.primaryMenuItem} ${styles.filterMenuItemActive}` : styles.primaryMenuItem}
                  onClick={() => setFeedMode('popular')}
                >
                  <span className={styles.filterMenuLabel}>
                    <Flame size={18} />
                    Най-популярни
                  </span>
                </button>
              </div>

              <div className={styles.filterHeader}>
                <strong>Филтри</strong>
                {hasActiveFilters ? (
                  <button type="button" className={styles.clearButton} onClick={clearFilters}>
                    Изчисти
                  </button>
                ) : null}
              </div>
              <div className={styles.filterMenu}>
                <button
                  type="button"
                  className={activeFilterMenu === 'periods' ? `${styles.filterMenuItem} ${styles.filterMenuItemActive}` : styles.filterMenuItem}
                  onClick={() => setActiveFilterMenu('periods')}
                >
                  <span className={styles.filterMenuLabel}>
                    <BookOpen size={18} />
                    Периодизация
                  </span>
                  <span className={styles.filterMenuMeta}>
                    {draftPeriods.length > 0 ? `${draftPeriods.length}` : ''}
                  </span>
                </button>
                <button
                  type="button"
                  className={activeFilterMenu === 'qualifications' ? `${styles.filterMenuItem} ${styles.filterMenuItemActive}` : styles.filterMenuItem}
                  onClick={() => setActiveFilterMenu('qualifications')}
                >
                  <span className={styles.filterMenuLabel}>
                    <Shapes size={18} />
                    Квалификация
                  </span>
                  <span className={styles.filterMenuMeta}>
                    {draftGeneralQualifications.length + draftSpecificQualifications.length > 0
                      ? `${draftGeneralQualifications.length + draftSpecificQualifications.length}`
                      : ''}
                  </span>
                </button>
              </div>

              {activeFilterMenu === 'periods' ? (
                <section className={styles.dropdownPanel}>
                  <div className={styles.dropdownTitle}>
                    <span>Периодизация</span>
                    <Landmark size={16} />
                  </div>
                  <div className={styles.dropdownContent}>
                    <div className={styles.filterBlock}>
                      {PERIODS.map((period) => (
                        <label key={period} className={styles.filterItem}>
                          <input
                            type="checkbox"
                            checked={draftPeriods.includes(period)}
                            onChange={() => toggleDraftPeriod(period)}
                          />
                          <span>{period}</span>
                        </label>
                      ))}
                    </div>
                  </div>
                </section>
              ) : null}

              {activeFilterMenu === 'qualifications' ? (
                <section className={styles.dropdownPanel}>
                  <div className={styles.dropdownTitle}>
                    <span>Квалификация</span>
                    <Shapes size={16} />
                  </div>
                  <div className={styles.dropdownContent}>
                    <div className={styles.subGroups}>
                      {QUALIFICATION_GROUPS.map((group) => {
                        const hasOpenSubfilters =
                          draftGeneralQualifications.includes(group.name)
                          || draftSpecificQualifications.some((key) => key.startsWith(`${group.name}::`));

                        const subFiltersClassName = hasOpenSubfilters
                          ? `${styles.subFilters} ${styles.subFiltersOpen}`
                          : styles.subFilters;

                        return (
                          <article key={group.name} className={styles.qualGroup}>
                            <label className={`${styles.filterItem} ${styles.qualGroupToggle}`}>
                              <input
                                type="checkbox"
                                checked={draftGeneralQualifications.includes(group.name)}
                                onChange={() => toggleDraftGeneralQualification(group.name)}
                              />
                              <span>{group.name}</span>
                              <span
                                className={
                                  hasOpenSubfilters
                                    ? `${styles.qualGroupArrow} ${styles.qualGroupArrowOpen}`
                                    : styles.qualGroupArrow
                                }
                                aria-hidden="true"
                              >
                                ▾
                              </span>
                            </label>
                            <div className={subFiltersClassName}>
                              <div className={styles.filterBlock}>
                                {group.items.map((item) => (
                                  <label key={`${group.name}-${item}`} className={styles.filterItem}>
                                    <input
                                      type="checkbox"
                                      checked={draftSpecificQualifications.includes(buildSpecificQualificationKey(group.name, item))}
                                      onChange={() => toggleDraftSpecificQualification(group.name, item)}
                                    />
                                    <span>{item}</span>
                                  </label>
                                ))}
                              </div>
                            </div>
                          </article>
                        );
                      })}
                    </div>
                  </div>
                </section>
              ) : null}

              <div className={styles.applySection}>
                <button
                  type="button"
                  className={styles.applyButton}
                  onClick={applyFilters}
                  disabled={!hasPendingChanges}
                >
                  Приложи
                </button>
              </div>
            </div>
          </aside>

          <section className={styles.feedColumn}>
            <div className={styles.searchBar}>
              <Search size={24} />
              <input
                type="text"
                value={searchText}
                onChange={(event) => setSearchText(event.target.value)}
                placeholder="Търсене"
              />
            </div>

            <article className={styles.createPostCard}>
              <div className={styles.avatar} aria-hidden="true">👨🏽</div>
              <div className={styles.createBody}>
                <button
                  type="button"
                  className={styles.createTrigger}
                  onClick={() => setComposerOpen(true)}
                >
                  Какво се случва?
                </button>
                <div className={styles.createActions}>
                  <div className={styles.createIcons}>
                    <Image size={18} />
                    <List size={18} />
                    <MessageCircle size={18} />
                    <FileText size={18} />
                  </div>
                  <button type="button" onClick={() => setComposerOpen(true)}>Публикувай</button>
                </div>
              </div>
            </article>

            {displayedPosts.length > 0 ? (
              displayedPosts.map((post) => (
                <article key={post.id} className={styles.postCard}>
                  <header className={styles.postHeader}>
                    <div className={styles.avatar} aria-hidden="true">👨🏽</div>
                    <div className={styles.postMeta}>
                      <div className={styles.authorRow}>
                        <strong
                          className={styles.authorLink}
                          role="button"
                          tabIndex={0}
                          onClick={() => onNavigate?.('profile')}
                          onKeyDown={(event) => {
                            if (event.key === 'Enter' || event.key === ' ') {
                              onNavigate?.('profile');
                            }
                          }}
                        >
                          {post.author}
                        </strong>
                        <UsersRound size={18} />
                      </div>
                      <span>{post.time}</span>
                    </div>
                    <button type="button" className={styles.moreBtn}>
                      <Ellipsis size={18} />
                    </button>
                  </header>

                  <p className={styles.postTitle}>{post.title}</p>
                  <p className={styles.postSubtitle}>
                    {post.category}
                    {post.periodizations?.[0] ? ` • ${post.periodizations[0]}` : ''}
                  </p>
                  {post.attachments?.length ? (
                    <button
                      type="button"
                      className={styles.attachmentHint}
                      onClick={() => openPostDetails(post)}
                    >
                      <Paperclip size={14} />
                      <span>
                        {post.attachments.length} прикачени файла
                        {countImageAttachments(post.attachments) ? ` • ${countImageAttachments(post.attachments)} изображения` : ''}
                      </span>
                    </button>
                  ) : null}
                  <p className={styles.postText}>{post.text}</p>

                  <footer className={styles.postFooter}>
                    <div>
                      <Heart size={19} /> {post.likes.toLocaleString()}
                    </div>
                    <div>
                      <Repeat2 size={19} /> {post.shares.toLocaleString()}
                    </div>
                    <div>
                      <MessageCircle size={19} /> {post.comments}
                    </div>
                  </footer>

                  <button type="button" className={styles.seeMore} onClick={() => openPostDetails(post)}>Виж още</button>
                </article>
              ))
            ) : (
              <article className={styles.emptyState}>
                <p>Няма публикации за избраните филтри.</p>
              </article>
            )}
          </section>

          <aside className={styles.rightPanel}>
            <article className={styles.profileCard}>
              <div className={styles.avatar} aria-hidden="true">👨🏽</div>
              <div>
                <strong
                  className={styles.authorLink}
                  role="button"
                  tabIndex={0}
                  onClick={() => onNavigate?.('profile')}
                  onKeyDown={(event) => {
                    if (event.key === 'Enter' || event.key === ' ') {
                      onNavigate?.('profile');
                    }
                  }}
                >
                  {currentUser?.full_name || currentUser?.email || 'Профил'}
                </strong>
                <p>{currentUser?.email ? `@${currentUser.email.split('@')[0]}` : ''}</p>
              </div>
              <div className={styles.profileMenuWrap}>
                <button
                  type="button"
                  className={styles.moreBtn}
                  onClick={() => setProfileMenuOpen((value) => !value)}
                >
                  <Ellipsis size={18} />
                </button>

                {profileMenuOpen ? (
                  <div className={styles.profileDropdown}>
                    <button
                      type="button"
                      onClick={() => {
                        setProfileMenuOpen(false);
                        onNavigate?.('profile');
                      }}
                    >
                      Профил
                    </button>
                    <button
                      type="button"
                      onClick={async () => {
                        setProfileMenuOpen(false);
                        await logoutUser();
                        onNavigate?.('login');
                      }}
                    >
                      Изход
                    </button>
                  </div>
                ) : null}
              </div>
            </article>

            <article className={styles.friendsCard}>
              <h3>Приятели</h3>
              <ul>
                {FRIENDS.map((friend) => (
                  <li key={friend.email}>
                    <div className={styles.friendAvatar} aria-hidden="true">👤</div>
                    <div>
                      <strong>{friend.name}</strong>
                      <p>{friend.email}</p>
                    </div>
                  </li>
                ))}
              </ul>
            </article>
          </aside>
        </section>
      </main>

      {composerOpen ? (
        <div className={styles.composerOverlay} onClick={resetComposer}>
          <div className={styles.composerModal} onClick={(event) => event.stopPropagation()}>
            <div className={styles.composerHeader}>
              <button
                type="button"
                className={styles.composerClose}
                onClick={resetComposer}
              >
                <X size={22} />
              </button>
              <button type="button" className={styles.composerDraftsButton}>
                Чернови
              </button>
            </div>

            <div className={styles.composerBody}>
              <div className={styles.avatar} aria-hidden="true">👨🏽</div>
              <div className={styles.composerFields}>
                <input
                  type="text"
                  value={composerTitle}
                  onChange={(event) => setComposerTitle(event.target.value)}
                  placeholder="Заглавие на поста"
                  className={styles.composerTitleInput}
                />
                <textarea
                  value={composerText}
                  onChange={(event) => setComposerText(event.target.value)}
                  placeholder="Какво се случва?"
                  className={styles.composerTextarea}
                />
                <div className={styles.categoryChooser}>
                  <strong>Категории</strong>
                  <div className={styles.categoryOptions}>
                    {categories.map((category) => (
                      <label key={category.id} className={styles.categoryPill}>
                        <input
                          type="checkbox"
                          checked={composerCategoryIds.includes(category.id)}
                          onChange={() => toggleComposerCategory(category.id)}
                        />
                        <span>{translateCategoryName(category.name)}</span>
                      </label>
                    ))}
                  </div>
                </div>
                {composerFiles.length ? (
                  <div className={styles.composerAttachments}>
                    {composerFiles.map((file, index) => (
                      <button
                        key={`${file.name}-${index}`}
                        type="button"
                        className={styles.attachmentChip}
                        onClick={() => removeComposerFile(index)}
                      >
                        <FileText size={16} />
                        <span>{file.name}</span>
                        <X size={14} />
                      </button>
                    ))}
                  </div>
                ) : null}
                {composerError ? <p className={styles.composerError}>{composerError}</p> : null}
              </div>
            </div>

            <div className={styles.composerFooter}>
              <div className={styles.createIcons}>
                <button type="button" className={styles.composerIconButton} onClick={() => imageInputRef.current?.click()}>
                  <Image size={18} />
                </button>
                <List size={18} />
                <Smile size={18} />
                <button type="button" className={styles.composerIconButton} onClick={() => fileInputRef.current?.click()}>
                  <FileText size={18} />
                </button>
              </div>
              <input
                ref={imageInputRef}
                type="file"
                accept="image/*"
                multiple
                hidden
                onChange={(event) => handleComposerFiles(event.target.files)}
              />
              <input
                ref={fileInputRef}
                type="file"
                accept=".pdf,.doc,.docx,.txt,.rtf,.odt,.xls,.xlsx,.ppt,.pptx"
                multiple
                hidden
                onChange={(event) => handleComposerFiles(event.target.files)}
              />
              <button
                type="button"
                className={styles.composerPostButton}
                onClick={submitComposer}
                disabled={composerSubmitting}
              >
                {composerSubmitting ? 'Публикуване...' : 'Публикувай'}
              </button>
            </div>
          </div>
        </div>
      ) : null}

      {previewImage ? (
        <div className={styles.imagePreviewOverlay} onClick={closeImagePreview}>
          <div className={styles.imagePreviewModal} onClick={(event) => event.stopPropagation()}>
            <button
              type="button"
              className={styles.imagePreviewClose}
              onClick={closeImagePreview}
            >
              <X size={22} />
            </button>
            {previewImage.images.length > 1 ? (
              <button
                type="button"
                className={`${styles.imagePreviewArrow} ${styles.imagePreviewArrowLeft}`}
                onClick={showPreviousPreviewImage}
                aria-label="Предишна снимка"
              >
                <ChevronLeft size={22} />
              </button>
            ) : null}
            <img
              src={getAttachmentUrl(previewImage.images[previewImage.index])}
              alt={previewImage.images[previewImage.index]?.file_name}
              className={styles.imagePreviewFull}
            />
            {previewImage.images.length > 1 ? (
              <button
                type="button"
                className={`${styles.imagePreviewArrow} ${styles.imagePreviewArrowRight}`}
                onClick={showNextPreviewImage}
                aria-label="Следваща снимка"
              >
                <ChevronRight size={22} />
              </button>
            ) : null}
          </div>
        </div>
      ) : null}

      {selectedPost ? (
        <div className={styles.postDetailOverlay} onClick={closePostDetails}>
          <div className={styles.postDetailModal} onClick={(event) => event.stopPropagation()}>
            <button
              type="button"
              className={styles.imagePreviewClose}
              onClick={closePostDetails}
            >
              <X size={22} />
            </button>

            <div className={styles.postDetailContent}>
              <div className={styles.postDetailHeader}>
                <div>
                  <p className={styles.postDetailEyebrow}>{selectedPost.category}</p>
                  <h2>{selectedPost.title}</h2>
                  <p className={styles.postDetailMeta}>{selectedPost.time}</p>
                </div>
              </div>

              <p className={styles.postDetailText}>{selectedPost.text}</p>

              {selectedPost.attachments?.length ? (
                <section className={styles.postDetailAttachments}>
                  <div className={styles.postDetailAttachmentsHeader}>
                    <strong>Прикачени файлове</strong>
                    <span>{selectedPost.attachments.length}</span>
                  </div>

                  {selectedPost.attachments.some((attachment) => attachment.kind !== 'image') ? (
                    <div className={styles.attachmentSubsection}>
                      <div className={styles.attachmentSubsectionTitle}>Документи</div>
                      <div className={styles.documentChipRow}>
                        {selectedPost.attachments
                          .filter((attachment) => attachment.kind !== 'image')
                          .map((attachment) => (
                            <a
                              key={attachment.id}
                              href={getAttachmentUrl(attachment)}
                              target="_blank"
                              rel="noreferrer"
                              className={styles.documentChip}
                            >
                              <FileText size={16} />
                              <span>{attachment.file_name}</span>
                              <ExternalLink size={14} />
                            </a>
                          ))}
                      </div>
                    </div>
                  ) : null}

                  {selectedPost.attachments.some((attachment) => attachment.kind === 'image') ? (
                    <div className={styles.attachmentSubsection}>
                      <div className={styles.attachmentSubsectionTitle}>Снимки</div>
                      <div className={styles.imageThumbGrid}>
                        {getImageAttachments(selectedPost.attachments).map((attachment, index) => (
                          <button
                            key={attachment.id}
                            type="button"
                            className={styles.imageThumbCard}
                            onClick={() => openImagePreview(selectedPost.attachments, index)}
                          >
                            <img
                              src={getAttachmentUrl(attachment)}
                              alt={attachment.file_name}
                              className={styles.imageThumb}
                            />
                          </button>
                        ))}
                      </div>
                    </div>
                  ) : null}
                </section>
              ) : null}
            </div>
          </div>
        </div>
      ) : null}
    </div>
  );
}
