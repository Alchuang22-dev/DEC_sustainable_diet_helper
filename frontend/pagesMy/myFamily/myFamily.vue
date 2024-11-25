<template>
  <view class="family-management">
    <view class="family-list">
      <view v-for="(member, index) in familyMembers" :key="index" class="family-member">
        <view class="member-info" @click="toggleDetails(index)">
          <image :src="member.avatar" alt="头像" class="avatar" />
          <view class="info">
            <text class="name">{{ member.name }}</text>
            <text class="role">{{ member.role }}</text>
          </view>
        </view>
        <view class="actions">
          <button @click.stop="editMember(index)" class="edit-btn">编辑</button>
          <button @click.stop="removeMember(index)" class="remove-btn">删除</button>
        </view>
        <view v-if="member.showDetails" class="details">
          <view class='wrapper'>   
            <view class="editor-wrapper">
              <editor id="editor" class="ql-container" placeholder="开始输入..." show-img-size show-img-toolbar
                show-img-resize @statuschange="onStatusChange" :read-only="readOnly" @ready="onEditorReady">
              </editor>
            </view>
          </view>
        </view>
      </view>
    </view>
    <view class="add-member">
      <button @click="addMember" class="add-btn">添加家庭成员</button>
    </view>

    <!-- 新增的 uni-drawer 用于编辑家庭成员信息 -->
    <uni-drawer v-if="drawerVisible" :show="drawerVisible" mode="right" @close="closeDrawer">
      <view class="drawer-content">
        <view class="input-wrapper">
          <input v-model="editedName" placeholder="请输入新的姓名" class="uni-input" />
        </view>
        <view class="input-wrapper">
          <input v-model="editedRole" placeholder="请输入新的家庭位置" class="uni-input" />
        </view>
        <button @click="confirmEdit" class="confirm-btn">确认修改</button>
      </view>
    </uni-drawer>
  </view>
</template>

<script setup>
import { ref } from 'vue';
import UniDrawer from "@/uni_modules/uni-drawer/components/uni-drawer/uni-drawer.vue";

const familyMembers = ref([
  { name: '张三', role: '父亲', avatar: 'path/to/avatar1.png', showDetails: false },
  { name: '李四', role: '母亲', avatar: 'path/to/avatar2.png', showDetails: false },
]);

const drawerVisible = ref(false);
const editedIndex = ref(null);
const editedName = ref('');
const editedRole = ref('');

const addMember = () => {
  const newMember = { name: '新成员', role: '未知', avatar: 'path/to/default_avatar.png', showDetails: false };
  familyMembers.value.push(newMember);
};

const removeMember = (index) => {
  familyMembers.value.splice(index, 1);
};

const editMember = (index) => {
  editedIndex.value = index;
  editedName.value = familyMembers.value[index].name;
  editedRole.value = familyMembers.value[index].role;
  drawerVisible.value = true;
};

const toggleDetails = (index) => {
  familyMembers.value[index].showDetails = !familyMembers.value[index].showDetails;
};

const closeDrawer = () => {
  drawerVisible.value = false;
};

const confirmEdit = () => {
  if (editedIndex.value !== null) {
    familyMembers.value[editedIndex.value].name = editedName.value;
    familyMembers.value[editedIndex.value].role = editedRole.value;
  }
  closeDrawer();
};
</script>

<style scoped>
.family-management {
  padding: 20px;
  background-color: #f5f5f5;
}
.family-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.family-member {
  display: flex;
  flex-direction: column; /* 改为列方向，以便各部分垂直排列 */
  align-items: flex-start;
  background-color: #ffffff;
  padding: 15px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: relative;
  box-sizing: border-box; /* 计算边框和内边距 */
  width: 100%; /* 确保家庭成员容器不会超出父级宽度 */
}

.member-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  width: 100%; /* 确保占据整个宽度 */
  margin-bottom: 10px; /* 增加底部间距，让详情部分与信息分隔开 */
  box-sizing: border-box; /* 防止边距导致溢出 */
}

.details {
  width: 100%; /* 占据整行宽度 */
  margin-top: 10px;
  padding: 10px;
  background-color: #f9f9f9;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  box-sizing: border-box; /* 防止内容溢出父级容器 */
}

.avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  margin-right: 15px;
}
.info {
  flex: 1;
}
.name {
  font-size: 18px;
  font-weight: bold;
}
.role {
  font-size: 14px;
  color: #888;
}
.actions {
  display: flex;
  gap: 10px;
}
.icon-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 16px;
}
.edit-btn {
  color: #007aff;
}
.remove-btn {
  color: #ff3b30;
}
.details {
  margin-top: 10px;
  padding: 10px;
  background-color: #f9f9f9;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
}
.detail-edit-btn {
  margin-top: 10px;
  background-color: #007aff;
  color: #fff;
  border: none;
  padding: 8px 12px;
  border-radius: 4px;
  cursor: pointer;
}
.add-member {
  margin-top: 20px;
  text-align: center;
}
.add-btn {
  background-color: #34c759;
  color: #fff;
  border: none;
  padding: 8px 12px;
  border-radius: 4px;
  cursor: pointer;
}

.drawer-content {
  padding: 20px;
}
.input-wrapper {
  margin-bottom: 15px;
}
.confirm-btn {
  background-color: #007aff;
  color: #fff;
  border: none;
  padding: 10px;
  border-radius: 4px;
  cursor: pointer;
  width: 100%;
}
</style>
