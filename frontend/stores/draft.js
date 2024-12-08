// stores/draft.js
import { defineStore } from 'pinia';

export const useDraftStore = defineStore('draft', {
  state: () => ({
    title: '',
    description: '',
    components: [], // 存储图文组件（文本或图片）
  }),
  actions: {
    // 设置标题
    setTitle(newTitle) {
      this.title = newTitle;
    },
    // 设置简介
    setDescription(newDescription) {
      this.description = newDescription;
    },
    // 添加文本组件
    addTextComponent(content) {
      this.components.push({
        type: 'text',
        content: content,
        id: this.components.length + 1,
      });
    },
    // 添加图片组件
    addImageComponent(imageUrl) {
      this.components.push({
        type: 'image',
        content: imageUrl,
        id: this.components.length + 1,
      });
    },
    // 删除组件
    removeComponent(id) {
      this.components = this.components.filter(item => item.id !== id);
    },
    // 清空草稿
    clearDraft() {
      this.title = '';
      this.description = '';
      this.components = [];
    },
    // 保存草稿（可以做本地存储或API请求）
    saveDraft() {
      const draft = {
        title: this.title,
        description: this.description,
        components: this.components,
      };
      // 例如可以用 `localStorage` 保存草稿
      uni.setStorageSync('draft', JSON.stringify(draft));
    },
    // 加载草稿
    loadDraft() {
      const savedDraft = localStorage.getItem('draft');
      if (savedDraft) {
        const parsedDraft = JSON.parse(savedDraft);
        this.title = parsedDraft.title;
        this.description = parsedDraft.description;
        this.components = parsedDraft.components;
      }
    },
  },
});
