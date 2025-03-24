import torch
import torch.nn as nn

class SimpleNet(nn.Module):
    def __init__(self):
        super(SimpleNet, self).__init__()
        # 第一个卷积层
        # 输入: 3通道(RGB图像) -> 输出: 16个特征图
        # 例如: 一张32x32的彩色图片(3x32x32) -> 16个32x32的特征图
        self.conv1 = nn.Conv2d(3, 16, 3, padding=1)  # 448参数
        
        # 第二个卷积层
        # 输入: 16个特征图 -> 输出: 32个特征图
        # 例如: 16个16x16的特征图 -> 32个16x16的特征图
        self.conv2 = nn.Conv2d(16, 32, 3, padding=1)  # 4,640参数
        
        # 全连接层
        # 输入: 32 * 8 * 8 = 2048个特征 -> 输出: 10个类别
        # 例如: 32个8x8的特征图展平后 -> 10个数字(0-9)的概率
        self.fc1 = nn.Linear(32 * 64 * 64, 101)  # 20,490参数

    def forward(self, x):
        # 输入x的形状: [batch_size, 3, 32, 32]
        # 例如: [32, 3, 32, 32] 表示32张32x32的RGB图片
        
        # 第一个卷积层 + ReLU激活
        # 输出形状: [batch_size, 16, 32, 32]
        # 例如: [32, 16, 32, 32] 表示32张图片，每张有16个32x32的特征图
        x = torch.relu(self.conv1(x))
        
        # 最大池化层，将特征图尺寸减半
        # 输出形状: [batch_size, 16, 16, 16]
        # 例如: [32, 16, 16, 16] 表示32张图片，每张有16个16x16的特征图
        x = torch.max_pool2d(x, 2)
        
        # 第二个卷积层 + ReLU激活
        # 输出形状: [batch_size, 32, 16, 16]
        # 例如: [32, 32, 16, 16] 表示32张图片，每张有32个16x16的特征图
        x = torch.relu(self.conv2(x))
        
        # 最大池化层，再次将特征图尺寸减半
        # 输出形状: [batch_size, 32, 8, 8]
        # 例如: [32, 32, 8, 8] 表示32张图片，每张有32个8x8的特征图
        x = torch.max_pool2d(x, 2)
        
        # 将特征图展平成一维向量
        # 输出形状: [batch_size, 32 * 8 * 8]
        # 例如: [32, 2048] 表示32张图片，每张图片的特征被展平成2048维向量
        x = x.view(x.size(0), -1)
        
        # 全连接层，得到最终的类别预测
        # 输出形状: [batch_size, 10]
        # 例如: [32, 10] 表示32张图片，每张图片对应10个类别的预测概率
        x = self.fc1(x)
        return x 