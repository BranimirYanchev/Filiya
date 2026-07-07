import React from 'react';
import styles from './About.module.scss';

export default function About({ user }) {
  const email = user?.email || 'Няма имейл';
  const username = user?.username || 'Няма потребителско име';
  const roleName = user?.role?.name || user?.Role?.Name || 'Member';
  const createdAt = user?.created_at
    ? new Intl.DateTimeFormat('bg-BG', { dateStyle: 'medium' }).format(new Date(user.created_at))
    : 'Няма дата';

  const details = [
    { icon: '✉️', value: email, href: email !== 'Няма имейл' ? `mailto:${email}` : null },
    { icon: '@', value: username },
    { icon: '🛡️', value: roleName },
    { icon: '🗓️', value: `Във Филия от ${createdAt}` },
  ];

  return (
    <div className={styles.about}>
      <h2 className={styles.heading}>Профил</h2>
      <div className={styles.contactList}>
        {details.map((item) => (
          item.href ? (
            <a key={item.value} href={item.href} className={styles.contactItem}>
              <span className={styles.icon}>{item.icon}</span>
              <span className={styles.contactText}>{item.value}</span>
            </a>
          ) : (
            <div key={item.value} className={styles.contactItem}>
              <span className={styles.icon}>{item.icon}</span>
              <span className={styles.contactText}>{item.value}</span>
            </div>
          )
        ))}
      </div>
    </div>
  );
}
