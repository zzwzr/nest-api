import { computed, ref } from 'vue'
import { fetchCurrentUser, login as loginApi, register as registerApi } from '@/api/auth'
import type { LoginPayload, RegisterPayload, UserProfile } from '@/types/auth'
import { clearAuth, getStoredUser, saveAuth } from '@/utils/auth-storage'

const user = ref<UserProfile | null>(getStoredUser())
const bootstrapped = ref(false)

export function useAuth() {
  const isLoggedIn = computed(() => !!user.value)

  async function bootstrap() {
    if (bootstrapped.value) return
    bootstrapped.value = true

    if (!getStoredUser()) return

    try {
      user.value = await fetchCurrentUser()
      saveAuth(
        localStorage.getItem('apinest_access_token') || '',
        localStorage.getItem('apinest_refresh_token') || '',
        user.value,
      )
    } catch {
      clearAuth()
      user.value = null
    }
  }

  async function login(payload: LoginPayload) {
    const result = await loginApi(payload)
    saveAuth(result.access_token, result.refresh_token, result.user)
    user.value = result.user
    return result
  }

  async function register(payload: RegisterPayload) {
    const result = await registerApi(payload)
    saveAuth(result.access_token, result.refresh_token, result.user)
    user.value = result.user
    return result
  }

  function logout() {
    clearAuth()
    user.value = null
  }

  return {
    user,
    isLoggedIn,
    bootstrap,
    login,
    register,
    logout,
  }
}
