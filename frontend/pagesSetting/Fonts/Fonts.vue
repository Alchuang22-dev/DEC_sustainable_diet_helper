<template>
  <view class="font-setting">
    <view class="preview" :style="{ fontSize: fontSize + 'px' }">
      在此处预览字体设置
    </view>
	<view class="preview" :style="{ fontSize: fontSize + 'px' }">
	  View the size of font here
	</view>
    <view class="slider-labels">
      <text v-for="size in fontSizeOptions" :key="size" :style="{ fontSize: size + 'px' }">
        A
      </text>
    </view>
    <slider
      class="font-slider"
      min="0"
      max="6"
      :value="fontSizeIndex"
      :step="1"
      @change="handleFontSizeChange"
      show-value
    ></slider>
    <button class="confirm-button" @click="confirmFontSize">完成</button>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue';

const fontSizeOptions = [10, 12, 14, 16, 18, 20, 22];
const fontSizeIndex = ref(fontSizeOptions.indexOf(16));
const fontSize = ref(fontSizeOptions[fontSizeIndex.value]);

function handleFontSizeChange(e) {
  fontSizeIndex.value = e.detail.value;
  fontSize.value = fontSizeOptions[fontSizeIndex.value];
}

function confirmFontSize() {
  // 在这里进行保存设置的操作
  uni.setStorage({
    key: 'fontSize',
    data: fontSize.value,
    success: () => {
      uni.showToast({
        title: '设置已保存',
        icon: 'success'
      });
    }
  });
}

onMounted(() => {
  // 初始化字体大小，如果之前有存储的字体大小，读取出来
  uni.getStorage({
    key: 'fontSize',
    success: (res) => {
      if (fontSizeOptions.includes(res.data)) {
        fontSizeIndex.value = fontSizeOptions.indexOf(res.data);
        fontSize.value = res.data;
      }
    }
  });
});
</script>

<style scoped>
.font-setting {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  padding: 20px;
  box-sizing: border-box;
}
.preview {
  margin-bottom: 20px;
  padding: 10px;
  border: 1px solid #ccc;
  width: 80%;
  text-align: center;
}
.slider-labels {
  display: flex;
  justify-content: space-between;
  width: 80%;
  margin-bottom: 10px;
}
.font-slider {
  width: 80%;
  margin-bottom: 20px;
}
.confirm-button {
  width: 40%;
  padding: 10px;
  background-color: #48c079;
  color: #fff;
  text-align: center;
  border-radius: 5px;
}
</style>
