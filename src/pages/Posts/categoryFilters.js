export const CIRCLE_CATEGORIES_STORAGE_KEY = 'circles-preferences-v1';

export const PERIODS = [
  'Античност',
  'Средновековие',
  'Ренесанс',
  'Барок',
  'Просвещение',
  'XIX век',
  'Модернизъм',
  'Авангард',
  'Постмодерност',
  'Съвременност'
];

export const QUALIFICATION_GROUPS = [
  {
    name: 'Философия',
    items: ['История на идеите', 'Етика', 'Метафизика', 'Политическа философия', 'Естетика', 'Други']
  },
  {
    name: 'Литература',
    items: ['Поезия', 'Проза', 'Драма', 'Есеистика']
  },
  {
    name: 'История',
    items: ['Политическа история', 'Социална история', 'Културна история', 'Военна история']
  },
  {
    name: 'Изкуства',
    items: ['Визуални изкуства', 'Музика', 'Театър', 'Кино', 'Архитектура и дизайн', 'Други']
  },
  {
    name: 'Хуманитарни и социални науки',
    items: ['Антропология', 'Социология', 'Психология', 'Лингвистика', 'Политология']
  }
];

export const buildSpecificQualificationKey = (groupName, itemName) => `${groupName}::${itemName}`;

export const getSavedCirclePreferences = () => {
  try {
    const raw = window.localStorage.getItem(CIRCLE_CATEGORIES_STORAGE_KEY);
    if (!raw) {
      return {
        periods: [],
        generalQualifications: [],
        specificQualifications: []
      };
    }

    const parsed = JSON.parse(raw);
    return {
      periods: Array.isArray(parsed.periods) ? parsed.periods : [],
      generalQualifications: Array.isArray(parsed.generalQualifications) ? parsed.generalQualifications : [],
      specificQualifications: Array.isArray(parsed.specificQualifications) ? parsed.specificQualifications : []
    };
  } catch (error) {
    return {
      periods: [],
      generalQualifications: [],
      specificQualifications: []
    };
  }
};
