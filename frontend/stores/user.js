// user.js
import { defineStore } from 'pinia'
import { reactive, watch } from 'vue'

const BASE_URL = 'https://xcxcs.uwdjl.cn'
export const UserStatus = {
  LOGGED_OUT: 'logged_out',
  LOGGED_IN: 'logged_in'
}
const STORAGE_KEY = 'user_store_data'
const SECURE_STORAGE_KEY = 'refresh_token_secure'
const ACCESS_TOKEN_EXPIRES_IN = 30 * 60 // 1800秒
let tokenRefreshTimer = null
let isRefreshing = false

const secureStorage = {
  setRefreshToken(token) {
    try {
      uni.setStorageSync(SECURE_STORAGE_KEY, token)
      uni.setStorageSync(`${SECURE_STORAGE_KEY}_timestamp`, Date.now())
    } catch (error) {
      console.error('Cannot store refreshToken securely:', error)
    }
  },
  getRefreshToken() {
    try {
      const token = uni.getStorageSync(SECURE_STORAGE_KEY)
      const timestamp = uni.getStorageSync(`${SECURE_STORAGE_KEY}_timestamp`)
      if (token && timestamp) {
        const currentTime = Date.now()
        const sevenDaysInMs = 7 * 24 * 60 * 60 * 1000
        if (currentTime - timestamp < sevenDaysInMs) return token
        this.removeRefreshToken()
        return null
      }
      return token
    } catch {
      return null
    }
  },
  removeRefreshToken() {
    try {
      uni.removeStorageSync(SECURE_STORAGE_KEY)
      uni.removeStorageSync(`${SECURE_STORAGE_KEY}_timestamp`)
    } catch {}
  }
}

function getStoredData() {
  try {
    const data = uni.getStorageSync(STORAGE_KEY)
    return data
      ? JSON.parse(data)
      : {
          uid: null,
          isLoggedIn: false,
          token: '',
          avatarUrl: '',
          nickName: '',
          status: UserStatus.LOGGED_OUT,
          tokenExpiry: null,
          registered_days: 0
        }
  } catch {
    return {
      uid: null,
      isLoggedIn: false,
      token: '',
      avatarUrl: '',
      nickName: '',
      status: UserStatus.LOGGED_OUT,
      tokenExpiry: null,
      registered_days: 0
    }
  }
}

export const useUserStore = defineStore('user', () => {
  const user = reactive(getStoredData())

  function saveToStorage() {
    try {
      uni.setStorageSync(STORAGE_KEY, JSON.stringify(user))
    } catch {}
  }

  function watchUser() {
    const keys = [
      'uid',
      'isLoggedIn',
      'token',
      'avatarUrl',
      'nickName',
      'status',
      'tokenExpiry',
      'registered_days'
    ]
    keys.forEach(key => {
      watch(
        () => user[key],
        () => {
          saveToStorage()
        },
        { deep: true }
      )
    })
  }

  function createRequestConfig(config) {
    return {
      ...config,
      header: {
        Authorization: `Bearer ${user.token || ''}`,
        ...(config.header || {})
      }
    }
  }

  function request(config) {
    return new Promise((resolve, reject) => {
      uni.request({
        ...config,
        success: res => {
          if (res.statusCode === 401) {
            if (!isRefreshing) {
              isRefreshing = true
              refreshToken()
                .then(() => {
                  isRefreshing = false
                  resolve(request(config))
                })
                .catch(err => {
                  isRefreshing = false
                  reject(err)
                })
            }
          } else {
            resolve(res)
          }
        },
        fail: err => reject(err)
      })
    })
  }

  async function login() {
    console.log('url:', `${BASE_URL}/users/auth`)
    try {
      const loginRes = await new Promise((resolve, reject) => {
        uni.login({ provider: 'weixin', onlyAuthorize: true, success: resolve, fail: reject })
      })
      if (!loginRes.code) throw new Error('微信登录失败，请重试')


      const authRes = await new Promise((resolve, reject) => {
        uni.request({
          url: `${BASE_URL}/users/auth`,
          method: 'POST',
          data: { code: loginRes.code },
          success: resolve,
          fail: reject
        })
      })

      if (authRes.statusCode !== 200) throw new Error(authRes.data?.message || '登录失败')

      const returnData = authRes.data
      user.token = returnData.access_token
      secureStorage.setRefreshToken(returnData.refresh_token)
      user.uid = returnData.user.id
      user.nickName = returnData.user.nickname
      user.avatarUrl = returnData.user.avatar_url || '/static/images/default_avatar.jpg'
      user.isLoggedIn = true
      user.status = UserStatus.LOGGED_IN
      user.tokenExpiry = Date.now() + ACCESS_TOKEN_EXPIRES_IN * 1000
      user.registered_days = returnData.registered_days || 0

      saveToStorage()
      scheduleTokenRefresh()
      return returnData
    } catch (error) {
      console.error('登录失败:', error)
      throw error
    }
  }

  async function logout() {
    try {
      const response = await request(
        createRequestConfig({
          url: `${BASE_URL}/users/logout`,
          method: 'POST',
          data: { refresh_token: secureStorage.getRefreshToken() }
        })
      )
      reset()
      return response.data
    } catch (error) {
      console.error('登出失败:', error)
      throw error
    }
  }

  async function setNickname(newNickname) {
    try {
      const response = await request(
        createRequestConfig({
          url: `${BASE_URL}/users/set_nickname`,
          method: 'PUT',
          data: { nickname: newNickname }
        })
      )
      user.nickName = response.data.nickname
      return response.data
    } catch (error) {
      console.error('设置昵称失败:', error)
      throw error
    }
  }

  async function setAvatar(avatarFile) {
    try {
      const uploadRes = await new Promise((resolve, reject) => {
        uni.uploadFile({
          url: `${BASE_URL}/users/set_avatar`,
          filePath: avatarFile,
          name: 'avatar',
          header: { Authorization: `Bearer ${user.token || ''}` },
          success: resolve,
          fail: reject
        })
      })
      if (uploadRes.statusCode !== 200) throw new Error(`头像上传失败：${uploadRes.data}`)
      const responseData = JSON.parse(uploadRes.data)
      user.avatarUrl = responseData.avatar_url
      saveToStorage()
      return responseData
    } catch (error) {
      console.error('设置头像失败:', error)
      throw error
    }
  }

  async function refreshToken() {
    try {
      const currentRefreshToken = secureStorage.getRefreshToken()
      if (!currentRefreshToken) throw new Error('refresh_token 不存在或已过期')

      const response = await request(
        createRequestConfig({
          url: `${BASE_URL}/users/refresh`,
          method: 'POST',
          data: { refresh_token: currentRefreshToken },
          header: { 'Content-Type': 'application/json' }
        })
      )

      if (response.statusCode === 200) {
        const { access_token, refresh_token: newRefreshToken, registered_days } = response.data
        user.token = access_token
        secureStorage.setRefreshToken(newRefreshToken)
        user.tokenExpiry = Date.now() + ACCESS_TOKEN_EXPIRES_IN * 1000
        user.isLoggedIn = true
        user.status = UserStatus.LOGGED_IN
        if (registered_days) user.registered_days = registered_days
        saveToStorage()
        scheduleTokenRefresh()
        return response.data
      } else {
        throw new Error(response.data.message || '刷新令牌失败')
      }
    } catch (error) {
      console.error('刷新令牌失败:', error)
      reset()
      uni.navigateTo({ url: '/pagesMy/login/login' })
      throw error
    }
  }

  function getUserID() {
    return user.uid
  }

  function clearStorage() {
    try {
      uni.removeStorageSync(STORAGE_KEY)
      secureStorage.removeRefreshToken()
    } catch {}
  }

  function reset() {
    user.uid = null
    user.isLoggedIn = false
    user.token = ''
    user.avatarUrl = ''
    user.nickName = ''
    user.status = UserStatus.LOGGED_OUT
    user.tokenExpiry = null
    user.registered_days = 0
    clearStorage()
    clearTokenRefreshTimer()
  }

  function loadFromLocalStorage() {
    try {
      const storedData = uni.getStorageSync(STORAGE_KEY)
      if (storedData) {
        const parsed = JSON.parse(storedData)
        user.uid = parsed.uid
        user.isLoggedIn = parsed.isLoggedIn
        user.token = parsed.token
        user.avatarUrl = parsed.avatarUrl
        user.nickName = parsed.nickName
        user.status = parsed.status
        user.tokenExpiry = parsed.tokenExpiry
        user.registered_days = parsed.registered_days || 0
      }
    } catch {}
  }

  function getStatusText() {
    const statusTexts = {
      [UserStatus.LOGGED_OUT]: '未登录',
      [UserStatus.LOGGED_IN]: '已登录'
    }
    return statusTexts[user.status] || '未知状态'
  }

  async function fetchBasicDetails() {
    try {
      const response = await request(
        createRequestConfig({
          url: `${BASE_URL}/users/basic_details`,
          method: 'GET',
          header: { 'Content-Type': 'application/json' }
        })
      )
      if (response.statusCode === 200) {
		  console.log('获取的用户信息：', response.data);
        const data = response.data
        user.uid = data.id
        user.nickName = data.nickname
        user.avatarUrl = data.avatar_url || '/static/images/default_avatar.jpg'
        user.registered_days = data.registered_days || 0
        user.isLoggedIn = true
        user.status = UserStatus.LOGGED_IN
        user.tokenExpiry = Date.now() + ACCESS_TOKEN_EXPIRES_IN * 1000
        saveToStorage()
      } else {
        throw new Error(response.data.message || '获取用户信息失败')
      }
    } catch (error) {
      console.error('获取用户信息失败:', error)
      throw error
    }
  }

  function scheduleTokenRefresh() {
    clearTokenRefreshTimer()
    if (user.tokenExpiry) {
      const currentTime = Date.now()
      const timeToExpiry = user.tokenExpiry - currentTime
      const timeBeforeRefresh = timeToExpiry - 20000
      if (timeBeforeRefresh > 0) {
        tokenRefreshTimer = setTimeout(() => {
          refreshToken().catch(err => console.error('自动刷新token失败:', err))
        }, timeBeforeRefresh)
      } else {
        refreshToken().catch(err => console.error('自动刷新token失败:', err))
      }
    }
  }

  function clearTokenRefreshTimer() {
    if (tokenRefreshTimer) {
      clearTimeout(tokenRefreshTimer)
      tokenRefreshTimer = null
    }
  }

  async function initialize() {
    if (user.isLoggedIn && secureStorage.getRefreshToken()) {
      try {
        await refreshToken()
      } catch (error) {
        console.error('应用启动时刷新token失败:', error)
      }
    }
  }

  watchUser()
  loadFromLocalStorage()
  fetchBasicDetails()
  initialize()

  return {
    user,
    UserStatus,
    getStatusText,
    login,
    setNickname,
    setAvatar,
    refreshToken,
    getUserID,
    logout,
    reset,
    clearStorage,
    loadFromLocalStorage,
    fetchBasicDetails
  }
})