<template>
  <view class="settings">
    <view class="header">
      <text class="title">{{ t('settings_app_settings') }}</text>
      <view class="header-actions">
        <button class="menu-icon"></button>
        <button class="camera-icon"></button>
      </view>
    </view>

    <view class="list">
      <view class="divider">
        <text class="divider-text">{{ t('settings_display') }}</text>
      </view>
      <view class="list-item" @click="navigateTo('language')">
        <text>{{ t('settings_language') }}</text>
        <text class="arrow">></text>
      </view>

      <view class="divider">
        <text class="divider-text">{{ t('settings_storage_data') }}</text>
      </view>
      <view class="list-item" @click="Seal">
        <text>{{ t('settings_clear_data') }}</text>
        <text class="arrow">></text>
      </view>

      <view class="divider">
        <text class="divider-text">{{ t('settings_accessibility') }}</text>
      </view>
      <view class="list-item" @click="Seal">
        <text>{{ t('settings_development_features') }}</text>
        <text class="arrow">></text>
      </view>

      <view class="divider">
        <text class="divider-text">{{ t('settings_about_software') }}</text>
      </view>
      <view class="list-item" @click="navigateTo('FeedBack')">
        <text>{{ t('settings_help_feedback') }}</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item" @click="navigateTo('SoftwareInfo')">
        <text>{{ t('settings_software_info') }}</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item" @click="navigateTo('ConnectUs')">
        <text>{{ t('settings_contact_us') }}</text>
        <text class="arrow">></text>
      </view>
    </view>
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import { useI18n } from 'vue-i18n'

/* ----------------- Setup ----------------- */
const { t } = useI18n()

/* ----------------- Methods ----------------- */
function navigateTo(link) {
  uni.navigateTo({
    url: `/pagesSetting/${link}/${link}`,
  })
}

function Seal() {
  uni.showToast({
    title: '正在开发',
    icon: 'error',
    duration: 2000,
  })
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
  position: relative;
  height: 20px;
  margin: 20px 0;
}
.divider::before {
  content: "";
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background-color: #ebebeb;
  transform: translateY(-50%);
}
.divider-text {
  position: relative;
  z-index: 1;
  background-color: #f8f8f8;
  padding: 0 10px;
  text-align: center;
  font-size: 14px;
  color: #888;
}
.red-text {
  color: red;
}
</style>