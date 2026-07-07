import React from 'react';
import styles from './ProfileVisitor.module.scss';
import ProfileHeaderVisitor from '../../components/ProfileHeader/ProfileHeaderVisitor';
import About from '../../components/About/About';
import Bio from '../../components/Bio/Bio';
import PostsFeed from '../../components/PostsFeed/PostsFeed';
import FriendsList from '../../components/FriendsList/FriendsList';

export default function ProfileVisitor() {
  return (
    <div className={styles.profile}>
      <div className={styles.container}>
        <ProfileHeaderVisitor />
        <div className={styles.mainContent}>
          <div className={styles.leftColumn}>
            <About />
            <Bio />
          </div>
          <div className={styles.centerColumn}>
            <PostsFeed />
          </div>
          <div className={styles.rightColumn}>
            <FriendsList />
          </div>
        </div>
      </div>
    </div>
  );
}

