<template>
  <view class="settings">
    <view class="header">
      <view @click="goBack" class="back-icon">{{ t('back') }}</view>
      <text class="title">{{ t('lanSettings') }}</text>
      <view class="header-actions">
        <button class="menu-icon"></button>
        <button class="camera-icon"></button>
      </view>
    </view>

    <view class="list">
      <view class="list-item" @click="switchToEn">
        <text>{{ t('lang.en') }}</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item" @click="switchToZhHans">
        <text>{{ t('lang.zh-hans') }}</text>
        <text class="arrow">></text>
      </view>
    </view>
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import { useI18n } from 'vue-i18n'

/* ----------------- Setup ----------------- */
const { t, locale } = useI18n()

/* ----------------- Methods ----------------- */
function goBack() {
  uni.navigateBack()
}

function switchToEn() {
  locale.value = 'en'
  uni.reLaunch({
    url: '/pagesSetting/language/language'
  })
}

function switchToZhHans() {
  locale.value = 'zh-Hans'
  uni.reLaunch({
    url: '/pagesSetting/language/language'
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

.arrow {
  color: #ccc;
}
</style>