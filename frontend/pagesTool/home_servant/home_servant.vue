<template>
  <view class="container">
    <!-- 全屏背景图片 -->
    <image
      src="/static/images/index/background_img.jpg"
      class="background-image"
    />

    <!-- 已加入家庭：顶部标题 -->
    <view
      v-if="family.status === FamilyStatus.JOINED"
      class="header"
    >
      <text class="header-title">
        {{ family.name || t('default_family_name') }}
      </text>
    </view>

    <!-- 家庭ID -->
    <text
      v-if="family.status === FamilyStatus.JOINED"
      class="list-title"
    >
      {{ t('family_info') + family.familyId }}
    </text>

    <!-- 用户未加入家庭 -->
    <view
      v-if="family.status === FamilyStatus.NOT_JOINED"
      class="no-family-view"
    >
      <image
        src="https://cdn.pixabay.com/photo/2017/01/13/02/31/family-1976162_1280.png"
        class="centered-image"
      />
      <view class="family-actions">
        <button
          class="action-button"
          @click="showCreateFamilyModal = true"
        >
          {{ t('create_family') }}
        </button>
        <button
          class="action-button"
          @click="showJoinFamilyModal = true"
        >
          {{ t('join_family') }}
        </button>
      </view>
    </view>

    <!-- 待审核状态 -->
    <view
      v-if="family.status === FamilyStatus.PENDING_APPROVAL"
      class="pending-view"
    >
      <image
        src="https://cdn.pixabay.com/photo/2021/09/20/22/15/hourglass-6641967_1280.png"
        class="pending-image"
      />
      <text class="pending-text">
        {{ t('pending_approval') }}
      </text>
      <button class="cancel-button" @click="handleCancelJoin">
        {{ t('cancel_join') }}
      </button>
    </view>

    <!-- 管理员提示：有等待审核的成员 -->
    <uni-card
      v-if="isCurrentUserAdmin && family.waiting_members?.length > 0"
      class="admin-notice"
      :is-shadow="false"
      :border="false"
    >
      <view class="notice-content">
        <view class="notice-left">
          <uni-icons
            type="info-filled"
            size="20"
            color="#3A86FF"
          />
          <text class="notice-text">
            {{ t('new_join_requests') }}
          </text>
        </view>
        <text class="notice-btn" @click="goToMyFamily">
          {{ t('handle_requests') }}
        </text>
      </view>
    </uni-card>

    <!-- 已加入家庭后的管理部分 -->
    <view
      v-if="family.status === FamilyStatus.JOINED"
      class="family-management"
    >
      <!-- 提议菜品模块 -->
      <view class="dish-proposal">
        <text class="section-title">{{ t('propose_dish') }}</text>
        <uni-combox
          :placeholder="t('dish_name_placeholder')"
          v-model="foodNameInput"
          :candidates="filteredFoods.map(item => displayName(item))"
          @input="onComboxInput"
        />
        <picker
          mode="selector"
          :range="dishPreferenceLevels"
          :value="newDish.preference"
          @change="onDishPreferenceChange"
          class="picker"
        >
          <view>
            {{ t('dish_preference') }}:
            {{ dishPreferenceLevels[newDish.preference] }}
          </view>
        </picker>
        <button
          class="submit-button"
          @click="submitDishProposal"
        >
          {{ t('submit_proposal') }}
        </button>
      </view>

      <!-- 家庭成员的提议列表 -->
      <view class="dish-list">
        <text class="section-title">
          {{ t('family_dish_proposals') }}
        </text>
        <scroll-view
          scroll-y
          class="proposals-scroll"
        >
          <uni-collapse
            :accordion="true"
            v-model="activeCollapse"
          >
            <uni-collapse-item
              v-for="dish in sortedDishProposals"
              :key="dish.id"
              :name="dish.id.toString()"
            >
              <!-- 自定义折叠面板标题 -->
              <template #title>
                <view class="collapse-title">
                  <text class="dish-name">
                    {{ dish.name }}
                  </text>
                  <view class="dish-info">
                    <text class="proposer-text">
                      {{ dish.proposer }}
                    </text>
                    <uni-tag
                      :text="dishPreferenceLevels[dish.preference]"
                      :type="getPreferenceTagType(dish.preference)"
                      size="small"
                    />
                  </view>
                </view>
              </template>
              <!-- 面板内容：删除操作 -->
              <view class="dish-actions">
                <uni-icons
                  v-if="canDeleteDish(dish)"
                  type="trash"
                  size="24"
                  color="#FF4D4F"
                  @click.stop="confirmDeleteDish(dish.id)"
                  class="delete-icon"
                />
              </view>
            </uni-collapse-item>
          </uni-collapse>
        </scroll-view>
      </view>

      <!-- 共享家庭碳排放 - 环形图 -->
      <view class="shared-data">
        <text class="section-title">
          {{ t('shared_family_data') }}
        </text>
        <view class="charts">
          <qiun-data-charts
            :canvas2d="true"
            canvas-id="familyCarbonChart"
            type="ring"
            :opts="carbonRingOpts"
            :chartData="carbonChartData"
          />
        </view>
      </view>

      <!-- 家庭成员的营养超标情况 -->
      <view class="nutrition-over-section">
        <text class="section-title">
          {{ t('nutrition_over_title') }}
        </text>
        <view
          v-for="member in membersNutritionOver"
          :key="member.userId"
          class="nutrition-over-row"
        >
          <image
            :src="`${BASE_URL}/static/${member.avatarUrl}`"
            class="member-avatar"
          />
          <text class="member-name">
            {{ member.nickname }}
          </text>
          <!-- 若有超标 -->
          <text
            class="over-list"
            v-if="member.overs.length > 0"
          >
            {{ member.overs.join('、') + t('nutrition_over_suffix') }}
          </text>
          <!-- 若无超标 -->
          <text
            class="over-list"
            v-else
          >
            {{ t('no_nutrition_over') }}
          </text>
        </view>
      </view>

      <!-- 家庭成员展示 -->
      <view class="family-info">
        <text class="section-title">
          {{ t('family_members') }}
        </text>
        <view class="family-members">
          <view
            v-for="member in sortedFamilyMembers"
            :key="member.id"
            class="member"
          >
            <image
              :src="`${BASE_URL}/static/${member.avatarUrl}`"
              class="member-avatar"
            />
            <text class="member-name">
              {{
                `${member.nickname} (${
                  member.role === 'admin'
                    ? t('admin')
                    : t('member')
                })`
              }}
            </text>
          </view>
        </view>
        <!-- 仅管理员可见的管理按钮 -->
        <button
          v-if="isCurrentUserAdmin"
          class="manage-members-button"
          @click="manageMembers"
        >
          {{ t('manage_members') }}
        </button>
      </view>

      <!-- 底部操作按钮 -->
      <view class="bottom-container">
        <uni-section
          class="action-section"
          type="line"
        >
          <view class="bottom-actions">
            <text
              class="action-btn leave-btn"
              @click="handleLeaveFamily"
            >
              <uni-icons
                type="close"
                size="16"
                color="#fff"
              />
              {{ t('leave_family') }}
            </text>
            <text
              v-if="isCurrentUserAdmin"
              class="action-btn break-btn"
              @click="handleBreakFamily"
            >
              <uni-icons
                type="trash"
                size="16"
                color="#fff"
              />
              {{ t('break_family') }}
            </text>
          </view>
        </uni-section>
      </view>
    </view>

    <!-- 创建家庭模态框 -->
    <view
      v-if="showCreateFamilyModal"
      class="modal"
    >
      <view class="modal-content">
        <text class="modal-title">
          {{ t('create_family') }}
        </text>
        <input
          v-model="newFamilyName"
          :placeholder="t('family_name_placeholder')"
          class="input"
        />
        <button
          class="modal-button"
          @click="createFamily"
        >
          {{ t('confirm') }}
        </button>
        <button
          class="modal-button cancel"
          @click="showCreateFamilyModal = false"
        >
          {{ t('cancel') }}
        </button>
      </view>
    </view>

    <!-- 加入家庭模态框 -->
    <view
      v-if="showJoinFamilyModal"
      class="modal"
    >
      <view class="modal-content">
        <text class="modal-title">
          {{ t('join_family') }}
        </text>
        <input
          v-model="joinFamilyId"
          :placeholder="t('family_id_placeholder')"
          class="input"
        />
        <button
          class="modal-button"
          @click="joinFamily"
        >
          {{ t('confirm') }}
        </button>
        <button
          class="modal-button cancel"
          @click="showJoinFamilyModal = false"
        >
          {{ t('cancel') }}
        </button>
      </view>
    </view>

    <!-- 删除确认模态框 -->
    <view
      v-if="showDeleteConfirm"
      class="modal"
    >
      <view class="modal-content">
        <text class="modal-title">
          {{ t('confirm_delete') }}
        </text>
        <button
          class="modal-button"
          @click="deleteDish(selectedDishId)"
        >
          {{ t('confirm') }}
        </button>
        <button
          class="modal-button cancel"
          @click="showDeleteConfirm = false"
        >
          {{ t('cancel') }}
        </button>
      </view>
    </view>
  </view>
</template>

<script setup>
/* ------------------ Imports ------------------ */
import { ref, computed, onMounted, watch } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useI18n } from 'vue-i18n'
import { useFamilyStore, FamilyStatus } from '../stores/family.js'
import { useUserStore } from '@/stores/user'
import { useFoodListStore } from '../stores/food_list'

/* ------------------ Setup & Refs ------------------ */
const { t, locale } = useI18n()
const familyStore = useFamilyStore()
const userStore = useUserStore()
const foodStore = useFoodListStore()

// 解构 family
const family = computed(() => familyStore.family)
const FamilyStatusEnum = FamilyStatus

const BASE_URL = 'http://xcxcs.uwdjl.cn:8080'

// 当前用户
if (!userStore.user.isLoggedIn) {
  uni.navigateTo({ url: '/pagesMy/login/login' })
}
const currentUserId = computed(() => userStore.user.uid)

// 判断是否管理员
const isCurrentUserAdmin = computed(() => {
  return familyStore.isAdmin(currentUserId.value)
})

// 家庭成员排序：当前用户排在最前
const allFamilyMembers = computed(() => family.value.allMembers || [])
const sortedFamilyMembers = computed(() => {
  const members = [...allFamilyMembers.value]
  const currentUserIndex = members.findIndex(
    m => m.id === currentUserId.value
  )
  if (currentUserIndex > -1) {
    const [currentUser] = members.splice(currentUserIndex, 1)
    members.unshift(currentUser)
  }
  return members
})

// 提议菜品
const newDish = ref({
  name: '',
  id: null,
  preference: 0
})
const dishPreferenceLevels = computed(() => [
  t('preference_low'),
  t('preference_medium'),
  t('preference_high')
])
function getPreferenceTagType(preference) {
  switch (preference) {
    case 0:
      return 'info'
    case 1:
      return 'warning'
    case 2:
      return 'success'
    default:
      return 'default'
  }
}

// 折叠面板
const activeCollapse = ref([])

// 模态框
const showCreateFamilyModal = ref(false)
const showJoinFamilyModal = ref(false)
const showDeleteConfirm = ref(false)

// 输入
const newFamilyName = ref('')
const joinFamilyId = ref('')

// 引入食物列表store
const { availableFoods, fetchAvailableFoods, getFoodName } = foodStore
const foodNameInput = ref('')
const showFoodList = ref(false)

/* ------------------ Computed ------------------ */
// 下拉候选
const filteredFoods = computed(() => {
  if (!foodNameInput.value) {
    // 当输入为空时
    const lang = locale.value
    return lang === 'zh-Hans'
      ? availableFoods.filter(f => f.name_zh)
      : availableFoods.filter(f => f.name_en)
  } else {
    const lang = locale.value
    return availableFoods.filter(f => {
      if (lang === 'zh-Hans') {
        return f.name_zh.includes(foodNameInput.value)
      } else {
        return f.name_en
          .toLowerCase()
          .includes(foodNameInput.value.toLowerCase())
      }
    })
  }
})

// 提议的菜品列表
const sortedDishProposals = computed(() => {
  if (!family.value.dishProposals) return []
  return family.value.dishProposals
    .map(item => ({
      id: item.dish_id,
      name: getFoodName(item.dish_id) || 'Unknown Dish',
      preference: item.level_of_desire,
      proposer: item.proposer_user
        ? item.proposer_user.nickname
        : 'Unknown'
    }))
    .sort((a, b) => b.preference - a.preference)
})

// 环形图
const carbonRingOpts = ref({
  title: {
    name: '',
    fontSize: 16,
    color: '#333'
  },
  subtitle: {
    name: '',
    fontSize: 14,
    color: '#999'
  },
  extra: {
    ring: {
      ringWidth: 10,
      activeOpacity: 0.5,
      activeRadius: 10,
      offsetAngle: 0,
      labelWidth: 15,
      border: false,
      borderWidth: 3,
      borderColor: '#FFFFFF'
    }
  }
})
const carbonChartData = ref({ series: [] })

// 每日数据变化时更新环形图
watch(
  () => family.value.memberDailyData,
  newVal => {
    if (!newVal || newVal.length === 0) {
      carbonChartData.value.series = []
      carbonRingOpts.value.title.name = t('actual_target')
      carbonRingOpts.value.subtitle.name = `0 / 0Kg`
      carbonRingOpts.value.subtitle.color = '#999'
      return
    }
    // 构造 series
    carbonChartData.value.series = newVal.map(member => ({
      name: member.nickname,
      data: member.carbon_intake_sum
    }))
    // 中心文字
    const { familySums } = family.value
    if (
      familySums &&
      familySums.carbon_intake_sum !== undefined &&
      familySums.carbon_goal_sum !== undefined
    ) {
      const actual = familySums.carbon_intake_sum || 0
      const target = familySums.carbon_goal_sum || 0
      carbonRingOpts.value.title.name = t('actual_target')
      carbonRingOpts.value.subtitle.name = `${actual} / ${target}Kg`
      carbonRingOpts.value.subtitle.color =
        actual > target ? '#FF4D4F' : '#999'
    } else {
      carbonRingOpts.value.title.name = t('actual_target')
      carbonRingOpts.value.subtitle.name = `0 / 0Kg`
      carbonRingOpts.value.subtitle.color = '#999'
    }
  },
  { immediate: true, deep: true }
)

// 计算营养超标
const nutritionKeys = [
  { key: 'calories', i18nKey: 'calories' },
  { key: 'protein', i18nKey: 'protein' },
  { key: 'fat', i18nKey: 'fat' },
  { key: 'carbohydrates', i18nKey: 'carbohydrates' },
  { key: 'sodium', i18nKey: 'sodium' }
]
const membersNutritionOver = computed(() => {
  if (!family.value.memberDailyData?.length) return []
  return family.value.memberDailyData.map(item => {
    const overs = []
    nutritionKeys.forEach(nut => {
      const goalVal = item.nutrition_goal[nut.key]
      const intakeVal = item.nutrition_intake_sum[nut.key]
      if (intakeVal > goalVal) {
        overs.push(t(nut.i18nKey))
      }
    })
    return {
      userId: item.user_id,
      nickname: item.nickname,
      avatarUrl: item.avatar_url,
      overs
    }
  })
})

/* ------------------ Lifecycle ------------------ */
onShow(async () => {
  try {
    await familyStore.getFamilyDetails()
    if (family.value.status === FamilyStatusEnum.JOINED) {
      await familyStore.getDesiredDishes()
    }
  } catch {
    uni.showToast({
      title: t('fetch_family_failed'),
      icon: 'none'
    })
  }
})
onMounted(async () => {
  if (availableFoods.length === 0) {
    await fetchAvailableFoods()
  }
})

/* ------------------ Methods ------------------ */
// 下拉输入
function onComboxInput(value) {
  foodNameInput.value = value
}

// 选择食物
function selectFood(foodItem) {
  newDish.value.name = foodItem.name_en
  newDish.value.id = foodItem.id
  foodNameInput.value = displayName(foodItem)
  showFoodList.value = false
}

// 工具函数：根据语言显示
function displayName(item) {
  return locale.value === 'zh-Hans'
    ? item.name_zh
    : item.name_en
}

// 提交菜品提议
async function submitDishProposal() {
  // 匹配
  const matchedFood = availableFoods.find(
    f => displayName(f) === foodNameInput.value
  )
  if (matchedFood) {
    selectFood(matchedFood)
  } else {
    uni.showToast({
      title: t('no_matching_food'),
      icon: 'none'
    })
    return
  }

  if (!newDish.value.name.trim()) {
    uni.showToast({
      title: t('dish_name_required'),
      icon: 'none'
    })
    return
  }

  try {
    await familyStore.addDishProposal({
      dishId: Number(newDish.value.id),
      preference: newDish.value.preference
    })
    // 重置
    newDish.value = { name: '', id: null, preference: 0 }
    foodNameInput.value = ''
    uni.showToast({ title: t('success'), icon: 'success' })
  } catch (error) {
    if (error.message === 'DISH_ALREADY_EXISTS') {
      uni.showToast({
        title: t('dish_already_exists'),
        icon: 'none'
      })
    } else {
      uni.showToast({
        title: t('failed'),
        icon: 'error'
      })
    }
  }
}

// 偏好更改
function onDishPreferenceChange(e) {
  newDish.value.preference = parseInt(e.detail.value, 10)
}

// 创建家庭
async function createFamily() {
  if (!newFamilyName.value.trim()) {
    uni.showToast({
      title: t('family_name_required'),
      icon: 'none'
    })
    return
  }
  try {
    await familyStore.createFamily(newFamilyName.value)
    newFamilyName.value = ''
    showCreateFamilyModal.value = false
    uni.showToast({ title: t('create_family_success'), icon: 'success' })
    await familyStore.getDesiredDishes()
  } catch {
    uni.showToast({
      title: t('create_family_failed'),
      icon: 'error'
    })
  }
}

// 加入家庭
async function joinFamily() {
  if (!joinFamilyId.value.trim()) {
    uni.showToast({
      title: t('family_id_required'),
      icon: 'none'
    })
    return
  }
  try {
    const searchResult = await familyStore.searchFamily(joinFamilyId.value)
    if (!searchResult || !searchResult.id) {
      uni.showToast({
        title: t('family_not_found'),
        icon: 'error'
      })
      return
    }
    await familyStore.joinFamily(searchResult.id)
    joinFamilyId.value = ''
    showJoinFamilyModal.value = false
    uni.showToast({
      title: t('join_request_sent'),
      icon: 'success'
    })
  } catch {
    uni.showToast({
      title: t('join_family_failed'),
      icon: 'error'
    })
  }
}

// 取消加入
async function handleCancelJoin() {
  try {
    await familyStore.cancelJoinRequest()
    uni.showToast({ title: t('success'), icon: 'success' })
  } catch {
    uni.showToast({ title: t('failed'), icon: 'error' })
  }
}

// 前往 myFamily 页面
function goToMyFamily() {
  uni.navigateTo({ url: '/pagesTool/myFamily/myFamily' })
}

// 退出家庭
async function handleLeaveFamily() {
  try {
    await familyStore.leaveFamily()
    uni.showToast({ title: t('success'), icon: 'success' })
  } catch {
    uni.showToast({ title: t('failed'), icon: 'error' })
  }
}

// 解散家庭
async function handleBreakFamily() {
  try {
    await familyStore.breakFamily()
    uni.showToast({ title: t('break_family_success'), icon: 'success' })
  } catch {
    uni.showToast({ title: t('break_family_failed'), icon: 'error' })
  }
}

// 管理成员
function manageMembers() {
  uni.navigateTo({ url: '/pagesTool/myFamily/myFamily' })
}

// 删除菜品
function confirmDeleteDish(dishId) {
  selectedDishId.value = dishId
  showDeleteConfirm.value = true
}
async function deleteDish(dishId) {
  try {
    await familyStore.deleteDesiredDish(dishId)
    uni.showToast({ title: t('delete_success'), icon: 'success' })
  } catch {
    uni.showToast({ title: t('delete_failed'), icon: 'error' })
  } finally {
    showDeleteConfirm.value = false
    selectedDishId.value = null
  }
}

// 判断是否可删除（提议者或管理员）
function canDeleteDish(dish) {
  return dish.proposer === userStore.user.nickName
}
</script>

<style scoped>
:root {
  --primary-color: #4caf50;
  --secondary-color: #2fc25b;
  --background-color: #f5f5f5;
  --card-background: rgba(255, 255, 255, 0.8);
  --text-color: #333;
  --shadow-color: rgba(0, 0, 0, 0.1);
  --font-size-title: 32rpx;
  --font-size-subtitle: 24rpx;
}

.container {
  display: flex;
  flex-direction: column;
  background-color: var(--background-color);
  min-height: 100vh;
  padding-bottom: 80rpx;
  position: relative;
  overflow: hidden;
}

/* 背景图 */
.background-image {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  z-index: 0;
  opacity: 0.08;
  pointer-events: none;
}

/* 已加入家庭：头部标题 */
.header {
  padding: 40rpx 20rpx 20rpx;
  text-align: center;
}
.header-title {
  font-size: 48rpx;
  color: var(--primary-color);
  font-weight: bold;
  animation: slideDown 1s ease-out;
}

.list-title {
  margin-left: 10rpx;
  font-size: 20rpx;
  font-weight: bold;
  color: var(--text-color);
  margin-bottom: 20rpx;
  text-align: center;
}

/* 用户未加入家庭 */
.no-family-view {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20rpx;
}
.centered-image {
  width: 200rpx;
  height: 200rpx;
  border-radius: 50%;
  margin-bottom: 30rpx;
  object-fit: cover;
}

/* 待审核状态 */
.pending-view {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20rpx;
}
.pending-image {
  width: 200rpx;
  height: 200rpx;
  border-radius: 50%;
  margin-bottom: 30rpx;
  object-fit: cover;
}
.pending-text {
  font-size: 32rpx;
  color: #666;
  margin-bottom: 40rpx;
}
.cancel-button {
  background-color: #ff4d4f;
  color: #fff;
  padding: 15rpx 40rpx;
  border-radius: 10rpx;
  font-size: 28rpx;
}

/* 管理员提示框 */
.admin-notice {
  margin: 20rpx;
  background: linear-gradient(to right, #eef2ff, #e0e7ff);
  border-radius: 16rpx;
}
.admin-notice :deep(.uni-card__content) {
  padding: 0 !important;
}
.notice-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24rpx;
}
.notice-left {
  display: flex;
  align-items: center;
  gap: 12rpx;
}
.notice-text {
  color: #3a86ff;
  font-size: 28rpx;
  font-weight: 500;
}
.notice-btn {
  color: #3a86ff;
  font-size: 28rpx;
  font-weight: 600;
  background-color: rgba(58, 134, 255, 0.1);
  padding: 8rpx 24rpx;
  border-radius: 30rpx;
  transition: all 0.3s ease;
}
.notice-btn:active {
  opacity: 0.8;
  transform: scale(0.98);
}

/* 家庭管理 */
.family-management {
  flex: 1;
  padding: 20rpx;
}

/* 创建/加入家庭 */
.family-actions {
  display: flex;
  justify-content: center;
  margin-top: 20rpx;
}
.action-button {
  background-color: var(--primary-color);
  color: #fff;
  padding: 15rpx 30rpx;
  margin: 0 10rpx;
  border-radius: 10rpx;
  font-size: var(--font-size-subtitle);
}

/* 提议菜品 */
.dish-proposal {
  background-color: var(--card-background);
  padding: 20rpx;
  border-radius: 10rpx;
  box-shadow: 0 2rpx 5rpx var(--shadow-color);
  margin-bottom: 20rpx;
}
.proposals-scroll {
  max-height: 600rpx;
  margin: 20rpx 0;
}
.section-title {
  font-size: var(--font-size-title);
  font-weight: bold;
  color: var(--primary-color);
  margin-bottom: 10rpx;
}
.picker {
  padding: 10rpx;
  border: 1rpx solid #ccc;
  border-radius: 5rpx;
  margin: 20rpx 0;
  color: #666;
  z-index: 1;
}
.submit-button {
  background-color: var(--primary-color);
  color: #fff;
  padding: 15rpx;
  border-radius: 10rpx;
  font-size: var(--font-size-subtitle);
  width: 100%;
  text-align: center;
}

/* 提议列表 */
.dish-list {
  background-color: var(--card-background);
  padding: 20rpx;
  border-radius: 10rpx;
  box-shadow: 0 2rpx 5rpx var(--shadow-color);
  margin-bottom: 20rpx;
}
.collapse-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 8rpx 0;
}
.dish-name {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
  margin-right: 16rpx;
}
.dish-info {
  display: flex;
  align-items: center;
  gap: 16rpx;
  flex-shrink: 0;
}
.proposer-text {
  font-size: 24rpx;
  color: #666;
}
.dish-actions {
  display: flex;
  justify-content: flex-end;
  padding: 16rpx 0;
}

/* 共享数据(碳排放) */
.shared-data {
  background-color: var(--card-background);
  padding: 20rpx;
  border-radius: 10rpx;
  box-shadow: 0 2rpx 5rpx var(--shadow-color);
  margin-bottom: 20rpx;
}
.charts {
  width: 100%;
  height: 400rpx;
  margin-bottom: 20rpx;
}

/* 营养超标 */
.nutrition-over-section {
  background-color: var(--card-background);
  padding: 20rpx;
  border-radius: 10rpx;
  box-shadow: 0 2rpx 5rpx var(--shadow-color);
  margin-bottom: 20rpx;
}
.nutrition-over-row {
  display: flex;
  align-items: center;
  margin-bottom: 20rpx;
}
.member-avatar {
  width: 60rpx;
  height: 60rpx;
  border-radius: 50%;
  margin-right: 20rpx;
}
.member-name {
  font-size: 26rpx;
  font-weight: 500;
  margin-right: 20rpx;
}
.over-list {
  font-size: 24rpx;
  color: #ff4d4f;
}

/* 家庭成员列表 */
.family-info {
  background-color: var(--card-background);
  padding: 20rpx;
  border-radius: 10rpx;
  box-shadow: 0 2rpx 5rpx var(--shadow-color);
  margin-bottom: 20rpx;
}
.family-members {
  display: flex;
  flex-wrap: wrap;
}
.member {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 25%;
  margin-bottom: 20rpx;
}
.member-avatar {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  margin-bottom: 5rpx;
}
.member-name {
  font-size: 20rpx;
  color: var(--text-color);
}
.manage-members-button {
  background-color: var(--secondary-color);
  color: #fff;
  padding: 15rpx;
  border-radius: 10rpx;
  font-size: var(--font-size-subtitle);
  width: 100%;
  text-align: center;
  margin-top: 10rpx;
}

/* 底部操作按钮 */
.bottom-container {
  padding: 20rpx 30rpx 40rpx;
  margin-top: auto;
}
.action-section :deep(.uni-section-header) {
  display: none;
}
.bottom-actions {
  display: flex;
  justify-content: center;
  gap: 30rpx;
  padding: 20rpx 0;
}
.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 8rpx;
  padding: 20rpx 40rpx;
  border-radius: 12rpx;
  font-size: 28rpx;
  font-weight: 500;
  transition: all 0.3s ease;
}
.action-btn:active {
  opacity: 0.8;
  transform: scale(0.98);
}
.leave-btn {
  background: linear-gradient(to right, #ff9800, #f57c00);
  color: #fff;
  box-shadow: 0 4rpx 12rpx rgba(255, 152, 0, 0.3);
}
.break-btn {
  background: linear-gradient(to right, #ff4b4b, #e53935);
  color: #fff;
  box-shadow: 0 4rpx 12rpx rgba(229, 57, 53, 0.3);
}

/* 模态框 */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}
.modal-content {
  background-color: #fff;
  padding: 40rpx;
  border-radius: 16rpx;
  width: 60%;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.1);
}
.modal-title {
  font-size: 36rpx;
  color: #333;
  margin-bottom: 40rpx;
  text-align: center;
  font-weight: 500;
}
.modal-button {
  background-color: var(--primary-color);
  color: #fff;
  padding: 20rpx;
  border-radius: 12rpx;
  font-size: 24rpx;
  width: 100%;
  text-align: center;
  margin-top: 20rpx;
  font-weight: 500;
}
.modal-button.cancel {
  background-color: #f5f5f5;
  color: #666;
}

/* 动画 */
@keyframes slideDown {
  from {
    transform: translateY(-20rpx);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}
</style>