<template>
  <view class="settings">
    <view class="header">
      <text class="title">{{ $t('settings_user_settings') }}</text>
      <view class="header-actions">
        <button class="menu-icon"></button>
        <button class="camera-icon"></button>
      </view>
    </view>

    <view class="list">
      <view class="list-item" @click="navigateTo('Idinfo')">
        <text>{{ $t('settings_account_security') }}</text>
        <text class="arrow">></text>
      </view>
      <view class="divider"></view>
      <view class="list-item" @click="navigateTo('Permissions')">
        <text>{{ $t('settings_personal_info_collection') }}</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item" @click="navigateTo('PIConnected')">
        <text>{{ $t('settings_third_party_info_collection') }}</text>
        <text class="arrow">></text>
      </view>
      <view class="divider"></view>
      <view class="list-item centered" @click="Seal()">  
        <text>{{ $t('settings_account_binding') }}</text>
      </view>
      <view class="list-item centered red-text" @click="Seal()">
        <text>{{ $t('settings_account_cancellation') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
function goBack() {
  uni.navigateBack();
}

function navigateTo(link) {
  uni.navigateTo({
    url: `/pagesSetting/${link}/${link}`,
  });
}

function Seal() {
	uni.showToast({
	  title: '正在开发',
	  icon: "error",
	  duration: 2000,
	});
}
</script>

<style scoped>
.settings {
  height: 100%;
  background: #f8f8f8;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 16px;
  height: 60px;
  background-color: #fff;
  border-bottom: 1px solid #ebebeb;
}
.title {
  font-size: 18px;
  font-weight: bold;
}
.header-actions button {
  background: none;
  border: none;
}
.list {
  margin-top: 10px;
}
.list-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background-color: #fff;
  border-bottom: 1px solid #ebebeb;
}
.list-item.centered {
  justify-content: center;
}
.arrow {
  color: #ccc;
}
.divider {
  height: 1px;
  background-color: #ebebeb;
  margin: 10px 0;
}
.red-text {
  color: red;
}
</style>
