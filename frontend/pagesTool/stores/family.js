// family.js
import { defineStore } from 'pinia'
import { reactive, watch } from 'vue'
import { useUserStore } from '../../stores/user.js'

const BASE_URL = 'https://dechelper.com'
export const FamilyStatus = {
  NOT_JOINED: 'empty',
  PENDING_APPROVAL: 'waiting',
  JOINED: 'family'
}
const STORAGE_KEY = 'family_store_data'

const TIMEZONE_MAPPING = {
  480: 'Asia/Shanghai',
  540: 'Asia/Tokyo',
  420: 'Asia/Bangkok',
  330: 'Asia/Kolkata',
  240: 'Asia/Dubai',
  180: 'Europe/Moscow',
  120: 'Europe/Berlin',
  60: 'Europe/London',
  0: 'Europe/London',
  '-60': 'Atlantic/Azores',
  '-120': 'America/Noronha',
  '-180': 'America/Sao_Paulo',
  '-240': 'America/New_York',
  '-300': 'America/Chicago',
  '-360': 'America/Denver',
  '-420': 'America/Los_Angeles',
  '-480': 'America/Anchorage',
  '-540': 'America/Adak',
  '-600': 'Pacific/Honolulu'
}

function getTimeZone() {
  const offset = -new Date().getTimezoneOffset()
  return TIMEZONE_MAPPING[offset] || 'Asia/Shanghai'
}

function createRequestConfig(config) {
  const userStore = useUserStore()
  return {
    ...config,
    header: {
      Authorization: `Bearer ${userStore.user.token || ''}`,
      ...(config.header || {})
    }
  }
}

function request(config) {
  const userStore = useUserStore()
  return new Promise((resolve, reject) => {
    uni.request({
      ...config,
      success: async res => {
        if (res.statusCode === 401) {
          try {
            await userStore.refreshToken()
            uni.request({
              ...config,
              header: {
                ...config.header,
                Authorization: `Bearer ${userStore.user.token}`
              },
              success: res2 => {
                if (res2.statusCode === 401) reject(new Error('Unauthorized'))
                else resolve(res2)
              },
              fail: err2 => reject(err2)
            })
          } catch {
            uni.navigateTo({ url: '/pagesMy/login/login' })
            reject(new Error('Unauthorized'))
          }
        } else {
          resolve(res)
        }
      },
      fail: err => reject(err)
    })
  })
}

function getInitialState() {
  try {
    const stored = uni.getStorageSync(STORAGE_KEY)
    return stored
      ? JSON.parse(stored)
      : {
          id: '',
          name: '',
          familyId: '',
          memberCount: 0,
          allMembers: [],
          waiting_members: [],
          dishProposals: [],
          status: FamilyStatus.NOT_JOINED,
          memberDailyData: [],
          familySums: {}
        }
  } catch {
    return {
      id: '',
      name: '',
      familyId: '',
      memberCount: 0,
      allMembers: [],
      waiting_members: [],
      dishProposals: [],
      status: FamilyStatus.NOT_JOINED,
      memberDailyData: [],
      familySums: {}
    }
  }
}

export const useFamilyStore = defineStore('family', () => {
  const family = reactive(getInitialState())

  function saveToStorage() {
    try {
      uni.setStorageSync(STORAGE_KEY, JSON.stringify(family))
    } catch {}
  }

  function watchFamily() {
    const watchKeys = [
      'id',
      'name',
      'familyId',
      'memberCount',
      'allMembers',
      'waiting_members',
      'status',
      'dishProposals',
      'memberDailyData',
      'familySums'
    ]
    watchKeys.forEach(key => {
      watch(
        () => family[key],
        () => {
          saveToStorage()
        },
        { deep: true }
      )
    })
  }

  async function createFamily(familyName) {
    try {
      const res = await request(
        createRequestConfig({
          url: `${BASE_URL}/families/create`,
          method: 'POST',
          data: { name: familyName }
        })
      )
      family.id = res.data.family.id
      family.name = res.data.family.name
      family.familyId = res.data.family.family_id
      family.status = FamilyStatus.JOINED
      await getFamilyDetails()
      return res.data
    } catch (error) {
      console.error('创建家庭失败:', error)
      throw error
    }
  }

  async function getFamilyDetails() {
    try {
      const res = await request(
        createRequestConfig({
          url: `${BASE_URL}/families/details?timezone=${encodeURIComponent(getTimeZone())}`,
          method: 'GET'
        })
      )
      const data = res.data
      family.status = data.status
      if (data.status === FamilyStatus.JOINED) {
        family.id = data.id
        family.name = data.name
        family.familyId = data.family_id
        family.memberCount = data.member_count
        const adminsWithRole = data.admins.map(a => ({
          id: a.id,
          nickname: a.nickname,
          avatarUrl: a.avatar_url,
          role: 'admin'
        }))
        const membersWithRole = data.members.map(m => ({
          id: m.id,
          nickname: m.nickname,
          avatarUrl: m.avatar_url,
          role: 'member'
        }))
        family.allMembers = [...adminsWithRole, ...membersWithRole]
        if (Array.isArray(data.waiting_members)) {
          family.waiting_members = data.waiting_members.map(m => ({
            id: m.id,
            nickname: m.nickname,
            avatarUrl: m.avatar_url
          }))
        } else {
          family.waiting_members = []
        }
        family.memberDailyData = data.member_daily_data || []
        family.familySums = data.family_sums || {}
      } else if (data.status === FamilyStatus.PENDING_APPROVAL) {
        family.id = data.id
        family.name = data.name
        family.familyId = data.family_id
        family.memberCount = 0
        family.allMembers = []
        family.waiting_members = []
      } else {
        reset()
      }
      return res
    } catch (error) {
      console.error('getFamilyDetails error:', error)
      throw error
    }
  }

  async function getDesiredDishes() {
    try {
      const res = await request(createRequestConfig({ url: `${BASE_URL}/families/desired_dishes`, method: 'GET' }))
      family.dishProposals = res.data
      return res.data
    } catch (error) {
      console.error('getDesiredDishes error:', error)
      throw error
    }
  }

  async function addDishProposal({ dishId, preference }) {
    try {
      const res = await request(
        createRequestConfig({
          url: `${BASE_URL}/families/add_desired_dish`,
          method: 'POST',
          data: { dish_id: dishId, level_of_desire: preference }
        })
      )
      if (res.statusCode === 200) {
        await getDesiredDishes()
        return res.data
      } else if (
        res.statusCode === 400 &&
        res.data.error === 'You have already desired this dish'
      ) {
        throw new Error('DISH_ALREADY_EXISTS')
      } else {
        throw new Error(res.data.error || 'Unknown error')
      }
    } catch (error) {
      console.error('addDishProposal error:', error)
      throw error
    }
  }

  async function deleteDesiredDish(dishId) {
    try {
      const res = await request(
        createRequestConfig({
          url: `${BASE_URL}/families/desired_dishes`,
          method: 'DELETE',
          data: { dish_id: dishId }
        })
      )
      await getDesiredDishes()
      return res.data
    } catch (error) {
      console.error('deleteDesiredDish error:', error)
      throw error
    }
  }

  async function searchFamily(familyId) {
    try {
      const res = await request(
        createRequestConfig({
          url: `${BASE_URL}/families/search?family_id=${familyId}`,
          method: 'GET'
        })
      )
      return res.data
    } catch (error) {
      console.error('searchFamily error:', error)
      throw error
    }
  }

  async function joinFamily(familyId) {
    try {
      const res = await request(
        createRequestConfig({
          url: `${BASE_URL}/families/${familyId}/join`,
          method: 'POST'
        })
      )
      await getFamilyDetails()
      return res.data
    } catch (error) {
      console.error('joinFamily error:', error)
      throw error
    }
  }

  async function cancelJoinRequest() {
    try {
      const res = await request(
        createRequestConfig({
          url: `${BASE_URL}/families/cancel_join`,
          method: 'DELETE'
        })
      )
      await getFamilyDetails()
      return res.data
    } catch (error) {
      console.error('cancelJoinRequest error:', error)
      throw error
    }
  }

  async function admitJoinRequest(userId) {
    try {
      const res = await request(
        createRequestConfig({
          url: `${BASE_URL}/families/admit`,
          method: 'POST',
          data: { user_id: userId }
        })
      )
      await getFamilyDetails()
      return res.data
    } catch (error) {
      console.error('admitJoinRequest error:', error)
      throw error
    }
  }

  async function rejectJoinRequest(userId) {
    try {
      const res = await request(
        createRequestConfig({
          url: `${BASE_URL}/families/reject`,
          method: 'POST',
          data: { user_id: userId }
        })
      )
      await getFamilyDetails()
      return res.data
    } catch (error) {
      console.error('rejectJoinRequest error:', error)
      throw error
    }
  }

  async function leaveFamily() {
    try {
      const userStore = useUserStore()
      const res = await request(
        createRequestConfig({
          url: `${BASE_URL}/families/leave_family`,
          method: 'DELETE',
          data: { user_id: userStore.user.uid }
        })
      )
      reset()
      return res.data
    } catch (error) {
      console.error('leaveFamily error:', error)
      throw error
    }
  }

  async function breakFamily() {
    try {
      const res = await request(createRequestConfig({ url: `${BASE_URL}/families/break`, method: 'DELETE' }))
      reset()
      return res.data
    } catch (error) {
      console.error('breakFamily error:', error)
      throw error
    }
  }

  async function setAdmin(userId) {
    try {
      const res = await request(
        createRequestConfig({
          url: `${BASE_URL}/families/set_admin`,
          method: 'PUT',
          data: { user_id: userId }
        })
      )
      await getFamilyDetails()
      return res.data
    } catch (error) {
      console.error('setAdmin error:', error)
      throw error
    }
  }

  async function removeFamilyMember(userId) {
    try {
      const res = await request(
        createRequestConfig({
          url: `${BASE_URL}/families/delete_family_member`,
          method: 'DELETE',
          data: { user_id: userId }
        })
      )
      await getFamilyDetails()
      return res.data
    } catch (error) {
      console.error('removeFamilyMember error:', error)
      throw error
    }
  }

  function isAdmin(userId) {
    return family.allMembers.some(member => member.id === userId && member.role === 'admin')
  }

  function clearStorage() {
    try {
      uni.removeStorageSync(STORAGE_KEY)
    } catch {}
  }

  function reset() {
    family.id = ''
    family.name = ''
    family.familyId = ''
    family.memberCount = 0
    family.allMembers = []
    family.waiting_members = []
    family.dishProposals = []
    family.status = FamilyStatus.NOT_JOINED
    family.memberDailyData = []
    family.familySums = {}
    clearStorage()
  }

  function getStatusText() {
    const statusTexts = {
      [FamilyStatus.NOT_JOINED]: '未加入',
      [FamilyStatus.PENDING_APPROVAL]: '待审核',
      [FamilyStatus.JOINED]: '已加入'
    }
    return statusTexts[family.status] || '未知状态'
  }

  watchFamily()

  return {
    family,
    FamilyStatus,
    getStatusText,
    createFamily,
    getFamilyDetails,
    searchFamily,
    joinFamily,
    cancelJoinRequest,
    admitJoinRequest,
    rejectJoinRequest,
    setAdmin,
    removeFamilyMember,
    leaveFamily,
    breakFamily,
    reset,
    clearStorage,
    isAdmin,
    addDishProposal,
    getDesiredDishes,
    deleteDesiredDish
  }
})