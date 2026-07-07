import React from 'react';
import PostShowcase from '../PostShowcase/PostShowcase';

export default function ImageCarousel({ posts, loading, error }) {
  return (
    <PostShowcase
      title="Открий"
      subtitle="Показваме реалните постове от общността. Зареждат се на порции по 4, така че секцията да се държи като жив feed, а не като статичен блок."
      posts={posts}
      loading={loading}
      error={error}
      emptyMessage="Все още няма постове за показване."
    />
  );
}
