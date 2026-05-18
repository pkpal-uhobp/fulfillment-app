<script setup>
import { computed, onMounted, reactive, ref } from 'vue'

import { apiFetch } from '@/shared/api/http'

const users = ref([])
const loading = ref(false)
const saving = ref(false)
const editing = ref(false)
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

const editForm = reactive({
  id: null,
  email: '',
  password: '',
  full_name: '',
  phone: '',
  role: 'client',
  is_active: true,
  is_blocked: false,
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

function openEdit(user) {
  error.value = ''
  notice.value = ''
  editForm.id = user.id
  editForm.email = user.email || ''
  editForm.password = ''
  editForm.full_name = user.full_name || ''
  editForm.phone = user.phone || ''
  editForm.role = user.role || 'client'
  editForm.is_active = Boolean(user.is_active)
  editForm.is_blocked = Boolean(user.is_blocked)
  editing.value = true
}

function closeEdit() {
  editing.value = false
  editForm.id = null
  editForm.password = ''
}

function editPayload() {
  const payload = {
    email: editForm.email,
    full_name: editForm.full_name,
    phone: editForm.phone || '',
    role: editForm.role,
    is_active: Boolean(editForm.is_active),
    is_blocked: Boolean(editForm.is_blocked),
  }

  const password = editForm.password.trim()
  if (password) payload.password = password

  return payload
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

async function saveUser() {
  if (!editForm.id) return

  saving.value = true
  error.value = ''
  notice.value = ''

  try {
    await apiFetch(`/users/${editForm.id}`, {
      method: 'PATCH',
      auth: true,
      body: editPayload(),
    })

    notice.value = 'Данные пользователя сохранены'
    closeEdit()
    await loadUsers()
  } catch (err) {
    error.value = err.message || 'Не удалось сохранить пользователя'
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

async function deactivateUser(user) {
  if (!confirm(`Деактивировать пользователя ${user.email}?`)) return

  error.value = ''
  notice.value = ''

  try {
    await apiFetch(`/users/${user.id}`, {
      method: 'DELETE',
      auth: true,
    })

    notice.value = 'Пользователь деактивирован'
    await loadUsers()
  } catch (err) {
    error.value = err.message || 'Не удалось деактивировать пользователя'
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
        <span>Создавайте пользователей и полностью редактируйте ФИО, email, телефон, пароль, роль и статус доступа.</span>
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

        <div class="form-grid">
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
        </div>

        <button class="red-btn" type="submit" :disabled="saving">
          {{ saving ? 'Сохраняем…' : 'Создать' }}
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
            <div class="user-info">
              <strong>{{ user.full_name || 'Без имени' }}</strong>
              <span>{{ user.email }}</span>
              <small v-if="user.phone">{{ user.phone }}</small>
            </div>

            <div class="badges">
              <em>{{ roleLabels[user.role] || user.role }}</em>
              <b :class="{ off: !user.is_active || user.is_blocked }">
                {{ user.is_blocked ? 'Заблокирован' : user.is_active ? 'Активен' : 'Неактивен' }}
              </b>
            </div>

            <div class="user-actions">
              <button type="button" class="small-btn primary" @click="openEdit(user)">
                Редактировать
              </button>

              <button
                type="button"
                class="small-btn"
                :class="{ danger: !user.is_blocked }"
                @click="patchUser(user, { is_blocked: !user.is_blocked })"
              >
                {{ user.is_blocked ? 'Разблокировать' : 'Заблокировать' }}
              </button>

              <button type="button" class="small-btn danger" @click="deactivateUser(user)">
                Деактивировать
              </button>
            </div>
          </article>
        </div>
      </section>
    </section>

    <teleport to="body">
      <div v-if="editing" class="modal-backdrop" @click.self="closeEdit">
        <form class="edit-modal" @submit.prevent="saveUser">
          <div class="modal-head">
            <div>
              <p class="eyebrow">Редактирование</p>
              <h2>Все данные пользователя</h2>
            </div>
            <button type="button" class="icon-btn" @click="closeEdit">×</button>
          </div>

          <div class="edit-grid">
            <label>
              <span>ID</span>
              <input :value="editForm.id" disabled type="text" />
            </label>

            <label>
              <span>Email</span>
              <input v-model.trim="editForm.email" required type="email" placeholder="user@example.com" />
            </label>

            <label>
              <span>ФИО</span>
              <input v-model.trim="editForm.full_name" required type="text" placeholder="Иванов Иван" />
            </label>

            <label>
              <span>Телефон</span>
              <input v-model.trim="editForm.phone" type="tel" placeholder="+7..." />
            </label>

            <label>
              <span>Новый пароль</span>
              <input v-model="editForm.password" minlength="6" type="password" placeholder="оставь пустым, если не менять" />
            </label>

            <label>
              <span>Роль</span>
              <select v-model="editForm.role">
                <option v-for="role in roles" :key="role.value" :value="role.value">
                  {{ role.label }}
                </option>
              </select>
            </label>
          </div>

          <div class="switches">
            <label class="switch-card">
              <input v-model="editForm.is_active" type="checkbox" />
              <span></span>
              <strong>Активный аккаунт</strong>
            </label>

            <label class="switch-card">
              <input v-model="editForm.is_blocked" type="checkbox" />
              <span></span>
              <strong>Заблокирован</strong>
            </label>
          </div>

          <div class="modal-actions">
            <button class="ghost-btn" type="button" @click="closeEdit">Отмена</button>
            <button class="red-btn" type="submit" :disabled="saving">
              {{ saving ? 'Сохраняем…' : 'Сохранить всё' }}
            </button>
          </div>
        </form>
      </div>
    </teleport>
  </section>
</template>

<style scoped>
.admin-page {
  display: grid;
  gap: 26px;
}

.hero-card,
.panel-card,
.stats-grid article,
.edit-modal {
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
.ghost-btn,
.small-btn,
.icon-btn {
  border: 0;
  border-radius: 20px;
  font-weight: 950;
  cursor: pointer;
}

.dark-btn,
.red-btn,
.ghost-btn {
  min-height: 58px;
  padding: 0 24px;
  font-size: 16px;
}

.dark-btn,
.small-btn,
.icon-btn {
  background: #061126;
  color: #fff;
}

.ghost-btn {
  background: #edf3fb;
  color: #061126;
}

.red-btn {
  background: #ff3f4d;
  color: #fff;
  box-shadow: 0 16px 32px rgba(255, 63, 77, .22);
}

.dark-btn:disabled,
.red-btn:disabled {
  opacity: .65;
  cursor: wait;
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
  background: #ecfdf3;
  color: #166534;
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
  grid-template-columns: 1fr;
  gap: 22px;
  align-items: stretch;
}

.panel-card {
  padding: 30px;
}

.form-card,
.filters-row,
.form-grid,
.edit-grid {
  display: grid;
  gap: 14px;
}

.form-grid,
.edit-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.form-card .red-btn {
  justify-self: start;
  min-width: 260px;
  margin-top: 4px;
}

label {
  display: grid;
  gap: 8px;
}

input,
select {
  width: 100%;
  min-height: 56px;
  border: 1px solid #dbe4ef;
  border-radius: 20px;
  background: #f8fbff;
  color: #061126;
  padding: 0 18px;
  font-size: 16px;
  font-weight: 850;
  outline: none;
  box-sizing: border-box;
}

input:disabled {
  opacity: .7;
  cursor: not-allowed;
}

input:focus,
select:focus {
  border-color: #ff3f4d;
  box-shadow: 0 0 0 5px rgba(255, 63, 77, .12);
  background: #fff;
}

select {
  appearance: none;
  -webkit-appearance: none;
  cursor: pointer;
  padding-right: 54px;
  background-color: #f8fbff;
  background-image:
    linear-gradient(45deg, transparent 50%, #061126 50%),
    linear-gradient(135deg, #061126 50%, transparent 50%),
    linear-gradient(135deg, rgba(255, 63, 77, .10), rgba(219, 234, 254, .55));
  background-position:
    calc(100% - 24px) 50%,
    calc(100% - 17px) 50%,
    100% 0;
  background-size: 7px 7px, 7px 7px, 56px 100%;
  background-repeat: no-repeat;
  transition: border-color .18s ease, box-shadow .18s ease, background-color .18s ease, transform .18s ease;
}

select:hover {
  border-color: rgba(255, 63, 77, .45);
  background-color: #fff;
  box-shadow: 0 10px 28px rgba(15, 23, 42, .08);
}

select option {
  color: #061126;
  background: #fff;
  font-weight: 850;
}

.list-card {
  overflow: hidden;
}

.list-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 20px;
}

.filters-row {
  grid-template-columns: 260px minmax(0, 1fr);
  margin-bottom: 18px;
}

.users-list {
  display: grid;
  gap: 12px;
  max-height: 720px;
  overflow-y: auto;
  overflow-x: hidden;
  padding-right: 10px;
}

.user-row {
  border-radius: 24px;
  background: #f6f9fd;
  padding: 18px;
  display: grid;
  grid-template-columns: minmax(260px, 1.4fr) minmax(200px, .8fr) minmax(360px, auto);
  gap: 14px;
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

.badges {
  display: grid;
  gap: 8px;
  justify-items: start;
  align-content: center;
}

em,
.badges b {
  border-radius: 999px;
  padding: 10px 14px;
  font-style: normal;
  font-weight: 950;
  white-space: nowrap;
}

em {
  background: #dbeafe;
  color: #1d4ed8;
}

.badges b {
  background: #dcfce7;
  color: #166534;
}

.badges b.off {
  background: #ffe4e6;
  color: #be123c;
}

.small-btn {
  min-height: 46px;
  padding: 0 14px;
}

.small-btn.primary {
  background: #2563eb;
}

.small-btn.danger {
  background: #ffe4e6;
  color: #be123c;
}

.empty {
  background: #f6f9fd;
  color: #64748b;
}

.modal-backdrop {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: grid;
  place-items: center;
  padding: 28px;
  background: rgba(6, 17, 38, .46);
  backdrop-filter: blur(8px);
}

.edit-modal {
  width: min(920px, 100%);
  max-height: calc(100vh - 56px);
  overflow-y: auto;
  padding: 30px;
}

.modal-head,
.modal-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.modal-head {
  margin-bottom: 22px;
}

.icon-btn {
  width: 52px;
  height: 52px;
  border-radius: 18px;
  font-size: 28px;
  line-height: 1;
}

.switches {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin-top: 16px;
}

.switch-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border-radius: 22px;
  background: #f6f9fd;
  cursor: pointer;
}

.switch-card input {
  display: none;
}

.switch-card span {
  width: 54px;
  height: 32px;
  border-radius: 999px;
  background: #dbe4ef;
  position: relative;
  flex: 0 0 auto;
}

.switch-card span::after {
  content: '';
  position: absolute;
  top: 5px;
  left: 5px;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: #fff;
  box-shadow: 0 2px 8px rgba(15, 23, 42, .2);
  transition: transform .18s ease;
}

.switch-card input:checked + span {
  background: #ff3f4d;
}

.switch-card input:checked + span::after {
  transform: translateX(22px);
}

.switch-card strong {
  color: #061126;
  font-weight: 950;
}

.modal-actions {
  margin-top: 22px;
  justify-content: flex-end;
}

@media (max-width: 1280px) {
  .stats-grid,
  .filters-row,
  .form-grid,
  .edit-grid,
  .switches {
    grid-template-columns: 1fr;
  }

  .user-row {
    grid-template-columns: 1fr;
  }

  .user-actions {
    justify-content: flex-start;
  }
}

@media (max-width: 760px) {
  .hero-card,
  .modal-head,
  .modal-actions {
    flex-direction: column;
    align-items: stretch;
  }

  .form-card .red-btn {
    width: 100%;
  }
}
</style>
