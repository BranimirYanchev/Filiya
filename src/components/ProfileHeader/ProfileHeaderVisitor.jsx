import React, { useState, useEffect } from 'react';
import styles from './ProfileHeader.module.scss';
import Button from '../Button/Button';

export default function ProfileHeaderVisitor({ user, friendStatus, onSendFriendRequest, onAcceptFriendRequest }) {
  const profilePicture = user?.profile_picture || user?.ProfilePicture || 
    "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=200&h=200&fit=crop";
  const fullName = user?.full_name || user?.FullName || user?.email || 'User';
  const role = user?.role?.name || user?.Role?.Name || 'Member';
  const [pendingRequestId, setPendingRequestId] = useState(null);

  useEffect(() => {
    // Fetch pending request ID if status is pending_received
    if (friendStatus === 'pending_received') {
      fetchPendingRequestId();
    }
  }, [friendStatus]);

  const fetchPendingRequestId = async () => {
    try {
      const { api } = await import('../../utils/api');
      const response = await api.get('/users/friends/requests/pending');
      if (response && Array.isArray(response)) {
        const userId = user?.id || user?.ID;
        const request = response.find(req => {
          const senderId = req.sender_id || req.SenderID;
          const status = req.status || req.Status;
          return senderId === userId && status === 'pending';
        });
        if (request) {
          setPendingRequestId(request.id || request.ID);
        }
      }
    } catch (err) {
      console.error('Failed to fetch pending request:', err);
    }
  };

  const renderFriendButton = () => {
    switch (friendStatus) {
      case 'friends':
        return (
          <Button className={styles.friendButton} disabled>
            <span className={styles.icon}>✓</span>
            Friends
          </Button>
        );
      case 'pending_sent':
        return (
          <Button className={styles.pendingButton} disabled>
            <span className={styles.icon}>⏳</span>
            Request Sent
          </Button>
        );
      case 'pending_received':
        return (
          <Button 
            className={styles.acceptButton}
            onClick={() => pendingRequestId && onAcceptFriendRequest(pendingRequestId)}
          >
            <span className={styles.icon}>✓</span>
            Accept Request
          </Button>
        );
      default:
        return (
          <Button 
            className={styles.addFriendButton}
            onClick={onSendFriendRequest}
          >
            <span className={styles.icon}>+</span>
            Add Friend
          </Button>
        );
    }
  };

  return (
    <div className={styles.profileHeader}>
      <div className={styles.coverPhoto}>
        <img
          src="/sculp.jpg"
          alt="Cover"
          className={styles.coverImage}
        />
        {/* No edit button for visitors */}
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
          {renderFriendButton()}
        </div>
      </div>
    </div>
  );
}

