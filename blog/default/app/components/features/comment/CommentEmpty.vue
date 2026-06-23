<template>
  <div class="comments-empty">
    <i class="ri-chat-off-line" />
    <p>暂无评论，快来发表第一条评论吧！</p>
  </div>
</template>

<style lang="scss" scoped>
.comments-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: var(--flec-lightText);
  text-align: center;

  i {
    font-size: 4rem;
    margin-bottom: 15px;
    opacity: 0.5;
  }

  p {
    font-size: 1rem;
    margin: 0;
  }
}

// 响应式设�?
@media screen and (max-width: 768px) {
  .comments-empty {
    padding: 40px 20px;

    i {
      font-size: 3rem;
    }
  }
}
</style>
