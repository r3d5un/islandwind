<script lang="ts" setup>
import { ref } from 'vue'
import { useLogger } from '@/ui/logging.ts'
import { login, Tokens, useAuthStore } from '@/api/auth.ts'
import { useRouter } from 'vue-router'

const logger = useLogger()
const authStore = useAuthStore()
const route = useRouter()

const username = ref('')
const password = ref('')
const loginError = ref(false)

const handleSubmit = async () => {
  logger.info('Username:', username.value)
  logger.info('Password:', password.value)

  logger.info('logging in', { username: username.value, password: password.value })
  const result = await login(username.value, password.value)
  if (result instanceof Tokens) {
    logger.info('logged in', { tokens: result })
    authStore.loggedIn = true
    authStore.tokens = result
    await route.push('/')
  } else {
    logger.error('Unable to login', { result: result })
    loginError.value = true
  }
}
</script>

<template>
  <div class="login-form">
    <h2>Login</h2>
    <div v-if="loginError" class="error-banner">
      Invalid username or password. Please try again.
    </div>
    <form @submit.prevent="handleSubmit">
      <div class="form-group">
        <label for="username">Username</label>
        <input id="username" v-model="username" type="text" required />
      </div>
      <div class="form-group">
        <label for="password">Password</label>
        <input id="password" v-model="password" type="password" required />
      </div>
      <button type="submit">Login</button>
    </form>
  </div>
</template>

<style scoped>
.login-form {
  font-family: Arial, sans-serif;
  max-width: 400px;
  margin: 0 auto;
  padding: 20px;
  border-radius: 5px;
}

.form-group {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 5px;
}

input {
  width: 100%;
  padding: 8px;
  box-sizing: border-box;
}
.error-banner {
  background-color: #ff4444;
  color: white;
  padding: 10px;
  border-radius: 5px;
  margin-bottom: 20px;
  text-align: center;
}
</style>
