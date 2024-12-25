<template>
  <view class="container">
    <view class="icon_holder">
      <img src="/static/images/settings/trash.png" class="delete-icon" />
    </view>
    <view class="text_holder" :style="{ fontSize: fontSizeBig + 'px' }">
      <text style="font-weight: bold;">注销账号</text>
    </view>
    <view class="text_holder" :style="{ fontSize: fontSize + 'px' }">
      <text>账号注销不可恢复，请谨慎操作。</text>
    </view>
    <view class="button_holder">
      <button class="confirm-button" @click="deleteId">注销账号</button>
    </view>
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import { ref } from 'vue'

/* ----------------- Setup ----------------- */
const fontSize = ref(16)
const fontSizeBig = fontSize.value + 4

/* ----------------- Methods ----------------- */
function deleteId() {
  uni.showToast({
    title: '不支持此操作',
    icon: 'error'
  })
}
</script>

<style scoped>
.container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  justify-content: space-between;
}

.icon_holder {
  text-align: center;
  padding-top: 100px;
}

.delete-icon {
  width: 200px;
  height: 200px;
}

.text_holder {
  text-align: center;
  padding-top: 10px;
}

.button_holder {
  text-align: center;
  margin-top: auto;
  padding-bottom: 30px;
}

.confirm-button {
  width: 40%;
  padding: 10px;
  background-color: red;
  color: #fff;
  text-align: center;
  border-radius: 5px;
}
</style>