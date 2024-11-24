<template>
  <view class="settings">
    <view class="list">
      <view class="list-item">
        <text>启用大模型菜谱推荐</text>
        <switch class="toggle-button" @change="toggleModelRecommendation" :checked="isModelRecommendationEnabled" />
      </view>
      <view class="list-item">
        <text>启用大模型图片识别</text>
        <switch class="toggle-button" @change="toggleImageRecognition" :checked="isImageRecognitionEnabled" />
      </view>
      <view class="list-item">
        <text>启用大模型资讯推荐</text>
        <switch class="toggle-button" @change="toggleInfoRecommendation" :checked="isInfoRecommendationEnabled" />
      </view>
      <view class="divider"></view>
      <view class="list-item">
        <text>启用自定义数据库</text>
        <switch class="toggle-button" @change="toggleCustomDatabase" :checked="isCustomDatabaseEnabled" />
      </view>
      <view class="list-item">
        <text>启用自定义菜谱</text>
        <switch class="toggle-button" @change="toggleCustomRecipe" :checked="isCustomRecipeEnabled" />
      </view>
      <view class="divider"></view>
      <view class="list-item">
        <text>引入公益项目</text>
        <switch class="toggle-button" @change="toggleCharityProject" :checked="isCharityProjectEnabled" />
      </view>
      <view class="divider"></view>
      <view class="list-item">
        <text>高级社区功能</text>
        <switch class="toggle-button" @change="toggleAdvancedCommunityFeatures" :checked="isAdvancedCommunityFeaturesEnabled" />
      </view>
      <view class="divider"></view>
      <view class="list-item centered">
        <text>开发者模式</text>
        <switch class="toggle-button" @change="toggleDeveloperMode" :checked="isDeveloperModeEnabled" />
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue';

const isModelRecommendationEnabled = ref(false);
const isImageRecognitionEnabled = ref(false);
const isInfoRecommendationEnabled = ref(false);
const isCustomDatabaseEnabled = ref(false);
const isCustomRecipeEnabled = ref(false);
const isCharityProjectEnabled = ref(false);
const isAdvancedCommunityFeaturesEnabled = ref(false);
const isDeveloperModeEnabled = ref(false);

function toggleModelRecommendation(event) {
  isModelRecommendationEnabled.value = event.detail.value;
}

function toggleImageRecognition(event) {
  isImageRecognitionEnabled.value = event.detail.value;
}

function toggleInfoRecommendation(event) {
  isInfoRecommendationEnabled.value = event.detail.value;
}

function toggleCustomDatabase(event) {
  isCustomDatabaseEnabled.value = event.detail.value;
}

function toggleCustomRecipe(event) {
  isCustomRecipeEnabled.value = event.detail.value;
}

function toggleCharityProject(event) {
  isCharityProjectEnabled.value = event.detail.value;
}

function toggleAdvancedCommunityFeatures(event) {
  isAdvancedCommunityFeaturesEnabled.value = event.detail.value;
}

function toggleDeveloperMode(event) {
  isDeveloperModeEnabled.value = event.detail.value;
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
.toggle-button {
  margin-left: auto;
}
</style>
