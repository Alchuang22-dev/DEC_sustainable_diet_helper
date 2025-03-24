import os
os.environ['KMP_DUPLICATE_LIB_OK'] = "TRUE"
import torch
import torch.nn as nn
import torch.optim as optim
from torchvision import datasets, transforms
from model import SimpleNet

def main():
    # 设置设备（GPU/CPU）
    device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
    print(f"Using device: {device}")

    # 数据预处理
    transform = transforms.Compose([
        transforms.Resize((256, 256)),  # 新增: 统一调整图像大小
        transforms.ToTensor(),
        transforms.Normalize((0.5, 0.5, 0.5), (0.5, 0.5, 0.5))
    ])

    # 加载数据集
    print("Loading Food-101 dataset...")
    trainset = datasets.Food101(root='./data', split='train', download=True, transform=transform)
    # 创建数据加载器，batch_size=32表示每次处理32张图片
    trainloader = torch.utils.data.DataLoader(trainset, batch_size=32, shuffle=True)

    # 初始化模型
    model = SimpleNet().to(device)  # 将模型移到指定设备
    criterion = nn.CrossEntropyLoss()  # 定义损失函数（交叉熵损失）
    optimizer = optim.SGD(model.parameters(), lr=0.01)  # 定义优化器（随机梯度下降）

    # 训练一个epoch
    print("Training for 1 epoch...")
    model.train()  # 设置为训练模式（启用dropout等训练特定层）
    
    # 遍历数据加载器
    for i, (inputs, labels) in enumerate(trainloader):
        # 将数据移到指定设备（GPU/CPU）
        inputs, labels = inputs.to(device), labels.to(device)
        
        # 前向传播
        optimizer.zero_grad()  # 清空梯度
        outputs = model(inputs)  # 模型前向传播
        loss = criterion(outputs, labels)  # 计算损失
        loss.backward()  # 反向传播，计算梯度
        optimizer.step()  # 更新模型参数

        # 每100个batch打印一次损失
        if (i + 1) % 100 == 0:
            print(f'Batch [{i + 1}], Loss: {loss.item():.4f}')

    print("Training finished!")

if __name__ == '__main__':
    main()
