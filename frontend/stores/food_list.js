// stores/food_list.js
import { defineStore } from 'pinia';
import { reactive, ref } from 'vue';

export const useFoodListStore = defineStore('foodList', () => {
    // 定义食物列表状态
    const foodList = reactive([
        {
            name: "西红柿",
            weight: "1kg",
            price: "5元",
            transportMethod: "陆运",
            foodSource: "本地",
            image: "",
            isAnimating: false,
            emission: 0,
            calories: 100,
            protein: 200,
            fat: 400,
            carbohydrates: 300,
            sodium: 500
        },
        // 可以预先添加更多食物项
    ]);

    // 定义一个加载状态
    const loaded = ref(false);

    // 加载食物列表从本地存储
    const loadFoodList = () => {
        const storedFoodList = uni.getStorageSync('foodDetails');
        if (storedFoodList && storedFoodList.length > 0) {
            foodList.splice(0, foodList.length, ...storedFoodList.map(food => ({
                ...food,
                isAnimating: false
            })));
        }
        loaded.value = true;
    };

    // 保存食物列表到本地存储
    const saveFoodList = () => {
        uni.setStorageSync('foodDetails', foodList);
    };

    // 添加一个新的食物项
    const addFood = (newFood) => {
        foodList.push(newFood);
        saveFoodList();
    };

    // 删除指定索引的食物项
    const deleteFood = (index) => {
        foodList.splice(index, 1);
        saveFoodList();
    };

    // 更新指定索引的食物项
    const updateFood = (index, updatedFood) => {
        if (foodList[index]) {
            Object.assign(foodList[index], updatedFood);
            saveFoodList();
        }
    };

    // 初始化加载
    if (!loaded.value) {
        loadFoodList();
    }

    return {
        foodList,
        addFood,
        deleteFood,
        updateFood,
        saveFoodList,
        loadFoodList,
        loaded
    };
});
