<template>
  <section class="profile-page">
    <article class="profile-card">
      <div class="head">
        <div>
          <p>Аккаунт</p>
          <h2>Данные пользователя</h2>
        </div>
        <button type="button" @click="load">Обновить</button>
      </div>

      <div class="info-list">
        <div class="info-row"><span>ФИО</span><b>{{ user?.full_name || '—' }}</b></div>
        <div class="info-row"><span>Email</span><b>{{ user?.email || '—' }}</b></div>
        <div class="info-row"><span>Телефон</span><b>{{ user?.phone || '—' }}</b></div>
        <div class="info-row"><span>Роль</span><b>{{ user?.role === 'admin' ? 'Администратор' : 'Логист' }}</b></div>
      </div>
    </article>
  </section>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { getCurrentUser, loadMe } from '@/shared/api/http'

const user = ref(getCurrentUser())

async function load() {
  try {
    user.value = await loadMe()
  } catch {
    user.value = getCurrentUser()
  }
}

onMounted(load)
</script>

<style scoped>
.profile-page{display:grid;gap:20px}.profile-card{background:white;border-radius:34px;padding:32px;box-shadow:0 18px 42px rgba(7,16,31,.08);max-width:920px}.head{display:flex;justify-content:space-between;gap:18px;align-items:start;margin-bottom:24px}.head p{margin:0 0 8px;color:#ff3f4d;text-transform:uppercase;letter-spacing:.22em;font-size:12px;font-weight:900}.head h2{margin:0;font-size:38px;letter-spacing:-.05em}.head button{height:48px;border:0;border-radius:16px;padding:0 18px;background:#edf2f7;color:#07101f;font-weight:900;cursor:pointer}.info-list{display:grid;gap:14px}.info-row{padding:22px;border-radius:24px;background:#f6f8fb;display:grid;gap:6px}.info-row span{color:#92a0b5;text-transform:uppercase;letter-spacing:.18em;font-weight:900;font-size:12px}.info-row b{font-size:20px}@media(max-width:640px){.head{flex-direction:column}.profile-card{padding:22px}.head h2{font-size:30px}}
</style>
