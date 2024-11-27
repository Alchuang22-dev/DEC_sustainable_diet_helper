<template>
  <view class="settings">
    <view class="header">
      <view @click="goBack" class="back-icon">{{$t('back')}}</view>
      <text class="title">{{$t('lanSettings')}}</text>
      <view class="header-actions">
        <button class="menu-icon"></button>
        <button class="camera-icon"></button>
      </view>
    </view>

    <view class="list">
      <view class="list-item" @click="switchToEn">
        <text>{{$t('lang.en')}}</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item" @click="switchToZhHans">
        <text>{{$t('lang.zh-hans')}}</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item centered red-text" @click="resetLocale">
        <text>{{$t('resetSystemLanguage')}}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { useI18n } from 'vue-i18n';

const { locale } = useI18n();

function goBack() {
  uni.navigateBack();
}

function switchToEn() {
  locale.value = 'en';
}

function switchToZhHans() {
  locale.value = 'zh-Hans';
}

function resetLocale() {
  const systemLocale = uni.getSystemInfoSync().language.toLowerCase();
  // 您可能需要对 systemLocale 进行处理，以匹配您的语言代码
  locale.value = systemLocale.includes('zh') ? 'zh-hans' : 'en';
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

