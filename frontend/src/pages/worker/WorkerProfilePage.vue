<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { clearAuth, getCurrentUser, loadMe } from '@/shared/api/http'

const router = useRouter()
const user = ref(getCurrentUser())
const loading = ref(false)
const error = ref('')

const initials = computed(() => {
  const source = user.value?.full_name || user.value?.email || 'Рабочий склада'
  return source.split(/\s+/).filter(Boolean).map((part) => part[0]).join('').slice(0, 2).toUpperCase() || 'РС'
})

const rows = computed(() => [
  { label: 'ФИО', value: user.value?.full_name || '—', icon: '◉' },
  { label: 'Email', value: user.value?.email || '—', icon: '✉' },
  { label: 'Телефон', value: user.value?.phone || '—', icon: '☎' },
  { label: 'Роль', value: 'Рабочий склада', icon: '▦' },
])

async function refresh() {
  loading.value = true
  error.value = ''
  try {
    user.value = await loadMe()
  } catch (err) {
    error.value = err.message || 'Не удалось обновить профиль'
    user.value = getCurrentUser()
  } finally {
    loading.value = false
  }
}

function logout() {
  clearAuth()
  router.push({ name: 'login' })
}

onMounted(refresh)
</script>

<template>
  <section class="profile-page">
    <div class="profile-hero">
      <div class="visual-card">
        <div class="visual-top">
          <span class="logo-box">FT</span>
          <div><strong>Fulfillment Transit</strong><small>Рабочая смена</small></div>
        </div>
        <div class="warehouse-illustration" aria-hidden="true">
          <span class="box box-a"></span>
          <span class="box box-b"></span>
          <span class="belt"></span>
          <span class="qr">▦</span>
          <span class="truck"></span>
        </div>
      </div>

      <div class="profile-info">
        <p class="eyebrow">Профиль</p>
        <h1>Рабочий склада</h1>
        <span>Данные аккаунта используются для фиксации операций с грузовыми местами и истории статусов.</span>
      </div>
    </div>

    <div v-if="error" class="alert error">{{ error }}</div>

    <section class="profile-card">
      <div class="profile-head">
        <div class="avatar">{{ initials }}</div>
        <div>
          <p class="eyebrow">Данные пользователя</p>
          <h2>{{ user?.full_name || 'Рабочий склада' }}</h2>
        </div>
        <button type="button" class="refresh-btn" :disabled="loading" @click="refresh">{{ loading ? 'Обновляем…' : 'Обновить' }}</button>
      </div>

      <div class="info-rows">
        <article v-for="row in rows" :key="row.label" class="info-row">
          <i>{{ row.icon }}</i>
          <span>{{ row.label }}</span>
          <strong>{{ row.value }}</strong>
        </article>
      </div>

      <div class="logout-strip">
        <div>
          <p class="eyebrow">Сессия</p>
          <span>После выхода потребуется снова войти по email и паролю.</span>
        </div>
        <button type="button" class="logout-btn" @click="logout">Выйти</button>
      </div>
    </section>
  </section>
</template>

<style scoped>
.profile-page { display: grid; gap: 26px; color: #061126; }
.profile-hero { border-radius: 34px; background: linear-gradient(135deg, #071222, #123247); color: #fff; padding: 32px; display: grid; grid-template-columns: minmax(320px, .9fr) minmax(0, 1.1fr); gap: 28px; box-shadow: 0 18px 62px rgba(15, 23, 42, .16); }
.visual-card { border: 1px solid rgba(255,255,255,.12); border-radius: 30px; padding: 24px; background: rgba(255,255,255,.06); }
.visual-top { display: flex; align-items: center; gap: 14px; }
.logo-box { width: 56px; height: 56px; border-radius: 18px; background: #ff3f4d; display: grid; place-items: center; font-weight: 950; box-shadow: 0 18px 42px rgba(255,63,77,.32); }
.visual-top strong, .visual-top small { display: block; }
.visual-top strong { font-size: 20px; font-weight: 950; }
.visual-top small { margin-top: 4px; color: #ffb3ba; font-size: 12px; font-weight: 950; letter-spacing: .22em; text-transform: uppercase; }
.warehouse-illustration { position: relative; margin-top: 24px; min-height: 330px; border-radius: 28px; background: #081525; overflow: hidden; border: 1px solid rgba(255,255,255,.08); }
.warehouse-illustration::before { content: ''; position: absolute; inset: 34px; border-radius: 24px; background: linear-gradient(135deg, rgba(255,63,77,.16), rgba(20,184,166,.13)); }
.box { position: absolute; width: 118px; height: 86px; border-radius: 20px; background: #202b3d; box-shadow: inset 0 0 0 1px rgba(255,255,255,.08); }
.box::before { content: ''; position: absolute; left: 20px; top: 20px; width: 34px; height: 34px; border-radius: 10px; background: #ff3f4d; }
.box::after { content: ''; position: absolute; left: 64px; top: 24px; width: 38px; height: 10px; border-radius: 999px; background: rgba(255,255,255,.44); box-shadow: 0 22px rgba(255,255,255,.22); }
.box-a { left: 58px; top: 78px; }
.box-b { right: 70px; bottom: 86px; transform: scale(1.08); }
.qr { position: absolute; right: 80px; top: 66px; width: 78px; height: 78px; border-radius: 24px; background: rgba(20,184,166,.16); color: #fff; display: grid; place-items: center; font-size: 34px; font-weight: 950; }
.belt { position: absolute; left: 86px; right: 86px; bottom: 62px; height: 8px; border-radius: 999px; background: rgba(255,255,255,.22); }
.truck { position: absolute; left: 50%; top: 45%; width: 104px; height: 68px; transform: translate(-50%, -50%); border-radius: 26px; background: rgba(255,255,255,.13); border: 1px solid rgba(255,255,255,.14); box-shadow: 0 22px 60px rgba(0,0,0,.22); }
.truck::before { content: '▧'; position: absolute; inset: 0; display: grid; place-items: center; color: #fff; font-size: 40px; }
.eyebrow { margin: 0 0 10px; color: #ff3f4d; font-size: 13px; font-weight: 950; letter-spacing: .28em; text-transform: uppercase; }
.profile-info { align-self: center; }
.profile-info .eyebrow { color: #ff9ca5; }
h1 { margin: 0; max-width: 640px; font-size: clamp(52px, 7vw, 92px); line-height: .88; font-weight: 950; letter-spacing: -.06em; }
.profile-info span { display: block; margin-top: 18px; max-width: 620px; color: #d5e4f3; font-size: 20px; line-height: 1.55; }
.profile-card { border-radius: 34px; background: #fff; padding: 32px; box-shadow: 0 18px 62px rgba(15, 23, 42, .08); }
.profile-head { display: grid; grid-template-columns: auto 1fr auto; align-items: center; gap: 18px; margin-bottom: 24px; }
.avatar { width: 74px; height: 74px; border-radius: 24px; background: #ff3f4d; color: #fff; display: grid; place-items: center; font-size: 24px; font-weight: 950; box-shadow: 0 18px 42px rgba(255, 63, 77, .24); }
h2 { margin: 0; font-size: clamp(32px, 4vw, 54px); line-height: 1; font-weight: 950; letter-spacing: -.05em; }
.refresh-btn, .logout-btn { min-height: 56px; border: 0; border-radius: 18px; padding: 0 22px; font-weight: 950; cursor: pointer; }
.refresh-btn { background: #eef3f9; color: #10223a; }
.logout-btn { background: #ff3f4d; color: #fff; box-shadow: 0 18px 42px rgba(255, 63, 77, .22); }
.info-rows { display: grid; gap: 12px; }
.info-row { min-height: 82px; border-radius: 24px; background: #f6f9fd; padding: 16px 18px; display: grid; grid-template-columns: 42px 160px 1fr; align-items: center; gap: 16px; }
.info-row i { width: 42px; height: 42px; border-radius: 14px; background: #fff; color: #ff3f4d; display: grid; place-items: center; font-style: normal; font-weight: 950; box-shadow: 0 10px 24px rgba(15,23,42,.06); }
.info-row span { color: #94a3b8; font-size: 12px; font-weight: 950; letter-spacing: .22em; text-transform: uppercase; }
.info-row strong { font-size: 20px; font-weight: 950; overflow-wrap: anywhere; }
.logout-strip { margin-top: 22px; border-radius: 26px; background: #061126; color: #fff; padding: 22px; display: flex; align-items: center; justify-content: space-between; gap: 20px; }
.logout-strip .eyebrow { color: #ff9ca5; }
.logout-strip span { color: #d5e4f3; font-size: 16px; font-weight: 800; }
.alert { padding: 18px 22px; border-radius: 22px; font-weight: 900; }
.alert.error { background: #fff0f1; color: #be123c; }
@media (max-width: 1060px) { .profile-hero { grid-template-columns: 1fr; } }
@media (max-width: 760px) { .profile-card, .profile-hero { padding: 20px; border-radius: 28px; } .profile-head, .info-row, .logout-strip { grid-template-columns: 1fr; align-items: start; } .logout-strip { display: grid; } .refresh-btn, .logout-btn { width: 100%; } }
</style>
