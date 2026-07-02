<template>
  <div class="auth-page">
    <div class="auth-page__left">
      <div class="auth-page__brand">
        <div class="auth-page__logo">N</div>
        <h1>ApiNest</h1>
        <p>开源 API 协作管理平台</p>
      </div>
    </div>

    <div class="auth-page__right">
      <div class="auth-panel">
        <router-view />
      </div>
    </div>
  </div>
</template>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(420px, 520px);
}

.auth-page__left {
  position: relative;
  background:
    radial-gradient(circle at 20% 20%, rgba(37, 99, 235, 0.28), transparent 40%),
    radial-gradient(circle at 80% 70%, rgba(124, 58, 237, 0.22), transparent 35%),
    linear-gradient(135deg, #0f172a, #1e293b 55%, #0f172a);
  color: #e2e8f0;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px;
}

.auth-page__brand {
  max-width: 420px;
}

.auth-page__logo {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  background: linear-gradient(135deg, #2563eb, #3b82f6);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 20px;
}

.auth-page__brand h1 {
  margin: 0 0 8px;
  font-size: 36px;
  font-weight: 700;
}

.auth-page__brand p {
  margin: 0;
  color: #94a3b8;
  font-size: 16px;
}

.auth-page__right {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px 24px;
  background: var(--color-surface);
}

.auth-panel {
  width: 100%;
  max-width: 380px;
}

@media (max-width: 900px) {
  .auth-page {
    grid-template-columns: 1fr;
  }

  .auth-page__left {
    display: none;
  }
}
</style>
