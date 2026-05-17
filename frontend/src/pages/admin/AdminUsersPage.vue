<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { apiFetch } from '@/shared/api/http'

const users = ref([])
const loading = ref(false)
const saving = ref(false)
const error = ref('')
const notice = ref('')
const filters = reactive({ role: 'all', search: '' })
const form = reactive({
  email: '',
  password: '',
  full_name: '',
  phone: '',
  role: 'client',
})

const roles = [
  { value: 'client', label: 'Клиент' },
  { value: 'logist', label: 'Логист' },
  { value: 'worker', label: 'Рабочий' },
  { value: 'admin', label: 'Администратор' },
]

const roleLabels = Object.fromEntries(roles.map((role) => [role.value, role.label]))

function collection(payload) {
  if (Array.isArray(payload)) return payload

  return payload?.users || payload?.items || payload?.data || []
}

const filteredUsers = computed(() => {
  const query = filters.search.trim().toLowerCase()

  return users.value.filter((user) => {
    const roleOk = filters.role === 'all' || user.role === filters.role
    const text = [user.full_name, user.email, user.phone, user.role, user.id].join(' ').toLowerCase()

    return roleOk && (!query || text.includes(query))
  })
})

const counters = computed(() => ({
  total: users.value.length,
  active: users.value.filter((user) => user.is_active && !user.is_blocked).length,
  blocked: users.value.filter((user) => user.is_blocked).length,
  admins: users.value.filter((user) => user.role === 'admin').length,
}))

function resetForm() {
  form.email = ''
  form.password = ''
  form.full_name = ''
  form.phone = ''
  form.role = 'client'
}

async function loadUsers() {
  loading.value = true
  error.value = ''

  try {
    const payload = await apiFetch('/users?limit=500', { auth: true })
    users.value = collection(payload)
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить пользователей'
  } finally {
    loading.value = false
  }
}

async function createUser() {
  saving.value = true
  error.value = ''
  notice.value = ''

  try {
    await apiFetch('/users', {
      method: 'POST',
      auth: true,
      body: {
        email: form.email,
        password: form.password,
        full_name: form.full_name,
        phone: form.phone || undefined,
        role: form.role,
      },
    })

    notice.value = 'Пользователь создан'
    resetForm()
    await loadUsers()
  } catch (err) {
    error.value = err.message || 'Не удалось создать пользователя'
  } finally {
    saving.value = false
  }
}

async function patchUser(user, patch) {
  error.value = ''
  notice.value = ''

  try {
    await apiFetch(`/users/${user.id}`, {
      method: 'PATCH',
      auth: true,
      body: patch,
    })

    notice.value = 'Пользователь обновлён'
    await loadUsers()
  } catch (err) {
    error.value = err.message || 'Не удалось обновить пользователя'
  }
}

async function deleteUser(user) {
  if (!confirm(`Удалить пользователя ${user.email}?`)) return

  error.value = ''
  notice.value = ''

  try {
    await apiFetch(`/users/${user.id}`, {
      method: 'DELETE',
      auth: true,
    })

    notice.value = 'Пользователь удалён'
    await loadUsers()
  } catch (err) {
    error.value = err.message || 'Не удалось удалить пользователя'
  }
}

onMounted(loadUsers)
</script>

<template>
  <section class="admin-page">
    <header class="hero-card">
      <div>
        <p class="eyebrow">Пользователи</p>
        <h1>Аккаунты и роли</h1>
        <span>Создавайте сотрудников, назначайте роли и блокируйте доступ без изменения кода.</span>
      </div>

      <button class="dark-btn" type="button" :disabled="loading" @click="loadUsers">
        {{ loading ? 'Загрузка…' : 'Обновить' }}
      </button>
    </header>

    <div v-if="error" class="alert error">{{ error }}</div>
    <div v-if="notice" class="alert success">{{ notice }}</div>

    <section class="stats-grid">
      <article>
        <span>Всего</span>
        <strong>{{ counters.total }}</strong>
      </article>

      <article>
        <span>Активных</span>
        <strong>{{ counters.active }}</strong>
      </article>

      <article>
        <span>Заблокировано</span>
        <strong>{{ counters.blocked }}</strong>
      </article>

      <article>
        <span>Админов</span>
        <strong>{{ counters.admins }}</strong>
      </article>
    </section>

    <section class="workspace-grid">
      <form class="panel-card form-card" @submit.prevent="createUser">
        <p class="eyebrow">Новый аккаунт</p>
        <h2>Создать пользователя</h2>

        <label>
          <span>ФИО</span>
          <input v-model.trim="form.full_name" required type="text" placeholder="Иванов Иван" />
        </label>

        <label>
          <span>Email</span>
          <input v-model.trim="form.email" required type="email" placeholder="user@example.com" />
        </label>

        <label>
          <span>Телефон</span>
          <input v-model.trim="form.phone" type="tel" placeholder="+7..." />
        </label>

        <label>
          <span>Пароль</span>
          <input v-model="form.password" required minlength="6" type="password" placeholder="минимум 6 символов" />
        </label>

        <label>
          <span>Роль</span>
          <select v-model="form.role">
            <option v-for="role in roles" :key="role.value" :value="role.value">
              {{ role.label }}
            </option>
          </select>
        </label>

        <button class="red-btn" type="submit" :disabled="saving">
          {{ saving ? 'Создаём…' : 'Создать' }}
        </button>
      </form>

      <section class="panel-card list-card">
        <div class="list-head">
          <div>
            <p class="eyebrow">Список</p>
            <h2>{{ filteredUsers.length }} пользователей</h2>
          </div>
        </div>

        <div class="filters-row">
          <label>
            <span>Роль</span>
            <select v-model="filters.role">
              <option value="all">Все роли</option>
              <option v-for="role in roles" :key="role.value" :value="role.value">
                {{ role.label }}
              </option>
            </select>
          </label>

          <label>
            <span>Поиск</span>
            <input v-model.trim="filters.search" type="text" placeholder="ФИО, email, телефон" />
          </label>
        </div>

        <div v-if="!filteredUsers.length" class="empty">Пользователи не найдены.</div>

        <div v-else class="users-list">
          <article v-for="user in filteredUsers" :key="user.id" class="user-row">
            <div>
              <strong>{{ user.full_name || 'Без имени' }}</strong>
              <span>{{ user.email }}</span>
              <small v-if="user.phone">{{ user.phone }}</small>
            </div>

            <em>{{ roleLabels[user.role] || user.role }}</em>

            <select :value="user.role" @change="patchUser(user, { role: $event.target.value })">
              <option v-for="role in roles" :key="role.value" :value="role.value">
                {{ role.label }}
              </option>
            </select>

            <div class="user-actions">
              <button
                type="button"
                class="small-btn"
                :class="{ danger: !user.is_blocked }"
                @click="patchUser(user, { is_blocked: !user.is_blocked })"
              >
                {{ user.is_blocked ? 'Разблокировать' : 'Заблокировать' }}
              </button>

              <button
                type="button"
                class="small-btn"
                @click="patchUser(user, { is_active: !user.is_active })"
              >
                {{ user.is_active ? 'Деактивировать' : 'Активировать' }}
              </button>

              <button type="button" class="small-btn danger" @click="deleteUser(user)">
                Удалить
              </button>
            </div>
          </article>
        </div>
      </section>
    </section>
  </section>
</template>

<style scoped>
.admin-page {
  display: grid;
  gap: 26px;
}

.hero-card,
.panel-card,
.stats-grid article {
  background: #fff;
  border-radius: 34px;
  box-shadow: 0 18px 62px rgba(15, 23, 42, .08);
}

.hero-card {
  padding: 34px;
  display: flex;
  justify-content: space-between;
  gap: 24px;
}

.eyebrow {
  margin: 0 0 12px;
  color: #ff3f4d;
  font-size: 13px;
  font-weight: 950;
  letter-spacing: .28em;
  text-transform: uppercase;
}

h1,
h2 {
  margin: 0;
  color: #061126;
  font-weight: 950;
  letter-spacing: -.06em;
}

h1 {
  font-size: clamp(48px, 7vw, 84px);
  line-height: .9;
}

h2 {
  font-size: clamp(28px, 3vw, 42px);
  line-height: 1;
}

.hero-card span {
  display: block;
  margin-top: 14px;
  color: #5d6d83;
  font-size: 18px;
  font-weight: 800;
  line-height: 1.5;
}

.dark-btn,
.red-btn,
.small-btn {
  border: 0;
  border-radius: 20px;
  font-weight: 950;
  cursor: pointer;
}

.dark-btn,
.red-btn {
  min-height: 58px;
  padding: 0 24px;
  color: #fff;
  font-size: 16px;
}

.dark-btn {
  background: #061126;
}

.red-btn {
  background: #ff3f4d;
  box-shadow: 0 18px 42px rgba(255, 63, 77, .22);
}

.alert,
.empty {
  padding: 18px 22px;
  border-radius: 22px;
  font-weight: 900;
}

.alert.error {
  background: #fff0f1;
  color: #be123c;
}

.alert.success {
  background: #e8fff5;
  color: #047857;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 18px;
}

.stats-grid article {
  min-height: 130px;
  padding: 24px;
  display: grid;
  align-content: center;
}

.stats-grid span,
label span {
  color: #97a5bb;
  font-size: 13px;
  font-weight: 950;
  letter-spacing: .18em;
  text-transform: uppercase;
}

.stats-grid strong {
  margin-top: 10px;
  font-size: 48px;
  line-height: 1;
  font-weight: 950;
}

.workspace-grid {
  display: grid;
  grid-template-columns: minmax(320px, .62fr) minmax(0, 1.38fr);
  gap: 22px;
  align-items: start;
}

.panel-card {
  padding: 30px;
}

.form-card,
.filters-row {
  display: grid;
  gap: 14px;
}

label {
  display: grid;
  gap: 8px;
}

input,
select {
  width: 100%;
  min-height: 54px;
  border: 1px solid #dbe4ef;
  border-radius: 18px;
  background: #f8fbff;
  color: #061126;
  padding: 0 16px;
  font-size: 16px;
  font-weight: 850;
  outline: none;
  box-sizing: border-box;
}

input:focus,
select:focus {
  border-color: #ff3f4d;
  box-shadow: 0 0 0 5px rgba(255, 63, 77, .12);
  background: #fff;
}

.list-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 20px;
}

.filters-row {
  grid-template-columns: 220px minmax(0, 1fr);
  margin-bottom: 18px;
}

.users-list {
  display: grid;
  gap: 12px;
  max-height: 620px;
  overflow-y: auto;
  overflow-x: hidden;
  padding-right: 6px;
}

.user-row {
  border-radius: 24px;
  background: #f6f9fd;
  padding: 16px;
  display: grid;
  grid-template-columns: minmax(220px, 1fr) auto minmax(170px, 220px) minmax(260px, auto);
  gap: 12px;
  align-items: center;
  min-width: 0;
}

.user-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  flex-wrap: wrap;
  min-width: 0;
}

.user-row strong {
  display: block;
  font-size: 17px;
  font-weight: 950;
}

.user-row span,
.user-row small {
  display: block;
  margin-top: 4px;
  color: #66758a;
  font-weight: 800;
  overflow-wrap: anywhere;
}

em {
  border-radius: 999px;
  padding: 10px 14px;
  background: #dbeafe;
  color: #1d4ed8;
  font-style: normal;
  font-weight: 950;
  white-space: nowrap;
}

.small-btn {
  min-height: 46px;
  padding: 0 14px;
  background: #061126;
  color: #fff;
}

.small-btn.danger {
  background: #ffe4e6;
  color: #be123c;
}

.empty {
  background: #f6f9fd;
  color: #64748b;
}

@media (max-width: 1280px) {
  .workspace-grid,
  .stats-grid,
  .filters-row {
    grid-template-columns: 1fr;
  }

  .user-row {
    grid-template-columns: 1fr 1fr;
  }

  .user-actions {
    grid-column: 1 / -1;
    justify-content: flex-start;
  }
}

@media (max-width: 760px) {
  .hero-card {
    flex-direction: column;
    align-items: stretch;
  }

  .user-row {
    grid-template-columns: 1fr;
  }

  .user-actions {
    grid-column: auto;
  }
}
</style>
