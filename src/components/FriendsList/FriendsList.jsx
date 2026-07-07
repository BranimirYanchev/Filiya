import React from 'react';
import { Link } from 'react-router-dom';
import styles from './FriendsList.module.scss';
import { resolveApiUrl } from '../../utils/api';

const DEFAULT_AVATAR = 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=50&h=50&fit=crop';

export default function FriendsList({ user }) {
  const friends = Array.isArray(user?.friends) ? user.friends : [];

  return (
    <div className={styles.friendsList}>
      <h2 className={styles.heading}>Приятели</h2>
      <div className={styles.friendsContainer}>
        {friends.length > 0 ? (
          friends.map((friend) => {
            const friendId = friend.id || friend.ID;
            const friendName = friend.full_name || friend.FullName || friend.email || 'Потребител';
            const friendEmail = friend.email || 'Няма имейл';
            const avatarUrl = resolveApiUrl(friend.profile_picture || friend.ProfilePicture || DEFAULT_AVATAR);

            return (
              <Link key={friendId || friendEmail} to={friendId ? `/profile/${friendId}` : '#'} className={styles.friendItem}>
                <img
                  src={avatarUrl}
                  alt={friendName}
                  className={styles.friendAvatar}
                />
                <div className={styles.friendInfo}>
                  <span className={styles.friendName}>{friendName}</span>
                  <span className={styles.friendEmail}>{friendEmail}</span>
                </div>
              </Link>
            );
          })
        ) : (
          <div className={styles.emptyState}>Все още няма видим списък с приятели за този профил.</div>
        )}
      </div>
    </div>
  );
}
