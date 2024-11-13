<template>
  <view class="container">
    <view class="header">
      <text class="back-button" @click="navigateBack">&larr;</text>
      <text class="title">家庭共享食谱</text>
      <text class="menu-button">&#9776;</text>
    </view>

    <view class="recipe-container">
      <view class="input-container">
        <input type="text" v-model="recipeName" placeholder="和你家人说一下你今天想吃什么吧~" />
        <button @click="addRecipe">+</button>
      </view>
      <view class="priority-container">
        <picker mode="selector" :range="priorityOptions" :value="priority" @change="updatePriority($event)">
          <view class="picker">{{ priorityOptions[priority] }}</view>
        </picker>
      </view>
    </view>

    <view class="daily-recipes">
      <view class="sort-button-container">
        <button class="sort-button" @click="sortRecipes">按优先级排序</button>
      </view>
      <view class="daily-recipe-list">
        <view v-for="(recipe, index) in recipes" :key="index" class="daily-recipe-item" :data-priority="recipe.priority">
          <text>{{ recipe.name }} - </text>
          <picker mode="selector" :range="portionOptions" :value="recipe.portion" @change="updatePortion(index, $event)">
            <view class="picker">{{ portionOptions[recipe.portion] }}</view>
          </picker>
          <view class="move-buttons">
            <button @click="moveRecipeUp(index)">&#8593;</button>
            <button @click="moveRecipeDown(index)">&#8595;</button>
          </view>
          <button class="delete-button" @click="deleteRecipe(index)">&times;</button>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
export default {
  data() {
    return {
      recipeName: '',
      priority: 0,
      priorityOptions: ['低优先级', '中优先级', '高优先级'],
      portionOptions: ['低优先级', '中优先级', '高优先级'],
      recipes: []
    };
  },
  methods: {
    updatePriority(event) {
      this.priority = event.target.value;
    },
    updatePortion(index, event) {
      this.recipes[index].portion = event.target.value;
    },
    addRecipe() {
      if (!this.recipeName.trim()) {
        uni.showToast({ title: '请输入食谱名称', icon: 'none' });
        return;
      }

      this.recipes.push({
        name: this.recipeName,
        priority: this.priority,
        portion: this.priority
      });
      this.recipeName = '';
      this.priority = 0; // Reset priority selection after adding
    },
    deleteRecipe(index) {
      this.recipes.splice(index, 1);
    },
    moveRecipeUp(index) {
      if (index > 0) {
        const temp = this.recipes[index];
        this.recipes.splice(index, 1);
        this.recipes.splice(index - 1, 0, temp);
      }
    },
    moveRecipeDown(index) {
      if (index < this.recipes.length - 1) {
        const temp = this.recipes[index];
        this.recipes.splice(index, 1);
        this.recipes.splice(index + 1, 0, temp);
      }
    },
    sortRecipes() {
      this.recipes.sort((a, b) => b.portion - a.portion);
    },
    navigateBack() {
      uni.navigateBack();
    }
  }
};
</script>

<style scoped>
.container {
  font-family: 'Arial', sans-serif;
  background-color: #f0f4f7;
  color: #333;
  padding: 0;
}
.header {
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: #ffffff;
  border-bottom: 1px solid #e0e0e0;
  justify-content: space-between;
}
.back-button,
.menu-button {
  font-size: 24px;
  cursor: pointer;
}
.title {
  font-size: 24px;
  font-weight: bold;
}
.recipe-container {
  margin: 20px;
  padding: 20px;
  background-color: #ffffff;
  border-radius: 10px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}
.input-container {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}
input[type='text'] {
  flex: 1;
  padding: 8px;
  border: 1px solid #e0e0e0;
  border-radius: 5px;
  font-size: 14px;
  margin-right: 10px;
}
button {
  padding: 8px;
  border: none;
  border-radius: 50%;
  background-color: #4caf50;
  color: #ffffff;
  font-size: 16px;
  cursor: pointer;
}
.delete-button {
  background-color: #f44336;
  border-radius: 50%;
  font-size: 16px;
  padding: 5px;
  margin-left: 5px;
}
.priority-container {
  margin-bottom: 10px;
}
.picker {
  padding: 8px;
  border-radius: 5px;
  border: 1px solid #e0e0e0;
  font-size: 14px;
}
.daily-recipes {
  margin: 20px;
  padding: 20px;
  background-color: #ffffff;
  border-radius: 10px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}
.daily-recipe-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 5px;
  border-bottom: 1px solid #e0e0e0;
  font-size: 14px;
}
.move-buttons {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-right: 5px;
}
.sort-button-container {
  text-align: center;
  margin-bottom: 10px;
}
.sort-button {
  padding: 8px 16px;
  border-radius: 5px;
  background-color: #4caf50;
  border: none;
  color: #ffffff;
  font-size: 14px;
  cursor: pointer;
}
</style>
