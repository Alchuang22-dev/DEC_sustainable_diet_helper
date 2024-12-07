<!-- home_servant.vue -->
<template>
  <view class="container">
    <!-- 全屏背景图片 -->
    <image src="/static/images/index/background_index_new.png" class="background-image"></image>

    <!-- 头部标题 -->
    <view v-if="family.status === FamilyStatus.JOINED" class="header">
      <text class="header-title">{{ family.name || $t('default_family_name') }}</text>
    </view>

    <!-- 家庭ID -->
    <text v-if="family.status === FamilyStatus.JOINED" class="list-title">{{ $t('family_info') + family.familyId }}</text>

    <!-- 用户未加入家庭时的视图 -->
    <view v-if="family.status === FamilyStatus.NOT_JOINED" class="no-family-view">
      <image src="https://cdn.pixabay.com/photo/2017/01/13/02/31/family-1976162_1280.png" class="centered-image"></image>
      <view class="family-actions">
        <button class="action-button" @click="showCreateFamilyModal = true">{{ t('create_family') }}</button>
        <button class="action-button" @click="showJoinFamilyModal = true">{{ t('join_family') }}</button>
      </view>
    </view>

    <!-- 待审核状态的视图 -->
    <view v-if="family.status === FamilyStatus.PENDING_APPROVAL" class="pending-view">
      <image src="https://cdn.pixabay.com/photo/2021/09/20/22/15/hourglass-6641967_1280.png" class="pending-image"></image>
      <text class="pending-text">{{ $t('pending_approval') }}</text>
      <button class="cancel-button" @click="handleCancelJoin">{{ $t('cancel_join') }}</button>
    </view>

    <!-- 如果已加入家庭并且是管理员且有waiting_members，则显示提示框 -->
    <view v-if="isCurrentUserAdmin && family.waiting_members && family.waiting_members.length > 0" class="admin-notice">
      <text>{{ $t('new_join_requests') }}</text>
      <button @click="goToMyFamily">{{ $t('handle_requests') }}</button>
    </view>

    <!-- 用户已加入家庭时的家庭管理部分 -->
    <view v-if="family.status === FamilyStatus.JOINED" class="family-management">
      <!-- 提出想吃的菜品 -->
      <view class="dish-proposal">
        <text class="section-title">{{ t('propose_dish') }}</text>
        <input v-model="newDish.name" :placeholder="t('dish_name_placeholder')" class="input"></input>
        <picker mode="selector" :range="dishPreferenceLevels" :value="newDish.preference"
                @change="onDishPreferenceChange" class="picker">
          <view>{{ $t('dish_preference') }}: {{ dishPreferenceLevels[newDish.preference] }}</view>
        </picker>
        <button class="submit-button" @click="submitDishProposal">{{ t('submit_proposal') }}</button>
      </view>

      <!-- 家庭成员的提议 -->
      <view class="dish-list">
        <text class="section-title">{{ $t('family_dish_proposals') }}</text>
        <scroll-view class="dish-scroll" scroll-y>
          <view v-for="dish in sortedDishProposals" :key="dish.id" class="dish-item">
            <uni-list>
              <uni-list-item :title="dish.name" :note="dishPreferenceLevels[dish.preference]"
                             :rightText="dish.proposer" />
            </uni-list>
          </view>
        </scroll-view>
      </view>

      <!-- 共享家庭成员的五大营养成分达标情况 -->
      <view class="shared-data">
        <text class="section-title">{{ $t('shared_family_data') }}</text>
        <!-- 添加家庭碳排放环形图 -->
        <view class="charts">
          <qiun-data-charts :canvas2d="true" canvas-id="familyCarbonChart" type="ring" :opts="carbonRingOpts"
                            :chartData="carbonChartData" />
        </view>
        <!-- 添加家庭五大营养成分达标情况的图表 -->
        <view class="charts">
          <qiun-data-charts :canvas2d="true" canvas-id="familyNutrientChart" type="column"
                            :opts="nutrientChartOpts" :chartData="nutrientChartData" />
        </view>
      </view>

      <!-- 家庭成员部分 -->
      <view class="family-info">
        <text class="section-title">{{ $t('family_members') }}</text>
        <view class="family-members">
          <view v-for="member in allFamilyMembers" :key="member.id" class="member">
            <image :src="`http://122.51.231.155:8080/static/${member.avatarUrl}`" class="member-avatar"></image>
            <text class="member-name">{{ `${member.nickname} (${member.role === 'admin' ? $t('admin') : $t('member')})` }}</text>
          </view>
        </view>
        <button class="manage-members-button" @click="manageMembers">{{ t('manage_members') }}</button>
      </view>

      <!-- 底部添加退出和解散家庭按钮 -->
      <view class="bottom-actions">
        <button @click="handleLeaveFamily" class="leave-btn">{{ $t('leave_family') }}</button>
        <button v-if="isCurrentUserAdmin" @click="handleBreakFamily" class="break-btn">{{ $t('break_family') }}</button>
      </view>
    </view>

    <!-- 创建家庭的模态框 -->
    <view v-if="showCreateFamilyModal" class="modal">
      <view class="modal-content">
        <text class="modal-title">{{ $t('create_family') }}</text>
        <input v-model="newFamilyName" :placeholder="t('family_name_placeholder')" class="input"></input>
        <button class="modal-button" @click="createFamily">{{ $t('confirm') }}</button>
        <button class="modal-button cancel" @click="showCreateFamilyModal = false">{{ $t('cancel') }}</button>
      </view>
    </view>

    <!-- 加入家庭的模态框 -->
    <view v-if="showJoinFamilyModal" class="modal">
      <view class="modal-content">
        <text class="modal-title">{{ $t('join_family') }}</text>
        <input v-model="joinFamilyId" :placeholder="t('family_id_placeholder')" class="input"></input>
        <button class="modal-button" @click="joinFamily">{{ $t('confirm') }}</button>
        <button class="modal-button cancel" @click="showJoinFamilyModal = false">{{ $t('cancel') }}</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useFamilyStore, FamilyStatus } from '../../stores/family.js';
import { useUserStore } from "../../stores/user.js";

// 国际化
const { t } = useI18n();

// Pinia 状态管理
const familyStore = useFamilyStore();
const family = computed(() => familyStore.family);
const FamilyStatusEnum = FamilyStatus;
const userStore = useUserStore();

// familyStore.reset();

// 确保用户已登录，否则跳转到登录页面
if (!userStore.user.isLoggedIn) {
  uni.navigateTo({
    url: '/pagesMy/login/login',
  });
}

// 获取当前用户ID
const currentUserId = computed(() => userStore.user.uid);

// 判断当前用户是否为管理员
const isCurrentUserAdmin = computed(() => {
  return familyStore.isAdmin(currentUserId.value);
});

// 合并后的家庭成员列表（包括管理员和普通成员）
const allFamilyMembers = computed(() => {
  return family.value.allMembers || [];
});

// 新菜品提议
const newDish = ref({
  name: '',
  preference: 0,
});

// 菜品偏好级别
const dishPreferenceLevels = computed(() => [
  t('preference_low'),
  t('preference_medium'),
  t('preference_high'),
]);

// 模态框显示状态
const showCreateFamilyModal = ref(false);
const showJoinFamilyModal = ref(false);

// 创建家庭相关
const newFamilyName = ref('');

// 加入家庭相关
const joinFamilyId = ref('');

// 定时器引用
let statusCheckTimer = null;

// 定时检查状态
const startStatusCheck = () => {
  if (statusCheckTimer) {
    clearInterval(statusCheckTimer);
  }
  statusCheckTimer = setInterval(async () => {
    if (family.value.status === FamilyStatusEnum.PENDING_APPROVAL) {
      await familyStore.getFamilyDetails();
    }
  }, 30000); // 30秒
};

// 取消加入申请
const handleCancelJoin = async () => {
  try {
    await familyStore.cancelJoinRequest();
    uni.showToast({
      title: t('cancel_success'),
      icon: 'success'
    });
  } catch (error) {
    uni.showToast({
      title: t('cancel_failed'),
      icon: 'error'
    });
  }
};

// 提交菜品提议
const submitDishProposal = async () => {
  if (newDish.value.name.trim() === '') {
    uni.showToast({
      title: t('dish_name_required'),
      icon: 'none'
    });
    return;
  }
  try {
    await familyStore.addDishProposal({
      name: newDish.value.name,
      preference: newDish.value.preference,
      proposer: userStore.user.nickName || 'You',
    });
    newDish.value.name = '';
    newDish.value.preference = 0;
    uni.showToast({
      title: t('submit_success'),
      icon: 'success'
    });
  } catch (error) {
    uni.showToast({
      title: t('submit_failed'),
      icon: 'error'
    });
  }
};

// 创建家庭
const createFamily = async () => {
  if (newFamilyName.value.trim() === '') {
    uni.showToast({
      title: t('family_name_required'),
      icon: 'none'
    });
    return;
  }

  try {
    const result = await familyStore.createFamily(newFamilyName.value);
    newFamilyName.value = '';
    showCreateFamilyModal.value = false;

    uni.showToast({
      title: t('create_family_success'),
      icon: 'success'
    });
  } catch (error) {
    uni.showToast({
      title: t('create_family_failed'),
      icon: 'error'
    });
  }
};

// 加入家庭
const joinFamily = async () => {
  if (joinFamilyId.value.trim() === '') {
    uni.showToast({
      title: t('family_id_required'),
      icon: 'none'
    });
    return;
  }

  try {
    const searchResult = await familyStore.searchFamily(joinFamilyId.value);

    if (!searchResult || !searchResult.id) {
      uni.showToast({
        title: t('family_not_found'),
        icon: 'error'
      });
      return;
    }

    await familyStore.joinFamily(searchResult.id);
    joinFamilyId.value = '';
    showJoinFamilyModal.value = false;
    startStatusCheck();

    uni.showToast({
      title: t('join_request_sent'),
      icon: 'success'
    });

  } catch (error) {
    uni.showToast({
      title: t('join_family_failed'),
      icon: 'error'
    });
  }
};

// 生命周期钩子
onMounted(async () => {
  try {
    await familyStore.getFamilyDetails();

    if (family.value.status === FamilyStatusEnum.PENDING_APPROVAL) {
      startStatusCheck();
    }

    if (family.value.status === FamilyStatusEnum.JOINED) {
      // 检查是否有待审核成员并且当前用户是管理员
      if (isCurrentUserAdmin.value && family.value.waiting_members.length > 0) {
        // 这里可以触发显示消息框的逻辑
        // 例如设置一个状态来显示消息框
      }
    }

  } catch (error) {
    uni.showToast({ title: t('fetch_family_failed'), icon: 'none' });
  }
});

// 监听用户登录状态变化
watch(() => userStore.user.isLoggedIn, (newVal) => {
  if (!newVal) {
    // 用户登出后，执行相应操作，如重置家庭状态
    familyStore.resetFamily();
    uni.navigateTo({
      url: '/pagesMy/login/login',
    });
  }
});

// 跳转到 myFamily 页面
const goToMyFamily = () => {
  uni.navigateTo({ url: '/pagesMy/myFamily/myFamily' });
};

// 退出家庭
const handleLeaveFamily = async () => {
  try {
    await familyStore.leaveFamily();
    uni.showToast({ title: t('leave_family_success'), icon: 'success' });
  } catch (error) {
    uni.showToast({ title: t('leave_family_failed'), icon: 'error' });
  }
};

// 解散家庭
const handleBreakFamily = async () => {
  try {
    await familyStore.breakFamily();
    uni.showToast({ title: t('break_family_success'), icon: 'success' });
  } catch (error) {
    uni.showToast({ title: t('break_family_failed'), icon: 'error' });
  }
};

// 处理成员管理
const manageMembers = () => {
  uni.navigateTo({ url: '/pagesMy/myFamily/myFamily' });
};

// 处理偏好选择变化
const onDishPreferenceChange = (e) => {
  newDish.value.preference = parseInt(e.detail.value, 10);
};

// 计算属性：排序后的菜品提议
const sortedDishProposals = computed(() => {
  if (!family.value.dishProposals) return [];
  return [...family.value.dishProposals].sort((a, b) => b.preference - a.preference);
});

// 计算属性：家庭成员五大营养成分达标情况图表数据（示例）
const carbonRingOpts = {
  // 图表配置
};
const carbonChartData = {
  // 图表数据
};
const nutrientChartOpts = {
  // 图表配置
};
const nutrientChartData = {
  // 图表数据
};

// // 监听用户 Token 过期并自动刷新（可选）
// watch(() => userStore.user.tokenExpiry, (newExpiry) => {
//   if (newExpiry) {
//     const currentTime = Date.now();
//     const timeToExpiry = newExpiry - currentTime;
//     const timeBeforeRefresh = timeToExpiry - (20 * 1000); // 20秒前刷新
//
//     if (timeBeforeRefresh > 0) {
//       setTimeout(async () => {
//         try {
//           await userStore.refreshToken();
//         } catch (error) {
//           console.error('Token 刷新失败:', error);
//         }
//       }, timeBeforeRefresh);
//     } else {
//       // 如果时间已经不足20秒，立即刷新
//       userStore.refreshToken().catch(err => {
//         console.error('Token 刷新失败:', err);
//       });
//     }
//   }
// });
</script>

<style scoped>
/* 全局样式变量 */
:root {
  --primary-color: #4CAF50;
  --secondary-color: #2fc25b;
  --background-color: #f5f5f5;
  --card-background: rgba(255, 255, 255, 0.8);
  --text-color: #333;
  --shadow-color: rgba(0, 0, 0, 0.1);
  --font-size-title: 32rpx;
  --font-size-subtitle: 24rpx;
}

/* 容器 */
.container {
  display: flex;
  flex-direction: column;
  background-color: var(--background-color);
  min-height: 100vh;
  padding-bottom: 80rpx;
  position: relative;
  overflow: hidden;
}

/* 全屏背景图片 */
.background-image {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  z-index: -1;
  opacity: 0.1;
}

/* 头部标题 */
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

/* 家庭ID */
.list-title {
  margin-left: 10rpx;
  font-size: 20rpx;
  font-weight: bold;
  color: var(--text-color);
  margin-bottom: 20rpx;
  text-align: center;
}

/* 用户未加入家庭时的视图 */
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

/* 待审核状态的样式 */
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
  background-color: #fff3cd;
  padding: 10rpx;
  margin: 10rpx;
  border-radius: 8rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.admin-notice text {
  color: #856404;
}

.admin-notice button {
  background-color: #ffeeba;
  border: none;
  padding: 8rpx 16rpx;
  border-radius: 5rpx;
  cursor: pointer;
}

/* 家庭管理部分 */
.family-management {
  flex: 1;
  padding: 20rpx;
}

/* 创建或加入家庭 */
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

/* 提出想吃的菜品 */
.dish-proposal {
  background-color: var(--card-background);
  padding: 20rpx;
  border-radius: 10rpx;
  box-shadow: 0 2rpx 5rpx var(--shadow-color);
  margin-bottom: 20rpx;
}

.section-title {
  font-size: var(--font-size-title);
  font-weight: bold;
  color: var(--primary-color);
  margin-bottom: 10rpx;
}

/* 输入框 */
.input {
  width: 100%;
  padding: 10rpx;
  border: 1rpx solid #ccc;
  border-radius: 5rpx;
  margin-bottom: 20rpx;
  margin-top: 20rpx;
}

.picker {
  width: 100%;
  padding: 10rpx;
  border: 1rpx solid #ccc;
  border-radius: 5rpx;
  margin-bottom: 20rpx;
  margin-top: 20rpx;
  color: #666666;
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

/* 菜品列表 */
.dish-list {
  background-color: var(--card-background);
  padding: 20rpx;
  border-radius: 10rpx;
  box-shadow: 0 2rpx 5rpx var(--shadow-color);
  margin-bottom: 20rpx;
}

/* Scroll-View 滚动区域 */
.dish-scroll {
  max-height: 300rpx;
  overflow-y: auto;
}

/* 菜品项 */
.dish-item {
  padding: 10rpx;
  border-bottom: 1rpx solid #eee;
}

/* 共享数据 */
.shared-data {
  background-color: var(--card-background);
  padding: 20rpx;
  border-radius: 10rpx;
  box-shadow: 0 2rpx 5rpx var(--shadow-color);
  margin-bottom: 20rpx;
}

.charts {
  width: 100%;
  height: 300rpx;
  margin-bottom: 20rpx;
}

/* 家庭成员部分 */
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

/* 管理成员按钮 */
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
.bottom-actions {
  margin-top: 20rpx;
  display: flex;
  justify-content: space-between;
}

.leave-btn {
  background-color: #ff9800;
  color: #fff;
  padding: 15rpx 30rpx;
  border-radius: 10rpx;
  font-size: 28rpx;
}

.break-btn {
  background-color: #f44336;
  color: #fff;
  padding: 15rpx 30rpx;
  border-radius: 10rpx;
  font-size: 28rpx;
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
  z-index: 10;
}

.modal-content {
  background-color: #fff;
  padding: 20rpx;
  border-radius: 10rpx;
  width: 80%;
}

.modal-title {
  font-size: var(--font-size-subtitle);
  color: var(--primary-color);
  margin-bottom: 15rpx;
  text-align: center;
}

.modal-button {
  background-color: var(--primary-color);
  color: #fff;
  padding: 15rpx;
  border-radius: 10rpx;
  font-size: var(--font-size-subtitle);
  width: 100%;
  text-align: center;
  margin-top: 10rpx;
}

.modal-button.cancel {
  background-color: #ccc;
}

/* 动画效果 */
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

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
