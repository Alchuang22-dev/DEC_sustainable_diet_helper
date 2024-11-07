<template>
  <view>
    <web-view :src="webUrl"></web-view>
  </view>
</template>

<script>
export default {
  data() {
    return {
      webUrl: '',
    };
  },
  onLoad(options) {
    if (options.url) {
      this.webUrl = decodeURIComponent(options.url);
    }
  },
};
</script>

<style scoped>
web-view {
  width: 100%;
  height: 100vh;
}
</style>
