import React from 'react';
import ProfileShowcase from '../ProfileShowcase/ProfileShowcase';

export default function PostGrid({ profiles, loading, error }) {
  return (
    <ProfileShowcase
      title="Най-популярни"
      subtitle="Тук профилите са подредени по реална активност: брой постове, харесвания и коментари по техните публикации."
      profiles={profiles}
      loading={loading}
      error={error}
      emptyMessage="Все още няма достатъчно активни профили за класацията."
      variant="popular"
    />
  );
}
