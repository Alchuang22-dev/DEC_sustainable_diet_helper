<!-- myFamily.vue -->
<template>
  <view class="family-management">
    <!-- 已加入的家庭成员列表 -->
    <view class="family-list">
      <view v-for="(member, index) in allFamilyMembers" :key="member.id" class="family-member">
        <view class="member-info">
          <image :src="`http://122.51.231.155:8080/static/${member.avatarUrl}`" alt="头像" class="avatar" />
          <view class="info">
            <text class="name">{{ member.nickname }}</text>
          </view>
        </view>
        <view class="actions">
          <button v-if="member.role !== 'admin'" @click.stop="setAsAdmin(member.id)" class="edit-btn">{{ $t('set_as_admin') }}</button>
          <button @click.stop="removeMember(member.id)" class="remove-btn">{{ $t('remove_member') }}</button>
        </view>
      </view>
    </view>

    <!-- 分割线 -->
    <view class="divider"></view>

    <!-- 等待加入的成员列表 -->
    <view class="family-list" v-if="family.waiting_members && family.waiting_members.length > 0">
      <view v-for="(member, index) in family.waiting_members" :key="member.id" class="family-member">
        <view class="member-info">
          <image :src="member.avatarUrl" alt="头像" class="avatar" />
          <view class="info">
            <text class="name">{{ member.nickname }}</text>
          </view>
        </view>
        <view class="actions">
          <button @click.stop="admitMember(member.id)" class="add-btn">{{ $t('admit') }}</button>
          <button @click.stop="rejectMember(member.id)" class="remove-btn">{{ $t('reject') }}</button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useFamilyStore } from '@/stores/family.js';

const { t } = useI18n();

const familyStore = useFamilyStore();
const family = computed(() => familyStore.family);

// 获取当前用户ID（假设存储在localStorage中）
const currentUserId = uni.getStorageSync('userId'); // 根据你的项目实际情况决定
const isCurrentUserAdmin = computed(() => {
  return familyStore.isAdmin(currentUserId);
});

// 合并后的家庭成员列表（包括管理员和普通成员）
const allFamilyMembers = computed(() => {
  return family.value.allMembers || [];
});

onMounted(async () => {
  try {
    await familyStore.getFamilyDetails();
  } catch (error) {
    uni.showToast({ title: t('fetch_family_details_failed'), icon: 'none' });
  }
});

// 设为管理员
const setAsAdmin = async (userId) => {
  try {
    await familyStore.setAdmin(userId);
    uni.showToast({ title: t('set_admin_success'), icon: 'success' });
  } catch (error) {
    uni.showToast({ title: t('set_admin_failed'), icon: 'error' });
  }
};

// 删除成员
const removeMember = async (userId) => {
  try {
    await familyStore.removeFamilyMember(userId);
    uni.showToast({ title: t('remove_member_success'), icon: 'success' });
  } catch (error) {
    uni.showToast({ title: t('remove_member_failed'), icon: 'error' });
  }
};

// 同意加入请求
const admitMember = async (userId) => {
  try {
    await familyStore.admitJoinRequest(userId);
    uni.showToast({ title: t('admit_success'), icon: 'success' });
  } catch (error) {
    uni.showToast({ title: t('admit_failed'), icon: 'error' });
  }
};

// 拒绝加入请求
const rejectMember = async (userId) => {
  try {
    await familyStore.rejectJoinRequest(userId);
    uni.showToast({ title: t('reject_success'), icon: 'success' });
  } catch (error) {
    uni.showToast({ title: t('reject_failed'), icon: 'error' });
  }
};
</script>

<style scoped>
.family-management {
  padding: 20rpx;
  background-color: #f5f5f5;
}

.family-list {
  display: flex;
  flex-direction: column;
  gap: 15rpx;
}

.family-member {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  background-color: #ffffff;
  padding: 15rpx;
  border-radius: 8rpx;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.1);
  position: relative;
  box-sizing: border-box;
  width: 100%;
}

.member-info {
  display: flex;
  align-items: center;
  width: 100%;
  margin-bottom: 10rpx;
  box-sizing: border-box;
}

.avatar {
  width: 50rpx;
  height: 50rpx;
  border-radius: 50%;
  margin-right: 15rpx;
}

.info {
  flex: 1;
}

.name {
  font-size: 18rpx;
  font-weight: bold;
}

.actions {
  display: flex;
  gap: 10rpx;
}

.edit-btn {
  color: #007aff;
  background-color: #e0f0ff;
  padding: 10rpx 20rpx;
  border-radius: 5rpx;
}

.remove-btn {
  color: #ff3b30;
  background-color: #ffe0e0;
  padding: 10rpx 20rpx;
  border-radius: 5rpx;
}

.add-btn {
  color: #34c759;
  background-color: #e0ffe0;
  padding: 10rpx 20rpx;
  border-radius: 5rpx;
}

.divider {
  height: 2rpx;
  background-color: #ddd;
  margin: 20rpx 0;
}
</style>
