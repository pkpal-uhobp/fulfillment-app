<script setup>
import { computed, onMounted, ref } from 'vue'
import { RefreshCcw, ShieldCheck, UserRound } from '@lucide/vue'
import { getCurrentUser, loadMe } from '@/shared/api/http'

const user = ref(getCurrentUser())
const loading = ref(false)
const error = ref('')

const initials = computed(() => {
  const name = user.value?.full_name || user.value?.email || 'Администратор'
  return name
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase())
    .join('') || 'АД'
})

const displayName = computed(() => user.value?.full_name || user.value?.email || 'Администратор')

async function refreshProfile() {
  loading.value = true
  error.value = ''
  try {
    user.value = await loadMe()
  } catch (err) {
    user.value = getCurrentUser()
    error.value = err?.message || 'Не удалось обновить профиль.'
  } finally {
    loading.value = false
  }
}

onMounted(refreshProfile)
</script>

<template>
  <section class="admin-profile-page">
    <div class="hero-card">
      <p class="eyebrow">профиль администратора</p>
      <h1>Профиль</h1>
      <p>Данные текущего администратора системы и параметры учетной записи.</p>
    </div>

    <div class="profile-grid">
      <aside class="visual-card">
        <div class="avatar">{{ initials }}</div>
        <strong>{{ displayName }}</strong>
        <span>Администратор системы</span>
      </aside>

      <main class="data-card">
        <div class="data-card__head">
          <div>
            <p class="eyebrow">данные пользователя</p>
            <h2>{{ displayName }}</h2>
          </div>
          <button type="button" :disabled="loading" @click="refreshProfile">
            <RefreshCcw class="h-5 w-5" :class="{ 'animate-spin': loading }" />
            Обновить
          </button>
        </div>

        <p v-if="error" class="error-box">{{ error }}</p>

        <dl class="profile-fields">
          <div>
            <dt>ФИО</dt>
            <dd>{{ user?.full_name || '—' }}</dd>
          </div>
          <div>
            <dt>Email</dt>
            <dd>{{ user?.email || '—' }}</dd>
          </div>
          <div>
            <dt>Телефон</dt>
            <dd>{{ user?.phone || '—' }}</dd>
          </div>
          <div>
            <dt>Роль</dt>
            <dd><ShieldCheck class="inline h-5 w-5 text-[#ff3f4d]" /> Администратор</dd>
          </div>
        </dl>
      </main>
    </div>
  </section>
</template>

<style scoped>
.admin-profile-page {
  display: grid;
  gap: 28px;
}

.hero-card,
.data-card,
.visual-card {
  border-radius: 32px;
  background: #fff;
  box-shadow: 0 24px 70px rgba(6, 17, 38, .08);
}

.hero-card {
  padding: 36px;
}

.eyebrow {
  color: #ff3f4d;
  font-size: 12px;
  font-weight: 950;
  letter-spacing: .42em;
  text-transform: uppercase;
}

.hero-card h1 {
  margin: 12px 0;
  font-size: clamp(44px, 7vw, 82px);
  line-height: .9;
  letter-spacing: -.07em;
}

.hero-card p:last-child {
  color: #64748b;
  font-size: 18px;
  font-weight: 800;
}

.profile-grid {
  display: grid;
  grid-template-columns: minmax(240px, 360px) minmax(0, 1fr);
  gap: 24px;
}

.visual-card {
  min-height: 420px;
  padding: 30px;
  background:
    radial-gradient(circle at 15% 15%, rgba(255, 63, 77, .25), transparent 28%),
    linear-gradient(145deg, #061126, #0b3a48);
  color: #fff;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
}

.avatar {
  width: 96px;
  height: 96px;
  border-radius: 28px;
  display: grid;
  place-items: center;
  background: #ff3f4d;
  font-size: 32px;
  font-weight: 950;
  box-shadow: 0 20px 50px rgba(255, 63, 77, .35);
}

.visual-card strong {
  margin-top: 22px;
  font-size: 28px;
  font-weight: 950;
}

.visual-card span {
  margin-top: 8px;
  color: #cbd5e1;
  font-weight: 900;
}

.data-card {
  padding: 34px;
}

.data-card__head {
  display: flex;
  justify-content: space-between;
  gap: 20px;
  align-items: center;
}

.data-card h2 {
  margin-top: 10px;
  font-size: clamp(34px, 4vw, 56px);
  letter-spacing: -.06em;
}

.data-card button {
  min-height: 56px;
  border: 0;
  border-radius: 18px;
  padding: 0 22px;
  display: inline-flex;
  align-items: center;
  gap: 10px;
  background: #eef3f9;
  color: #061126;
  font-weight: 950;
  cursor: pointer;
}

.profile-fields {
  margin-top: 28px;
  border: 1px solid #d7e1ee;
  border-radius: 24px;
  overflow: hidden;
}

.profile-fields div {
  display: grid;
  grid-template-columns: 210px 1fr;
  border-bottom: 1px solid #d7e1ee;
}

.profile-fields div:last-child {
  border-bottom: 0;
}

.profile-fields dt,
.profile-fields dd {
  padding: 20px;
}

.profile-fields dt {
  color: #94a3b8;
  font-size: 12px;
  font-weight: 950;
  letter-spacing: .32em;
  text-transform: uppercase;
}

.profile-fields dd {
  font-size: 18px;
  font-weight: 950;
}

.error-box {
  margin-top: 20px;
  border-radius: 18px;
  background: #fee2e2;
  color: #b91c1c;
  padding: 14px 18px;
  font-weight: 900;
}

@media (max-width: 900px) {
  .profile-grid {
    grid-template-columns: 1fr;
  }

  .profile-fields div {
    grid-template-columns: 1fr;
  }
}
</style>
