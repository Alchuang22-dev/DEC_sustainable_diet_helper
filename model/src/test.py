import torch
import torchvision.models as models
from PIL import Image, ImageDraw, ImageFont
import torchvision.transforms as transforms
import matplotlib.pyplot as plt
import numpy as np

# 创建 ResNet101 模型
model = models.resnet101(pretrained=False)

# 修改最后一层全连接层为2000类
model.fc = torch.nn.Linear(model.fc.in_features, 2000)

# 加载保存的模型参数
state = torch.load('food2k_resnet101.pth', map_location='cpu')

# 如果 'state_dict' 存在，说明保存的文件包含了模型的权重
if 'state_dict' in state:
    state = state['state_dict']

new_state = {}
for k, v in state.items():
    if k.startswith('module.'):
        new_state[k[len('module.'):]] = v
    else:
        new_state[k] = v

# 加载到模型中
model.load_state_dict(new_state)
model.eval()

# 加载标签文件
def load_labels(label_file='data/label.txt'):
    labels = {}
    with open(label_file, 'r') as f:
        for line in f:
            parts = line.strip().split(' ', 1)
            # print(parts)
            labels[int(parts[0])] = parts[1]
    return labels

labels = load_labels('data/labels.txt')

# 定义图像预处理方法
preprocess = transforms.Compose([
    transforms.Resize(256),
    transforms.CenterCrop(224),
    transforms.ToTensor(),
    transforms.Normalize(mean=[0.485, 0.456, 0.406],
                         std=[0.229, 0.224, 0.225]),
])

# 加载图像并进行预处理
img_path = 'food_image_2.jpg'  # 注意替换
img = Image.open(img_path).convert('RGB')
input_tensor = preprocess(img)
input_batch = input_tensor.unsqueeze(0)  # 增加 batch 维度

# 确保输入图像的形状正确
assert input_batch.shape == (1, 3, 224, 224)

with torch.no_grad():
    output = model(input_batch)
    probabilities = torch.softmax(output, dim=1)
    top_prob, top_catid = torch.topk(probabilities, 10)  # 获取前10个预测

# 检查不确定度
max_prob = top_prob[0][0].item()
if max_prob < 0.10:
    print("Warning: The highest prediction confidence is lower than 10%.")

# 输出前10个预测类别名称和概率
print("Top 10 predictions:")
for i in range(top_catid.size(1)):
    class_id = top_catid[0][i].item()
    class_name = labels.get(class_id, 'Unknown')
    print(f"类别: {class_name}, 概率: {top_prob[0][i].item():.4f}")

# 创建可视化呈现
plt.figure(figsize=(16, 8))
if max_prob < 0.10:
    plt.title("Warning: It seems not like cooked food!")

# 1. 显示输入图像
plt.subplot(1, 2, 1)
plt.imshow(np.array(img))
plt.title('Image', fontsize=15)
plt.axis('off')

# 2. 绘制前10个预测概率的条形图
plt.subplot(1, 2, 2)
top_classes = [labels.get(top_catid[0][i].item(), 'Unknown') for i in range(top_catid.size(1))]
top_probs = [top_prob[0][i].item() for i in range(top_catid.size(1))]

# 反转列表使得最高概率显示在顶部
top_classes.reverse()
top_probs.reverse()

# 画水平条形图
bars = plt.barh(range(len(top_probs)), top_probs, color='skyblue')
plt.yticks(range(len(top_probs)), top_classes, fontsize=12)
plt.xlabel('Probabolity', fontsize=12)
plt.title('Top 10 predictions', fontsize=15)
plt.xlim(0, 1.0)
plt.grid(axis='x', linestyle='--', alpha=0.7)

# 为每个条形添加标签
for i, v in enumerate(top_probs):
    plt.text(v + 0.01, i, f'{v:.4f}', va='center', fontsize=10)

# 调整布局
plt.tight_layout()

# 保存图像
plt.savefig('prediction_result_3.png', dpi=300, bbox_inches='tight')

# 显示图像
plt.show()

