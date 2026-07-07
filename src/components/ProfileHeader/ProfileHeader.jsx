import React from 'react';
import styles from './ProfileHeader.module.scss';
import Button from '../Button/Button';

export default function ProfileHeader({ user }) {
  const profilePicture = user?.profile_picture || user?.ProfilePicture || 
    "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=200&h=200&fit=crop";
  const fullName = user?.full_name || user?.FullName || user?.email || 'User';
  const role = user?.role?.name || user?.Role?.Name || 'Member';

  return (
    <div className={styles.profileHeader}>
      <div className={styles.coverPhoto}>
        <img
          src="/sculp.jpg"
          alt="Cover"
          className={styles.coverImage}
        />
        <Button className={styles.editCoverButton}>
          <span className={styles.icon}>✏️</span>
          Edit Cover Photo
        </Button>
      </div>
      <div className={styles.profileInfo}>
        <div className={styles.profilePicture}>
          <img
            src={profilePicture}
            alt={fullName}
          />
        </div>
        <div className={styles.userInfo}>
          <h1 className={styles.name}>{fullName}</h1>
          <p className={styles.profession}>{role}</p>
        </div>
        <div className={styles.actionButtons}>
          <Button className={styles.editButton}>
            <span className={styles.icon}>🎨</span>
            Edit colours
          </Button>
          <Button className={styles.editButton}>
            <span className={styles.icon}>✏️</span>
            Edit Profile
          </Button>
        </div>
      </div>
    </div>
  );
}

