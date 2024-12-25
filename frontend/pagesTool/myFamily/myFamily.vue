<template>
  <view class="family-management">
    <view class="family-list">
      <!-- 已加入的成员列表 -->
      <uni-card
        v-for="member in sortedFamilyMembers"
        :key="member.id"
        :padding="0"
        spacing="0"
        class="member-card"
      >
        <view class="member-wrapper">
          <view class="member-info">
            <image
              :src="`http://122.51.231.155:8080/static/${member.avatarUrl}`"
              mode="aspectFill"
              class="avatar"
            />
            <text class="nickname">{{ member.nickname }}</text>
            <uni-tag
              v-if="member.role === 'admin'"
              text="管理员"
              type="primary"
              size="small"
              class="role-tag"
            />
          </view>
          <view class="actions" v-if="member.id !== currentUserId">
            <text
              v-if="member.role !== 'admin'"
              class="action-btn promote-btn"
              @click="setAsAdmin(member.id)"
            >
              {{ $t('set_as_admin') }}
            </text>
            <text
              class="action-btn remove-btn"
              @click="removeMember(member.id)"
            >
              {{ $t('remove_member') }}
            </text>
          </view>
        </view>
      </uni-card>
    </view>

    <!-- 待审核成员 -->
    <uni-section
      title="待审核成员"
      type="line"
      v-if="family.waiting_members?.length > 0"
    >
      <view class="family-list">
        <uni-card
          v-for="member in family.waiting_members"
          :key="member.id"
          :padding="0"
          spacing="0"
          class="member-card"
        >
          <view class="member-wrapper">
            <view class="member-info">
              <image
                :src="`http://122.51.231.155:8080/static/${member.avatarUrl}`"
                mode="aspectFill"
                class="avatar"
              />
              <text class="nickname">{{ member.nickname }}</text>
              <uni-tag
                text="待审核"
                type="warning"
                size="small"
                class="role-tag"
              />
            </view>
            <view class="actions">
              <text
                class="action-btn approve-btn"
                @click="admitMember(member.id)"
              >
                {{ $t('admit') }}
              </text>
              <text
                class="action-btn reject-btn"
                @click="rejectMember(member.id)"
              >
                {{ $t('reject') }}
              </text>
            </view>
          </view>
        </uni-card>
      </view>
    </uni-section>
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import {computed} from 'vue'
import {onShow} from '@dcloudio/uni-app'
import {useI18n} from 'vue-i18n'
import {useFamilyStore} from '../stores/family.js'
import {useUserStore} from '@/stores/user'

/* ----------------- Setup ----------------- */
const {t} = useI18n()
const familyStore = useFamilyStore()
const family = computed(() => familyStore.family)
const userStore = useUserStore()
const currentUserId = computed(() => userStore.user.uid)

/* ----------------- Computed ----------------- */
const isCurrentUserAdmin = computed(() => familyStore.isAdmin(currentUserId.value))
const allFamilyMembers = computed(() => family.value.allMembers || [])
const sortedFamilyMembers = computed(() => {
  const members = [...allFamilyMembers.value]
  const currentUserIndex = members.findIndex(m => m.id === currentUserId.value)
  if (currentUserIndex > -1) {
    const [currentUser] = members.splice(currentUserIndex, 1)
    members.unshift(currentUser)
  }
  return members
})

/* ----------------- Lifecycle ----------------- */
onShow(async () => {
  try {
    await familyStore.getFamilyDetails()
  } catch {
    uni.showToast({title: t('failed'), icon: 'none'})
  }
})

/* ----------------- Methods ----------------- */
async function setAsAdmin(userId) {
  try {
    await familyStore.setAdmin(userId)
    uni.showToast({title: t('success'), icon: 'success'})
  } catch {
    uni.showToast({title: t('failed'), icon: 'error'})
  }
}

async function removeMember(userId) {
  try {
    await familyStore.removeFamilyMember(userId)
    uni.showToast({title: t('success'), icon: 'success'})
  } catch {
    uni.showToast({title: t('failed'), icon: 'error'})
  }
}

async function admitMember(userId) {
  try {
    await familyStore.admitJoinRequest(userId)
    uni.showToast({title: t('success'), icon: 'success'})
  } catch {
    uni.showToast({title: t('failed'), icon: 'error'})
  }
}

async function rejectMember(userId) {
  try {
    await familyStore.rejectJoinRequest(userId)
    uni.showToast({title: t('success'), icon: 'success'})
  } catch {
    uni.showToast({title: t('failed'), icon: 'error'})
  }
}
</script>

<style scoped>
.family-management {
  padding: 30rpx;
  background-color: #f8f8f8;
}

.family-list {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.member-card {
  border-radius: 12rpx;
  overflow: hidden;
}

.member-wrapper {
  padding: 24rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #ffffff;
}

.member-info {
  display: flex;
  align-items: center;
  gap: 20rpx;
  flex: 1;
}

.avatar {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  background-color: #f0f0f0;
}

.nickname {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
}

.role-tag {
  margin-left: 16rpx;
}

.actions {
  display: flex;
  gap: 16rpx;
  align-items: center;
}

.action-btn {
  padding: 12rpx 24rpx;
  border-radius: 6rpx;
  font-size: 24rpx;
  transition: opacity 0.2s;
}

.action-btn:active {
  opacity: 0.8;
}

.promote-btn {
  color: #2979ff;
  background-color: rgba(41, 121, 255, 0.1);
}

.remove-btn {
  color: #ff5252;
  background-color: rgba(255, 82, 82, 0.1);
}

.approve-btn {
  color: #00c853;
  background-color: rgba(0, 200, 83, 0.1);
}

.reject-btn {
  color: #ff5252;
  background-color: rgba(255, 82, 82, 0.1);
}
</style>